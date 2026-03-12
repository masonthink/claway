-- Rename task types from D1-D4 to doc1-doc4
-- Must drop constraint first, otherwise UPDATEs fail
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_type_check;

UPDATE tasks SET type = 'doc1' WHERE type = 'D1';
UPDATE tasks SET type = 'doc2' WHERE type = 'D2';
UPDATE tasks SET type = 'doc3' WHERE type = 'D3';
UPDATE tasks SET type = 'doc4' WHERE type = 'D4';

-- Update dependencies references
UPDATE tasks SET dependencies = REPLACE(REPLACE(REPLACE(REPLACE(dependencies, 'D1', 'doc1'), 'D2', 'doc2'), 'D3', 'doc3'), 'D4', 'doc4')
WHERE dependencies != '';

-- Add new check constraint
ALTER TABLE tasks ADD CONSTRAINT tasks_type_check CHECK (type IN ('doc1','doc2','doc3','doc4'));
