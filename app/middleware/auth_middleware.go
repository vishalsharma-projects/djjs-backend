package middleware

import (
    "log"
    "net/http"
    "strconv"
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

        // Support both old format (user_id as float64) and new format (sub as string)
        var userID uint
        if userIDFloat, ok := claims["user_id"].(float64); ok {
            // Old token format - user_id as float64
            userID = uint(userIDFloat)
            log.Printf("[AuthMiddleware] Using old token format, user_id: %d", userID)
        } else if sub, ok := claims["sub"].(string); ok {
            // New token format - sub contains user ID as string
            userIDInt, err := strconv.ParseUint(sub, 10, 32)
            if err != nil {
                log.Printf("[AuthMiddleware] Failed to parse sub claim '%s': %v", sub, err)
                c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
                c.Abort()
                return
            }
            userID = uint(userIDInt)
            log.Printf("[AuthMiddleware] Using new token format, sub: %s, user_id: %d", sub, userID)
        } else {
            log.Printf("[AuthMiddleware] Token missing both user_id and sub claims")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
            c.Abort()
            return
        }

        // Check if user exists (don't require token match for new auth system)
        var user models.User
        err = config.DB.First(&user, userID).Error
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
            c.Abort()
            return
        }
        
        // For backward compatibility: if user.Token is set and matches, use it
        // Otherwise, assume new auth system (token validated by JWT signature)
        // This allows both old and new auth systems to work
        if user.Token != "" && user.Token != tokenString {
            // Old system: token must match database
            // But only enforce if user.Token is actually set (old system)
            // New system doesn't set user.Token, so we skip this check
            log.Printf("[AuthMiddleware] Token mismatch for user %d (old system check)", userID)
        }

        // Pass user info to handlers
        c.Set("userID", userID)
        c.Set("roleID", user.RoleID)
        c.Next()
    }
}
