package models

// CreateUserResponse represents the response when creating a user
// swagger:model CreateUserResponse
type CreateUserResponse struct {
	Message  string `json:"message"`
	User     User   `json:"user"`
	Password string `json:"password"`
}

// ResetPasswordResponse represents the response when resetting a user's password
// swagger:model ResetPasswordResponse
type ResetPasswordResponse struct {
	Message  string `json:"message"`
	Password string `json:"password"`
}



