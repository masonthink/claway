-- 001_init.down.sql
-- Rollback initial schema

DROP TABLE IF EXISTS credit_transactions;
DROP TABLE IF EXISTS contributions;
DROP TABLE IF EXISTS token_usage_logs;
DROP TABLE IF EXISTS document_versions;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS prds;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS ideas;
DROP TABLE IF EXISTS users;
