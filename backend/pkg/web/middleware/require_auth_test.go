package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fastenhealth/fasten-onprem/backend/pkg"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/auth"
	mock_config "github.com/fastenhealth/fasten-onprem/backend/pkg/config/mock"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/models"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testJWTKey = "test_signing_key"

func tokenFor(t *testing.T, username string) string {
	t.Helper()
	tok, err := auth.JwtGenerateFastenTokenFromUser(models.User{Username: username}, testJWTKey)
	require.NoError(t, err)
	return tok
}

func runRequireAuth(t *testing.T, build func(req *http.Request)) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockConfig := mock_config.NewMockInterface(ctrl)
	mockConfig.EXPECT().GetString("jwt.issuer.key").Return(testJWTKey).AnyTimes()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(pkg.ContextKeyTypeConfig, mockConfig)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/secure/x", nil)
	build(c.Request)

	middleware.RequireAuth()(c)
	return c, w
}

// Cookie fallback: when there's no Authorization header, the HttpOnly session cookie is accepted (#103).
func TestRequireAuth_CookieFallback(t *testing.T) {
	c, _ := runRequireAuth(t, func(req *http.Request) {
		req.AddCookie(&http.Cookie{Name: pkg.SessionCookieName, Value: tokenFor(t, "alice")})
	})
	require.False(t, c.IsAborted(), "a valid cookie token should be accepted")
	require.Equal(t, "alice", c.GetString(pkg.ContextKeyTypeAuthUsername))
}

// The Authorization header is primary: it wins even if a (different) cookie is also present (RFC 6750).
func TestRequireAuth_HeaderWinsOverCookie(t *testing.T) {
	c, _ := runRequireAuth(t, func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+tokenFor(t, "bob"))
		req.AddCookie(&http.Cookie{Name: pkg.SessionCookieName, Value: tokenFor(t, "alice")})
	})
	require.False(t, c.IsAborted())
	require.Equal(t, "bob", c.GetString(pkg.ContextKeyTypeAuthUsername), "header token must take precedence")
}

// No header and no cookie → 401.
func TestRequireAuth_NoToken(t *testing.T) {
	c, w := runRequireAuth(t, func(req *http.Request) {})
	require.True(t, c.IsAborted())
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
