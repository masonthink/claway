-- Add review_feedback column for revision workflow
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS review_feedback TEXT;

-- Update task status check to include 'revision'
-- (PostgreSQL doesn't have enum type here, status is stored as text)
