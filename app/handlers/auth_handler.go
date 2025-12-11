package handlers

import (
	"log"
	"net/http"

	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/gin-gonic/gin"
)

// LoginRequest represents login payload
// swagger:parameters loginRequest
type LoginRequest struct {
	// in: body
	// required: true
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler godoc
// @Summary Login user and return JWT token
// @Description Authenticate user using email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	token, err := services.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Println("Generated Token:", token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// LogoutHandler godoc
// @Summary Logout user
// @Description Clears user's token and ends session
// @Tags Auth
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /logout [post]
func LogoutHandler(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	userID := userIDInterface.(uint)
	if err := services.Logout(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
