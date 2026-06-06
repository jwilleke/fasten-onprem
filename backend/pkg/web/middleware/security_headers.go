package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware sets baseline security response headers (issue #105 / H4).
//
// The Content-Security-Policy is now ENFORCING. script-src 'self' (no 'unsafe-inline',
// no nonce) is the key XSS defense and the main mitigation for the localStorage-stored
// session token (#103 / H2); the index.html bootstrap scripts were externalized so this
// is safe. Notes on the directives:
//   - style-src allows 'unsafe-inline' (Angular emits inline component styles).
//   - img-src allows https: so externally-referenced images in imported FHIR records
//     still render (images don't execute, so the XSS risk is negligible).
//   - connect-src allows the hello.coop origins for the IdpConnect (third-party) login;
//     everything else (backend API, SSE, the SMART flow) is same-origin. SMART provider
//     auth happens via top-level navigation/backend, not browser fetch.
//
// HSTS is only emitted when HTTPS is enabled (it's meaningless/inappropriate over plain HTTP).
func SecurityHeadersMiddleware(httpsEnabled bool) gin.HandlerFunc {
	const csp = "default-src 'self'; " +
		"script-src 'self'; " +
		"style-src 'self' 'unsafe-inline'; " +
		"img-src 'self' data: https:; " +
		"font-src 'self' data:; " +
		"connect-src 'self' https://wallet.hello.coop https://issuer.hello.coop; " +
		"object-src 'none'; " +
		"frame-ancestors 'none'; " +
		"base-uri 'self'; " +
		"form-action 'self'"

	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Content-Security-Policy", csp)
		if httpsEnabled {
			// HTTPS is on by default (web.listen.https.enabled). 1 year, include subdomains.
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}
