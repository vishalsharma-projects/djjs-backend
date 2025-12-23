package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/followCode/djjs-event-reporting-backend/app/services/auth"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes sets up authentication routes with proper middleware
func SetupAuthRoutes(r *gin.RouterGroup) {
	// Initialize auth service
	mailer := auth.NewStubMailer()
	authService := auth.NewAuthService(mailer)
	authHandler := handlers.NewAuthHandler(authService)

	// Public routes
	authGroup := r.Group("/auth")
	{
		// Registration
		authGroup.POST("/register",
			middleware.StrictJSONBinding(),
			middleware.RateLimiter(middleware.RateLimitConfig{
				MaxRequests:   5,
				Window:        config.RateLimitWindow,
				IdentifierKey: "ip",
			}),
			authHandler.Register,
		)

		// Email verification
		authGroup.POST("/verify-email",
			middleware.StrictJSONBinding(),
			authHandler.VerifyEmail,
		)

		// Login (rate limited by IP)
		authGroup.POST("/login",
			middleware.StrictJSONBinding(),
			middleware.RateLimiter(middleware.RateLimitConfig{
				MaxRequests:   config.RateLimitLoginPerIP,
				Window:        config.RateLimitWindow,
				IdentifierKey: "ip",
			}),
			authHandler.Login,
		)

		// Refresh token (CSRF optional - uses HttpOnly cookie for security)
		// CSRF is checked but refresh can proceed if cookie is valid even without header
		authGroup.POST("/refresh",
			middleware.OptionalCSRFProtection(),
			authHandler.Refresh,
		)

		// Logout
		authGroup.POST("/logout", authHandler.Logout)

		// Forgot password (rate limited by IP)
		authGroup.POST("/forgot-password",
			middleware.StrictJSONBinding(),
			middleware.RateLimiter(middleware.RateLimitConfig{
				MaxRequests:   config.RateLimitForgotPasswordPerIP,
				Window:        config.RateLimitWindow,
				IdentifierKey: "ip",
			}),
			authHandler.ForgotPassword,
		)

		// Reset password
		authGroup.POST("/reset-password",
			middleware.StrictJSONBinding(),
			authHandler.ResetPassword,
		)
	}

	// Protected routes
	protected := r.Group("/auth")
	protected.Use(middleware.AuthRequired())
	{
		// Get current user
		protected.GET("/me", authHandler.Me)

		// Change password
		protected.POST("/change-password",
			middleware.StrictJSONBinding(),
			authHandler.ChangePassword,
		)

		// Session management
		protected.GET("/sessions", authHandler.GetSessions)
		protected.DELETE("/sessions/:id", authHandler.RevokeSession)
	}
}

