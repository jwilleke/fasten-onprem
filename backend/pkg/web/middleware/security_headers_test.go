package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newSecurityHeadersEngine(httpsEnabled bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeadersMiddleware(httpsEnabled))
	r.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return r
}

func Test_SecurityHeadersMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	newSecurityHeadersEngine(true).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	require.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	require.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	require.Equal(t, "no-referrer", w.Header().Get("Referrer-Policy"))
	//CSP is now enforcing (not report-only)
	csp := w.Header().Get("Content-Security-Policy")
	require.Contains(t, csp, "default-src 'self'")
	require.Contains(t, csp, "script-src 'self'")
	require.Contains(t, csp, "frame-ancestors 'none'")
	require.Empty(t, w.Header().Get("Content-Security-Policy-Report-Only"), "should be enforcing, not report-only")
	//HSTS present when HTTPS enabled
	require.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
}

func Test_SecurityHeadersMiddleware_NoHSTSWithoutHTTPS(t *testing.T) {
	w := httptest.NewRecorder()
	newSecurityHeadersEngine(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	//HSTS must be omitted over plain HTTP
	require.Empty(t, w.Header().Get("Strict-Transport-Security"))
	//other headers still present
	require.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
}
