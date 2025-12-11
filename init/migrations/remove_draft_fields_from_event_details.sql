-- Migration: Remove draft JSONB fields from event_details table
-- Date: 2024
-- Description: Remove draft fields since we're now using separate event_drafts table

-- Drop indexes if they exist
DROP INDEX IF EXISTS idx_event_details_general_details_draft;
DROP INDEX IF EXISTS idx_event_details_media_promotion_draft;
DROP INDEX IF EXISTS idx_event_details_special_guests_draft;
DROP INDEX IF EXISTS idx_event_details_volunteers_draft;

-- Remove JSONB columns from event_details table
ALTER TABLE event_details
DROP COLUMN IF EXISTS general_details_draft,
DROP COLUMN IF EXISTS media_promotion_draft,
DROP COLUMN IF EXISTS special_guests_draft,
DROP COLUMN IF EXISTS volunteers_draft;



