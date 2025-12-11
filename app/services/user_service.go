package services

import (
	"errors"
	"math/rand"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"golang.org/x/crypto/bcrypt"
)

// Helper: Generate random 8-character alphanumeric password
func generateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies a plain password against a hashed password
func VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// CreateUser inserts a new user record
func CreateUser(user *models.User) error {
	// Validate that role exists
	var role models.Role
	if err := config.DB.First(&role, user.RoleID).Error; err != nil {
		return errors.New("invalid role_id: role does not exist")
	}

	// Validate email uniqueness
	var existingUser models.User
	if err := config.DB.Where("email = ? AND is_deleted = ?", user.Email, false).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	}

	plainPassword := generateRandomPassword()
	hashedPassword, err := HashPassword(plainPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.CreatedOn = time.Now()
	user.UpdatedOn = nil

	if err := config.DB.Create(user).Error; err != nil {
		return err
	}

	// Return the plain password to the caller for display
	user.Password = plainPassword
	return nil
}

// GetAllUsers fetches all users (excluding deleted)
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := config.DB.Preload("Role").Where("is_deleted = ?", false).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserSearch fetches a single user by email, contact
func GetUserSearch(email, contact string) ([]models.User, error) {
	var users []models.User
	query := config.DB.Model(&models.User{}).Preload("Role")

	// Dynamically build WHERE conditions
	if email != "" && contact != "" {
		query = query.Where("email = ? AND contact_number = ?", email, contact)
	} else if email != "" {
		query = query.Where("email = ?", email)
	} else if contact != "" {
		query = query.Where("contact_number = ?", contact)
	}

	// Execute query
	if err := query.Find(&users).Error; err != nil {
		return nil, errors.New("no users found")
	}

	if len(users) == 0 {
		return nil, errors.New("no users found")
	}

	return users, nil
}

// UpdateUser updates user details
func UpdateUser(userID uint, updatedData map[string]interface{}) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Check if user is deleted
	if user.IsDeleted {
		return errors.New("cannot update a deleted user")
	}

	// Validate email uniqueness if email is being updated
	if email, ok := updatedData["email"]; ok {
		var existingUser models.User
		if err := config.DB.Where("email = ? AND id != ? AND is_deleted = ?", email, userID, false).First(&existingUser).Error; err == nil {
			return errors.New("email already exists")
		}
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&user).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser performs soft delete (sets is_deleted=true)
func DeleteUser(userID uint) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	user.IsDeleted = true
	now := time.Now()
	user.UpdatedOn = &now

	if err := config.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// ChangePassword changes a user's password (requires old password verification)
func ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !VerifyPassword(user.Password, oldPassword) {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	now := time.Now()
	if err := config.DB.Model(&user).Updates(map[string]interface{}{
		"password":   hashedPassword,
		"updated_on": &now,
	}).Error; err != nil {
		return err
	}

	return nil
}

// ResetPassword resets a user's password (admin only, generates new password)
func ResetPassword(userID uint) (string, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return "", errors.New("user not found")
	}

	plainPassword := generateRandomPassword()
	hashedPassword, err := HashPassword(plainPassword)
	if err != nil {
		return "", err
	}

	// Update password
	now := time.Now()
	if err := config.DB.Model(&user).Updates(map[string]interface{}{
		"password":   hashedPassword,
		"updated_on": &now,
	}).Error; err != nil {
		return "", err
	}

	return plainPassword, nil
}
