package handlers

import (
    "log"
	"net/http"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/followCode/djjs-event-reporting-backend/app/services/auth"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *auth.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

// RegisterRequest represents registration payload
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Validate
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	ip := middleware.GetClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	if err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.Name, ip, userAgent); err != nil {
		if err == auth.ErrUserNotFound {
			// Generic error - don't reveal if user exists
			c.JSON(http.StatusConflict, gin.H{"error": "account already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registration successful. please verify your email"})
}

// VerifyEmailRequest represents email verification payload
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmail handles email verification
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		switch err {
		case auth.ErrInvalidToken:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		case auth.ErrTokenExpired:
			c.JSON(http.StatusBadRequest, gin.H{"error": "token expired"})
		case auth.ErrTokenUsed:
			c.JSON(http.StatusBadRequest, gin.H{"error": "token already used"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify email"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

// LoginRequest represents login payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken string       `json:"accessToken"`
	User        UserResponse `json:"user"`
	CsrfToken   string       `json:"csrfToken"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ip := middleware.GetClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	user, accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password, ip, userAgent)
	if err != nil {
		// Log the actual error for debugging (remove in production)
		// fmt.Printf("Login error: %v\n", err)
		
		// Generic error message - don't reveal if email exists
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Set refresh token cookie
	h.setRefreshTokenCookie(c, refreshToken)

	// Set CSRF token cookie and get token value
	csrfToken := middleware.SetCSRFToken(c)

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: accessToken,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		CsrfToken: csrfToken,
	})
}

// RefreshResponse represents refresh token response
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
	CsrfToken   string `json:"csrfToken"`
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		// Log the error for debugging
		log.Printf("[Refresh] Failed to get refresh_token cookie: %v", err)
		log.Printf("[Refresh] Request cookies: %v", c.Request.Cookies())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	if refreshToken == "" {
		log.Printf("[Refresh] Refresh token cookie is empty")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	accessToken, newRefreshToken, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		log.Printf("[Refresh] Refresh token validation failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// Set new refresh token cookie
	h.setRefreshTokenCookie(c, newRefreshToken)

	// Update CSRF token cookie and get token value
	csrfToken := middleware.SetCSRFToken(c)

	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken: accessToken,
		CsrfToken:   csrfToken,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		// No user context, just clear cookies
		h.clearAuthCookies(c)
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
		return
	}

	// Get refresh token from cookie
	refreshToken, _ := c.Cookie("refresh_token")

	// Revoke session
	_ = h.authService.Logout(c.Request.Context(), refreshToken, userID)

	// Clear cookies
	h.clearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// MeResponse represents current user info
type MeResponse struct {
	User UserResponse `json:"user"`
}

// Me returns current user information
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get user from database
	var user auth.User
	err := config.AuthDB.QueryRow(c.Request.Context(),
		`SELECT id, email, name FROM users WHERE id = $1 AND is_deleted = false`,
		userID).Scan(&user.ID, &user.Email, &user.Name)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, MeResponse{
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

// ForgotPasswordRequest represents forgot password payload
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPassword handles password reset request
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ip := middleware.GetClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	// Always return 200 to avoid email enumeration
	_ = h.authService.ForgotPassword(c.Request.Context(), req.Email, ip, userAgent)

	c.JSON(http.StatusOK, gin.H{"message": "if an account exists with that email, a password reset link has been sent"})
}

// ResetPasswordRequest represents password reset payload
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		switch err {
		case auth.ErrInvalidToken:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		case auth.ErrTokenExpired:
			c.JSON(http.StatusBadRequest, gin.H{"error": "token expired"})
		case auth.ErrTokenUsed:
			c.JSON(http.StatusBadRequest, gin.H{"error": "token already used"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
}

// ChangePasswordRequest represents change password payload
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		switch err {
		case auth.ErrInvalidPassword:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid current password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// GetSessionsResponse represents sessions list response
type GetSessionsResponse struct {
	Sessions []SessionResponse `json:"sessions"`
}

// SessionResponse represents session data in API responses
type SessionResponse struct {
	ID               string    `json:"id"`
	UserAgent        string    `json:"userAgent"`
	IP               string    `json:"ip"`
	CreatedAt        time.Time `json:"createdAt"`
	LastUsedAt       time.Time `json:"lastUsedAt"`
	IsCurrentSession bool      `json:"isCurrentSession"`
}

// GetSessions returns all active sessions for the current user
func (h *AuthHandler) GetSessions(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sessionID, _ := middleware.GetSessionID(c)

	sessions, err := h.authService.GetSessions(c.Request.Context(), userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sessions"})
		return
	}

	sessionResponses := make([]SessionResponse, len(sessions))
	for i, s := range sessions {
		sessionResponses[i] = SessionResponse{
			ID:               s.ID,
			UserAgent:        s.UserAgent,
			IP:               s.IP,
			CreatedAt:        s.CreatedAt,
			LastUsedAt:       s.LastUsedAt,
			IsCurrentSession: s.IsCurrentSession,
		}
	}

	c.JSON(http.StatusOK, GetSessionsResponse{
		Sessions: sessionResponses,
	})
}

// RevokeSession revokes a specific session (session ID from URL parameter)
func (h *AuthHandler) RevokeSession(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session id required"})
		return
	}

	if err := h.authService.RevokeSession(c.Request.Context(), userID, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session revoked"})
}

// Helper methods

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	maxAge := int(config.RefreshTokenTTL.Seconds())
	
	// For localhost development, don't set domain (empty string)
	// Empty domain allows cookie to work on localhost
	domain := ""
	
	c.SetCookie(
		"refresh_token",
		token,
		maxAge,
		config.CookiePath,
		domain, // Empty for localhost
		config.CookieSecure, // Should be false for localhost HTTP
		true, // HttpOnly
	)
	
	log.Printf("[setRefreshTokenCookie] Cookie set: path=%s, domain='%s', secure=%v, maxAge=%d", 
		config.CookiePath, domain, config.CookieSecure, maxAge)
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	// Clear refresh token cookie
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		config.CookiePath,
		"",
		config.CookieSecure,
		true,
	)

	// Clear CSRF token cookie
	middleware.ClearCSRFToken(c)
}



