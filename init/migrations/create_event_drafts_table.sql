-- Migration: Create event_drafts table for storing draft data
-- Date: 2024
-- Description: Separate table for event draft data that will be deleted after submission

-- Create event_drafts table
CREATE TABLE IF NOT EXISTS event_drafts (
    id BIGSERIAL PRIMARY KEY,

    -- Draft data for each step (JSONB)
    general_details_draft JSONB DEFAULT '{}'::jsonb,
    media_promotion_draft JSONB DEFAULT '{}'::jsonb,
    special_guests_draft JSONB DEFAULT '{}'::jsonb,
    volunteers_draft JSONB DEFAULT '{}'::jsonb,

    -- Optional: Link to event if draft is associated with an existing event
    event_id BIGINT REFERENCES event_details(id) ON DELETE CASCADE,

    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMPTZ
);

-- Create indexes on JSONB fields for better query performance
CREATE INDEX IF NOT EXISTS idx_event_drafts_general_details_draft ON event_drafts USING GIN (general_details_draft);
CREATE INDEX IF NOT EXISTS idx_event_drafts_media_promotion_draft ON event_drafts USING GIN (media_promotion_draft);
CREATE INDEX IF NOT EXISTS idx_event_drafts_special_guests_draft ON event_drafts USING GIN (special_guests_draft);
CREATE INDEX IF NOT EXISTS idx_event_drafts_volunteers_draft ON event_drafts USING GIN (volunteers_draft);

-- Create index on event_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_event_drafts_event_id ON event_drafts(event_id);

-- Create index on created_on for cleanup of old drafts
CREATE INDEX IF NOT EXISTS idx_event_drafts_created_on ON event_drafts(created_on);



