package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/internal/models"
	"github.com/ioverpi/personal-site/internal/services"
)

const (
	SessionCookieName = "session"
	UserContextKey    = "user"
)

// AuthMiddleware validates session tokens and sets user in context
func AuthMiddleware(authService *services.AuthService, secureCookies bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow access to login page and public routes
		path := c.Request.URL.Path
		if path == "/admin/login" || path == "/register" {
			c.Next()
			return
		}

		// Check session cookie
		token, err := c.Cookie(SessionCookieName)
		if err != nil {
			redirectToLogin(c)
			return
		}

		// Validate session and get user
		user, err := authService.GetUserBySession(token)
		if err != nil {
			// Clear invalid cookie
			clearSessionCookie(c, secureCookies)
			redirectToLogin(c)
			return
		}

		// Store user in context for handlers
		c.Set(UserContextKey, user)
		c.Next()
	}
}

// RequireAdmin ensures user has admin role
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUser(c)
		if user == nil || !user.IsAdmin() {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

// GetUser retrieves the authenticated user from context
func GetUser(c *gin.Context) *models.User {
	if user, exists := c.Get(UserContextKey); exists {
		if u, ok := user.(*models.User); ok {
			return u
		}
	}
	return nil
}

// SetSessionCookie sets the session cookie with secure settings
func SetSessionCookie(c *gin.Context, token string, maxAge int, secureCookies bool) {
	sameSite := http.SameSiteLaxMode

	c.SetSameSite(sameSite)
	c.SetCookie(
		SessionCookieName,
		token,
		maxAge,
		"/",
		"",
		secureCookies, // Secure flag (HTTPS only)
		true,          // HttpOnly (no JS access)
	)
}

// clearSessionCookie removes the session cookie
func clearSessionCookie(c *gin.Context, secureCookies bool) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		SessionCookieName,
		"",
		-1,
		"/",
		"",
		secureCookies,
		true,
	)
}

func redirectToLogin(c *gin.Context) {
	c.Redirect(http.StatusFound, "/admin/login")
	c.Abort()
}

// AdminAuth is deprecated - kept for backwards compatibility during migration
// TODO: Remove after Phase 4 controller updates
func AdminAuth(password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/admin/login" {
			c.Next()
			return
		}

		if cookie, err := c.Cookie("admin_session"); err == nil {
			if cookie == password {
				c.Next()
				return
			}
		}

		c.Redirect(http.StatusFound, "/admin/login")
		c.Abort()
	}
}
