-- Migration: Create user_preferences table
-- Description: Stores user-specific preferences like column visibility, pinned columns, etc.

CREATE TABLE IF NOT EXISTS user_preferences (
    id BIGSERIAL PRIMARY KEY,
    user_email VARCHAR(255) NOT NULL,
    preference_type VARCHAR(100) NOT NULL,
    preference_data JSONB NOT NULL,
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_preference UNIQUE (user_email, preference_type)
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_email ON user_preferences(user_email);
CREATE INDEX IF NOT EXISTS idx_user_preferences_type ON user_preferences(preference_type);
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_type ON user_preferences(user_email, preference_type);

-- Add comment to table
COMMENT ON TABLE user_preferences IS 'Stores user-specific preferences like column visibility, pinned columns, dashboard layouts, etc.';
COMMENT ON COLUMN user_preferences.user_email IS 'Email of the user who owns these preferences';
COMMENT ON COLUMN user_preferences.preference_type IS 'Type of preference (e.g., events_list_columns, dashboard_layout)';
COMMENT ON COLUMN user_preferences.preference_data IS 'JSONB data containing the actual preference values';

