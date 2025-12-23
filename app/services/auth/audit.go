package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/google/uuid"
)

// AuditEventType represents the type of audit event
type AuditEventType string

const (
	AuditEventLogin            AuditEventType = "login"
	AuditEventLoginFailed      AuditEventType = "login_failed"
	AuditEventLogout           AuditEventType = "logout"
	AuditEventRegister         AuditEventType = "register"
	AuditEventEmailVerified    AuditEventType = "email_verified"
	AuditEventPasswordReset    AuditEventType = "password_reset"
	AuditEventPasswordChanged  AuditEventType = "password_changed"
	AuditEventSessionRevoked   AuditEventType = "session_revoked"
	AuditEventTokenRefreshed   AuditEventType = "token_refreshed"
)

// LogAuditEvent logs an authentication event for security auditing
// Never logs sensitive data (passwords, tokens)
func LogAuditEvent(ctx context.Context, eventType AuditEventType, userID *int64, ip, userAgent string, metadata map[string]interface{}) error {
	var metadataJSON []byte
	if metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	var userIDVal interface{}
	if userID != nil {
		userIDVal = *userID
	}

	eventID := uuid.New().String()
	query := `
		INSERT INTO auth_audit_events (id, user_id, type, ip, user_agent, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err := config.AuthDB.Exec(ctx, query, eventID, userIDVal, string(eventType), ip, userAgent, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to insert audit event: %w", err)
	}

	return nil
}

