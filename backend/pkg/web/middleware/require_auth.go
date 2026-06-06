package middleware

import (
	"net/http"
	"strings"

	"github.com/fastenhealth/fasten-onprem/backend/pkg"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/auth"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/config"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/database"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		appConfig := c.MustGet(pkg.ContextKeyTypeConfig).(config.Interface)

		// Token transport: Authorization: Bearer is primary (RFC 6750 / SMART App Launch);
		// fall back to the HttpOnly session cookie for browser clients. The header wins if
		// both are present, so non-browser/desktop/CLI clients are unaffected (#103 / H2).
		tokenString := ""
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenString = parts[1]
			}
		}
		if tokenString == "" {
			if cookieVal, err := c.Cookie(pkg.SessionCookieName); err == nil {
				tokenString = cookieVal
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "request does not contain a valid token"})
			c.Abort()
			return
		}
		
		claims, err := auth.JwtValidateFastenToken(appConfig.GetString("jwt.issuer.key"), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
			c.Abort()
			return
		}
		
		if claims.TokenType == "access" {
			databaseRepo := c.MustGet(pkg.ContextKeyTypeDatabase).(database.DatabaseRepository)
			token, err := databaseRepo.GetAccessTokenByTokenIDAndUsername(c, claims.ID, claims.Subject)
			if err != nil || token == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Invalid access token"})
				c.Abort()
				return
			}
		}
		
		// Set context for both regular and access tokens
		c.Set(pkg.ContextKeyTypeAuthToken, tokenString)
		c.Set(pkg.ContextKeyTypeAuthUsername, claims.Subject)
		c.Next()
	}
}
