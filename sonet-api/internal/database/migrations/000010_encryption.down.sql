DROP TABLE IF EXISTS user_sessions;

ALTER TABLE notes DROP COLUMN is_encrypted;
ALTER TABLE notes DROP COLUMN title_iv;
ALTER TABLE notes DROP COLUMN title_encrypted;
ALTER TABLE notes DROP COLUMN content_iv;
ALTER TABLE notes DROP COLUMN content_encrypted;

ALTER TABLE workspaces DROP COLUMN is_encrypted;

ALTER TABLE workspace_members DROP COLUMN encrypted_dek;

ALTER TABLE users DROP COLUMN recovery_dek;
ALTER TABLE users DROP COLUMN public_key;
ALTER TABLE users DROP COLUMN encrypted_private_key;
ALTER TABLE users DROP COLUMN user_salt;
