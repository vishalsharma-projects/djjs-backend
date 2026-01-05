package handlers

import (
    "log"
	"net/http"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/followCode/djjs-event-reporting-backend/app/services/auth"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
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

// Register godoc
// @Summary Register a new user
// @Description Register a new user account. An email verification link will be sent to the provided email address.
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body RegisterRequest true "Registration payload"
// @Success 201 {object} map[string]string "Registration successful"
// @Failure 400 {object} map[string]string "Invalid request or validation failed"
// @Failure 409 {object} map[string]string "Account already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/register [post]
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

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify user email address using the verification token sent via email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param verifyEmailRequest body VerifyEmailRequest true "Email verification payload"
// @Success 200 {object} map[string]string "Email verified successfully"
// @Failure 400 {object} map[string]string "Invalid token, expired token, or token already used"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/verify-email [post]
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

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access token. Refresh token is set as HttpOnly cookie.
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /api/auth/login [post]
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

// Refresh godoc
// @Summary Refresh access token
// @Description Refresh the access token using the refresh token from HttpOnly cookie.
// @Tags Auth
// @Produce json
// @Success 200 {object} RefreshResponse "Token refreshed successfully"
// @Failure 401 {object} map[string]string "Refresh token missing or invalid"
// @Router /api/auth/refresh [post]
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

// Logout godoc
// @Summary Logout user
// @Description Logout user and revoke current session. Clears authentication cookies.
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string "Logged out successfully"
// @Router /api/auth/logout [post]
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

// Me godoc
// @Summary Get current user information
// @Description Get the currently authenticated user's information.
// @Tags Auth
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} MeResponse "User information"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Router /api/auth/me [get]
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

// ForgotPassword godoc
// @Summary Request password reset
// @Description Request a password reset link to be sent to the provided email address. Always returns 200 to prevent email enumeration.
// @Tags Auth
// @Accept json
// @Produce json
// @Param forgotPasswordRequest body ForgotPasswordRequest true "Password reset request"
// @Success 200 {object} map[string]string "Password reset link sent (if account exists)"
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /api/auth/forgot-password [post]
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

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using the token received via email from forgot-password endpoint. New password must meet strength requirements.
// @Tags Auth
// @Accept json
// @Produce json
// @Param resetPasswordRequest body ResetPasswordRequest true "Password reset payload"
// @Success 200 {object} map[string]string "Password reset successful"
// @Failure 400 {object} map[string]string "Invalid token, expired token, token already used, or password doesn't meet requirements"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate new password strength
	if err := validators.ValidatePasswordStrength(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// ChangePassword godoc
// @Summary Change password
// @Description Change password for the currently authenticated user. Requires current password. New password must meet strength requirements.
// @Tags Auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param changePasswordRequest body ChangePasswordRequest true "Password change payload"
// @Success 200 {object} map[string]string "Password changed successfully"
// @Failure 400 {object} map[string]string "Invalid request or password doesn't meet requirements"
// @Failure 401 {object} map[string]string "Unauthorized or invalid current password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/change-password [post]
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

	// Validate new password strength
	if err := validators.ValidatePasswordStrength(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if new password is different from current password
	if req.CurrentPassword == req.NewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new password must be different from current password"})
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

// GetSessions godoc
// @Summary Get all active sessions
// @Description Get all active sessions for the currently authenticated user.
// @Tags Auth
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} GetSessionsResponse "List of active sessions"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/sessions [get]
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

// RevokeSession godoc
// @Summary Revoke a session
// @Description Revoke a specific session by session ID for the currently authenticated user.
// @Tags Auth
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} map[string]string "Session revoked successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/sessions/{id} [delete]
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



