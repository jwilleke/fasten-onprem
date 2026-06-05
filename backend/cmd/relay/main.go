// Command relay is the YourPHR SMART on FHIR OAuth store-and-poll relay (EPIC #20, issue #50).
//
// It is a small, stateless public bouncer for the SMART authorization code: the provider
// redirects the user's browser to /callback with ?code&state; the relay stores {state -> code}
// in memory with a short TTL; the (possibly non-public) YourPHR instance polls
// /pending?state= (gated by a shared secret) to retrieve the code and completes the token
// exchange itself. The relay never sees access/refresh tokens and holds no provider app
// registration — it is provider-agnostic and client-agnostic (per-user/BYO model).
//
// See docs/planning/smart-on-fhir/oauth-gateway.md.
package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const defaultTTL = 60 * time.Second

type codeEntry struct {
	code   string
	expiry time.Time
}

// store is an in-memory, TTL'd map of OAuth state -> authorization code. Safe for concurrent use.
type store struct {
	mu      sync.Mutex
	entries map[string]codeEntry
	ttl     time.Duration
	now     func() time.Time // injectable for tests
}

func newStore(ttl time.Duration) *store {
	return &store{entries: map[string]codeEntry{}, ttl: ttl, now: time.Now}
}

func (s *store) put(state, code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[state] = codeEntry{code: code, expiry: s.now().Add(s.ttl)}
}

// take returns the code for state and removes it. ok is false if missing or expired.
func (s *store) take(state string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[state]
	if !ok {
		return "", false
	}
	delete(s.entries, state)
	if s.now().After(e.expiry) {
		return "", false
	}
	return e.code, true
}

func (s *store) evictExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	for k, e := range s.entries {
		if now.After(e.expiry) {
			delete(s.entries, k)
		}
	}
}

// metrics holds the relay's Prometheus counters. Hand-rolled (text exposition) to keep the
// relay pure-stdlib / zero-dependency. Exposed on a separate, in-cluster-only metrics port.
type metrics struct {
	callbacks           atomic.Int64 // codes received + stored at /callback ("a secret arrived")
	callbackErrors      atomic.Int64 // /callback rejected (provider error or missing code/state)
	pendingFound        atomic.Int64 // /pending delivered a code
	pendingNotFound     atomic.Int64 // /pending had no (live) code for the state
	pendingUnauthorized atomic.Int64 // /pending rejected (missing/bad shared secret)
}

// writeProm emits the counters in Prometheus text exposition format (v0.0.4).
func (m *metrics) writeProm(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	fmt.Fprintf(w, "# HELP yourphr_relay_callbacks_total Authorization codes received and stored at /callback.\n")
	fmt.Fprintf(w, "# TYPE yourphr_relay_callbacks_total counter\n")
	fmt.Fprintf(w, "yourphr_relay_callbacks_total %d\n", m.callbacks.Load())
	fmt.Fprintf(w, "# HELP yourphr_relay_callback_errors_total Rejected /callback requests (provider error or missing params).\n")
	fmt.Fprintf(w, "# TYPE yourphr_relay_callback_errors_total counter\n")
	fmt.Fprintf(w, "yourphr_relay_callback_errors_total %d\n", m.callbackErrors.Load())
	fmt.Fprintf(w, "# HELP yourphr_relay_pending_total Poll requests at /pending, by result.\n")
	fmt.Fprintf(w, "# TYPE yourphr_relay_pending_total counter\n")
	fmt.Fprintf(w, "yourphr_relay_pending_total{result=\"found\"} %d\n", m.pendingFound.Load())
	fmt.Fprintf(w, "yourphr_relay_pending_total{result=\"not_found\"} %d\n", m.pendingNotFound.Load())
	fmt.Fprintf(w, "yourphr_relay_pending_total{result=\"unauthorized\"} %d\n", m.pendingUnauthorized.Load())
}

