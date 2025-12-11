-- Migration: Add draft JSONB fields to event_details table
-- Date: 2024
-- Description: Adds JSONB fields for auto-saving draft data for each step

-- Add JSONB columns for draft data
ALTER TABLE event_details
ADD COLUMN IF NOT EXISTS general_details_draft JSONB DEFAULT '{}'::jsonb,
ADD COLUMN IF NOT EXISTS media_promotion_draft JSONB DEFAULT '{}'::jsonb,
ADD COLUMN IF NOT EXISTS special_guests_draft JSONB DEFAULT '{}'::jsonb,
ADD COLUMN IF NOT EXISTS volunteers_draft JSONB DEFAULT '{}'::jsonb;

-- Create indexes on JSONB fields for better query performance (optional)
CREATE INDEX IF NOT EXISTS idx_event_details_general_details_draft ON event_details USING GIN (general_details_draft);
CREATE INDEX IF NOT EXISTS idx_event_details_media_promotion_draft ON event_details USING GIN (media_promotion_draft);
CREATE INDEX IF NOT EXISTS idx_event_details_special_guests_draft ON event_details USING GIN (special_guests_draft);
CREATE INDEX IF NOT EXISTS idx_event_details_volunteers_draft ON event_details USING GIN (volunteers_draft);



