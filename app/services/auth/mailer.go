package auth

// Mailer interface for sending authentication-related emails
type Mailer interface {
	// SendVerification sends an email verification token to the user
	SendVerification(email, token string) error

	// SendPasswordReset sends a password reset token to the user
	SendPasswordReset(email, token string) error
}

// StubMailer is a stub implementation that logs instead of sending emails
// Replace this with your actual email service (SendGrid, SES, etc.)
type StubMailer struct{}

func NewStubMailer() *StubMailer {
	return &StubMailer{}
}

func (m *StubMailer) SendVerification(email, token string) error {
	// In production, replace this with actual email sending
	// NEVER log the token in production logs
	// fmt.Printf("[STUB] Would send verification email to %s with token: %s\n", email, token)
	
	// For security, we don't log the token at all
	// fmt.Printf("[STUB] Would send verification email to %s\n", email)
	return nil
}

func (m *StubMailer) SendPasswordReset(email, token string) error {
	// In production, replace this with actual email sending
	// NEVER log the token in production logs
	// fmt.Printf("[STUB] Would send password reset email to %s with token: %s\n", email, token)
	
	// For security, we don't log the token at all
	// fmt.Printf("[STUB] Would send password reset email to %s\n", email)
	return nil
}



