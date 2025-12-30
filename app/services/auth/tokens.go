package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateRandomToken generates a cryptographically secure random token
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// HashRefreshToken hashes a refresh token with pepper for storage
// Returns the hash as bytes (for BYTEA column)
func HashRefreshToken(token string) []byte {
	// Hash = SHA256(token + pepper)
	h := sha256.New()
	h.Write([]byte(token))
	h.Write(config.TokenPepper)
	return h.Sum(nil)
}

// GenerateAccessToken generates a JWT access token
func GenerateAccessToken(userID int64, sessionID string, roleID int64, roleName string) (string, error) {
	now := time.Now()
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"sub":       fmt.Sprintf("%d", userID), // Subject (user ID)
		"sid":       sessionID,                 // Session ID
		"jti":       jti,                       // JWT ID (unique token identifier)
		"iat":       now.Unix(),                // Issued at
		"exp":       now.Add(config.JWTTTL).Unix(), // Expiration
		"iss":       config.JWTIssuer,          // Issuer
		"aud":       config.JWTAudience,        // Audience
		"role_id":   roleID,                    // Role ID
		"role_name": roleName,                  // Role Name
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyAccessToken verifies and parses an access token, returning the claims
func VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWTSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// ParseUserIDFromToken extracts user ID from JWT claims
func ParseUserIDFromToken(claims jwt.MapClaims) (int64, error) {
	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid sub claim")
	}

	var userID int64
	_, err := fmt.Sscanf(sub, "%d", &userID)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID in sub claim: %w", err)
	}

	return userID, nil
}

// ParseSessionIDFromToken extracts session ID from JWT claims
func ParseSessionIDFromToken(claims jwt.MapClaims) (string, error) {
	sid, ok := claims["sid"].(string)
	if !ok {
		return "", fmt.Errorf("invalid sid claim")
	}
	return sid, nil
}

// ParseRoleIDFromToken extracts role ID from JWT claims
func ParseRoleIDFromToken(claims jwt.MapClaims) (int64, error) {
	// Try to get as float64 (default JSON number type)
	if roleIDFloat, ok := claims["role_id"].(float64); ok {
		return int64(roleIDFloat), nil
	}
	
	// Try to get as int64
	if roleID, ok := claims["role_id"].(int64); ok {
		return roleID, nil
	}
	
	return 0, fmt.Errorf("invalid role_id claim")
}

// ParseRoleNameFromToken extracts role name from JWT claims
func ParseRoleNameFromToken(claims jwt.MapClaims) (string, error) {
	roleName, ok := claims["role_name"].(string)
	if !ok {
		return "", fmt.Errorf("invalid role_name claim")
	}
	return roleName, nil
}

// HashToken hashes a token (for verification/reset tokens) with pepper
func HashToken(token string) []byte {
	h := sha256.New()
	h.Write([]byte(token))
	h.Write(config.TokenPepper)
	return h.Sum(nil)
}

// ConstantTimeCompare compares two byte slices in constant time
func ConstantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}


