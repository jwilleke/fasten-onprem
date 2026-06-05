package relay

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

const testSecret = "test-secret"

func newClient(srv *httptest.Server) Client {
	return Client{BaseURL: srv.URL, Secret: testSecret, HTTPClient: srv.Client()}
}

func TestPollSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pending" {
			t.Errorf("path = %q, want /pending", r.URL.Path)
		}
		if got := r.Header.Get("X-Yourphr-Token"); got != testSecret {
			t.Errorf("token header = %q, want %q", got, testSecret)
		}
		if got := r.URL.Query().Get("state"); got != "S1" {
			t.Errorf("state = %q, want S1", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"ABC123"}`))
	}))
	defer srv.Close()

	code, err := newClient(srv).Poll(context.Background(), "S1")
	if err != nil {
		t.Fatalf("Poll: %v", err)
	}
	if code != "ABC123" {
		t.Errorf("code = %q, want ABC123", code)
	}
}

func TestPollNotReady(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := newClient(srv).Poll(context.Background(), "S1")
	if !errors.Is(err, ErrNotReady) {
		t.Errorf("err = %v, want ErrNotReady", err)
	}
}

func TestPollUnauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer srv.Close()

	if _, err := newClient(srv).Poll(context.Background(), "S1"); err == nil || errors.Is(err, ErrNotReady) {
		t.Errorf("expected a hard auth error, got %v", err)
	}
}

func TestPollValidation(t *testing.T) {
	if _, err := (Client{}).Poll(context.Background(), "S1"); err == nil {
		t.Error("expected error when BaseURL/Secret unset")
	}
	if _, err := (Client{BaseURL: "http://x", Secret: "s"}).Poll(context.Background(), ""); err == nil {
		t.Error("expected error when state empty")
	}
}

func TestPollUntilSucceedsAfterRetries(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 404 for the first two polls, then the code on the third.
		if atomic.AddInt32(&hits, 1) < 3 {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(`{"code":"LATE"}`))
	}))
	defer srv.Close()

	code, err := newClient(srv).PollUntil(context.Background(), "S1", 5*time.Millisecond, 2*time.Second)
	if err != nil {
		t.Fatalf("PollUntil: %v", err)
	}
	if code != "LATE" {
		t.Errorf("code = %q, want LATE", code)
	}
	if atomic.LoadInt32(&hits) < 3 {
		t.Errorf("expected at least 3 polls, got %d", hits)
	}
}

func TestPollUntilTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := newClient(srv).PollUntil(context.Background(), "S1", 5*time.Millisecond, 40*time.Millisecond)
	if err == nil || errors.Is(err, ErrNotReady) {
		t.Errorf("expected timeout error, got %v", err)
	}
}
