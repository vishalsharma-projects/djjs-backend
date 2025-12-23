package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
)

const (
	csrfCookieName = "csrf-token"
	csrfHeaderName = "X-CSRF-Token"
)

// CSRFProtection middleware implements double-submit cookie CSRF protection
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only protect state-changing methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Get CSRF token from cookie
		cookie, err := c.Cookie(csrfCookieName)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "csrf token missing"})
			c.Abort()
			return
		}

		// Get CSRF token from header
		headerToken := c.GetHeader(csrfHeaderName)
		if headerToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "X-CSRF-Token header missing"})
			c.Abort()
			return
		}

		// Compare tokens (constant-time comparison)
		if !constantTimeEqual(cookie, headerToken) {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid csrf token"})
			c.Abort()
			return
		}

		// Validate Origin header
		origin := c.GetHeader("Origin")
		if origin != "" && origin != config.FrontendOrigin {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid origin"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SetCSRFToken sets a CSRF token cookie and returns the token value
func SetCSRFToken(c *gin.Context) string {
	// Generate new token (always generate new one for rotation)
	token := generateCSRFToken()

	// Set cookie (not HttpOnly, so JS can read it for double-submit)
	// Use same expiry as refresh token for consistency
	maxAge := int(config.RefreshTokenTTL.Seconds())
	c.SetCookie(
		csrfCookieName,
		token,
		maxAge,
		"/",
		"",
		config.CookieSecure,
		false, // Not HttpOnly
	)

	// Store in context for handlers that need it
	c.Set("csrf_token", token)
	
	return token
}

// ClearCSRFToken clears the CSRF token cookie
func ClearCSRFToken(c *gin.Context) {
	c.SetCookie(
		csrfCookieName,
		"",
		-1,
		"/",
		"",
		config.CookieSecure,
		false,
	)
}

func generateCSRFToken() string {
	bytes := make([]byte, 32) // 256 bits
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// OptionalCSRFProtection middleware implements CSRF protection but allows requests to proceed
// if CSRF token is missing (useful for refresh endpoint which uses HttpOnly cookies)
func OptionalCSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only protect state-changing methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Get CSRF token from cookie
		cookie, err := c.Cookie(csrfCookieName)
		if err != nil {
			// No CSRF cookie - allow request to proceed (refresh uses HttpOnly cookie for security)
			c.Next()
			return
		}

		// Get CSRF token from header
		headerToken := c.GetHeader(csrfHeaderName)
		if headerToken == "" {
			// No CSRF header - allow request to proceed if cookie exists
			// This handles cases where frontend hasn't set header yet
			c.Next()
			return
		}

		// Compare tokens (constant-time comparison)
		if !constantTimeEqual(cookie, headerToken) {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid csrf token"})
			c.Abort()
			return
		}

		// Validate Origin header if present
		origin := c.GetHeader("Origin")
		if origin != "" && origin != config.FrontendOrigin {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid origin"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func constantTimeEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}


