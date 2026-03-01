package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds essential HTTP security headers to all responses.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Enforce HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable basic XSS filtering (mostly legacy browsers, but good practice)
		c.Header("X-XSS-Protection", "1; mode=block")

		// Control framing: Only allow Telegram Web Apps / telegram.org to frame this application
		// "self" allows same-origin framing if needed
		c.Header("Content-Security-Policy", "frame-ancestors 'self' https://*.telegram.org https://telegram.org")
		c.Header("X-Frame-Options", "SAMEORIGIN")

		// Prevent leaking exact route/path details when crossing origin
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Block access to web features that the API definitively does not need
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}
