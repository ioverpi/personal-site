package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuth(password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check session cookie first
		if cookie, err := c.Cookie("admin_session"); err == nil {
			if subtle.ConstantTimeCompare([]byte(cookie), []byte(password)) == 1 {
				c.Next()
				return
			}
		}

		// Check if this is a login attempt
		if c.Request.Method == "POST" && c.Request.URL.Path == "/admin/login" {
			c.Next()
			return
		}

		// Allow access to login page
		if c.Request.URL.Path == "/admin/login" {
			c.Next()
			return
		}

		// Redirect to login
		c.Redirect(http.StatusFound, "/admin/login")
		c.Abort()
	}
}
