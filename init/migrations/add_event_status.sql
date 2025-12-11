-- Add status column to event_details table
ALTER TABLE event_details
ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'incomplete' CHECK (status IN ('complete', 'incomplete'));

-- Update existing events to have 'incomplete' status if they don't have one
UPDATE event_details
SET status = 'incomplete'
WHERE status IS NULL;

-- Add index on status for faster filtering
CREATE INDEX IF NOT EXISTS idx_event_details_status ON event_details(status);



