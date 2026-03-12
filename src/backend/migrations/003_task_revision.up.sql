-- Add review_feedback column for revision workflow
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS review_feedback TEXT;

-- Update task status check to include 'revision'
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_status_check;
ALTER TABLE tasks ADD CONSTRAINT tasks_status_check CHECK (status IN ('open','claimed','submitted','approved','rejected','revision'));
