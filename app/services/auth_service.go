package services

import (
    "errors"
    "time"

    "github.com/followCode/djjs-event-reporting-backend/config"
    "github.com/followCode/djjs-event-reporting-backend/app/models"
    //"golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
)

func Login(email, password string) (string, error) {
    var user models.User
    err := config.DB.Preload("Role").Where("email = ? AND is_deleted = false", email).First(&user).Error
    if err != nil {
        return "", errors.New("user not found")
    }

    // Compare hashed password - will add this later
    // if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
    //     return "", errors.New("invalid password")
    // }

	if user.Password != password {
    	return "", errors.New("invalid password")
	}

    now := time.Now()
    // Set first login if not set
    if user.FirstLoginOn == nil {
        user.FirstLoginOn = &now
    }

    user.LastLoginOn = &now
    expiry := now.Add(24 * time.Hour)
    user.ExpiredOn = &expiry

    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role_id": user.RoleID,
        "exp":     expiry.Unix(),
    })

    tokenString, err := token.SignedString(config.JWTSecret)
    if err != nil {
        return "", err
    }

    // Save token and audit fields in DB
    user.Token = tokenString
    user.UpdatedOn = &now
    user.UpdatedBy = user.Email // or system user
    if err := config.DB.Save(&user).Error; err != nil {
        return "", err
    }

    return tokenString, nil
}

func Logout(userID uint) error {
    var user models.User
    err := config.DB.First(&user, userID).Error
    if err != nil {
        return err
    }

    now := time.Now()
    user.Token = ""
    user.UpdatedOn = &now
    user.UpdatedBy = user.Email
    return config.DB.Save(&user).Error
}
