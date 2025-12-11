package middleware

import (
    "net/http"
    "strings"

    "github.com/followCode/djjs-event-reporting-backend/config"
    "github.com/followCode/djjs-event-reporting-backend/app/models"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return config.JWTSecret, nil // dynamic secret from config
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
            c.Abort()
            return
        }

        userID := uint(claims["user_id"].(float64))

        var user models.User
        err = config.DB.First(&user, userID).Error
        if err != nil || user.Token != tokenString {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
            c.Abort()
            return
        }

        // Pass user info to handlers
        c.Set("userID", userID)
        c.Set("roleID", user.RoleID)
        c.Next()
    }
}
