package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		// - default-src 'self': Only load resources from same origin
		// - script-src 'self' https://unpkg.com: Allow scripts from self and htmx CDN
		// - style-src 'self' 'unsafe-inline': Allow styles from self and inline (for theme toggle)
		// - img-src 'self' data:: Allow images from self and data URIs
		// - connect-src 'self': Only allow AJAX/fetch to same origin
		// - frame-ancestors 'none': Prevent embedding in iframes (clickjacking protection)
		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' https://unpkg.com; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data:; "+
				"connect-src 'self'; "+
				"frame-ancestors 'none'")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Enable browser XSS filter
		c.Header("X-XSS-Protection", "1; mode=block")

		// Only send referrer for same-origin requests
		c.Header("Referrer-Policy", "same-origin")

		c.Next()
	}
}
