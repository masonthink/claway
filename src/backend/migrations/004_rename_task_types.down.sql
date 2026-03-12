-- Revert task types from doc1-doc4 back to D1-D4
UPDATE tasks SET type = 'D1' WHERE type = 'doc1';
UPDATE tasks SET type = 'D2' WHERE type = 'doc2';
UPDATE tasks SET type = 'D3' WHERE type = 'doc3';
UPDATE tasks SET type = 'D4' WHERE type = 'doc4';

UPDATE tasks SET dependencies = REPLACE(REPLACE(REPLACE(REPLACE(dependencies, 'doc1', 'D1'), 'doc2', 'D2'), 'doc3', 'D3'), 'doc4', 'D4')
WHERE dependencies != '';

ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_type_check;
ALTER TABLE tasks ADD CONSTRAINT tasks_type_check CHECK (type IN ('D1','D2','D3','D4'));