// newServer builds the relay HTTP handler. secret gates /pending; ttl is the code lifetime.
// The returned *store is exposed for tests (TTL injection) and the background janitor; the
// returned *metrics is served on the separate metrics port and asserted in tests.
func newServer(secret string, ttl time.Duration) (http.Handler, *store, *metrics) {
	st := newStore(ttl)
	m := &metrics{}
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// "/" is the ServeMux catch-all: respond with a friendly note at the exact root
	// (a human landing here in a browser), but keep a real 404 for unknown paths.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<!doctype html><h1>YourPHR SMART relay</h1>" +
			"<p>OAuth store-and-poll relay — no content here. " +
			"See <a href=\"https://github.com/jwilleke/yourphr/issues/50\">YourPHR #50</a>.</p>"))
	})

	// /callback: the provider redirects the user's browser here with ?code&state. Open by
	// design (the provider must reach it); it only stores a short-lived code keyed by state.
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if e := q.Get("error"); e != "" {
			m.callbackErrors.Add(1)
			log.Printf("relay: /callback provider error for state=%q: %s — %s", q.Get("state"), e, q.Get("error_description"))
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("<!doctype html><h1>Authorization failed</h1><p>" +
				html.EscapeString(e+" "+q.Get("error_description")) + "</p>"))
			return
		}
		code, state := q.Get("code"), q.Get("state")
		if code == "" || state == "" {
			m.callbackErrors.Add(1)
			http.Error(w, "missing code or state", http.StatusBadRequest)
			return
		}
		st.put(state, code)
		m.callbacks.Add(1)
		// Log that a code arrived (state only — never the code, which is the secret).
		log.Printf("relay: stored authorization code for state=%q (ttl %s)", state, ttl)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<!doctype html><h1>Connected</h1>" +
			"<p>You may close this window and return to YourPHR.</p>"))
	})

	// /pending: the YourPHR instance polls here (shared-secret gated) to retrieve the code.
	mux.HandleFunc("/pending", func(w http.ResponseWriter, r *http.Request) {
		if !secretEqual(r.Header.Get("X-Yourphr-Token"), secret) {
			m.pendingUnauthorized.Add(1)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		state := r.URL.Query().Get("state")
		if state == "" {
			http.Error(w, "missing state", http.StatusBadRequest)
			return
		}
		code, ok := st.take(state)
		if !ok {
			m.pendingNotFound.Add(1)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		m.pendingFound.Add(1)
		log.Printf("relay: delivered authorization code for state=%q", state)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"code": code})
	})

	return mux, st, m
}

// secretEqual is a constant-time comparison that also rejects an empty configured secret.
func secretEqual(got, want string) bool {
	if want == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}

func main() {
	secret := os.Getenv("YOURPHR_RELAY_SECRET")
	if secret == "" {
		log.Fatal("YOURPHR_RELAY_SECRET is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9090"
	}

	handler, st, m := newServer(secret, defaultTTL)

	// Background janitor: evict expired codes so the map cannot grow unbounded.
	go func() {
		t := time.NewTicker(defaultTTL)
		defer t.Stop()
		for range t.C {
			st.evictExpired()
		}
	}()

	// Metrics server on a separate port — scraped in-cluster only, NOT exposed via the public
	// tunnel (which routes the relay's main :PORT). Keeps callback/poll counts off the internet.
	go func() {
		mmux := http.NewServeMux()
		mmux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) { m.writeProm(w) })
		mmux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
		msrv := &http.Server{Addr: ":" + metricsPort, Handler: mmux, ReadHeaderTimeout: 5 * time.Second}
		log.Printf("yourphr relay metrics listening on :%s", metricsPort)
		if err := msrv.ListenAndServe(); err != nil {
			log.Printf("relay: metrics server stopped: %v", err)
		}
	}()

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("yourphr relay listening on :%s (code TTL %s)", port, defaultTTL)
	log.Fatal(srv.ListenAndServe())
}
