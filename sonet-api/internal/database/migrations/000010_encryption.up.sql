-- Пользователь: ключи шифрования
ALTER TABLE users ADD COLUMN user_salt VARCHAR(64) DEFAULT NULL;
ALTER TABLE users ADD COLUMN encrypted_private_key TEXT DEFAULT NULL;
ALTER TABLE users ADD COLUMN public_key TEXT DEFAULT NULL;
ALTER TABLE users ADD COLUMN recovery_dek TEXT DEFAULT NULL;

-- Workspace members: зашифрованный DEK для каждого участника
ALTER TABLE workspace_members ADD COLUMN encrypted_dek TEXT DEFAULT NULL;

-- Workspace: флаг шифрования
ALTER TABLE workspaces ADD COLUMN is_encrypted BOOLEAN NOT NULL DEFAULT FALSE;

-- Заметки: зашифрованные поля
ALTER TABLE notes ADD COLUMN content_encrypted TEXT DEFAULT NULL;
ALTER TABLE notes ADD COLUMN content_iv VARCHAR(32) DEFAULT NULL;
ALTER TABLE notes ADD COLUMN title_encrypted TEXT DEFAULT NULL;
ALTER TABLE notes ADD COLUMN title_iv VARCHAR(32) DEFAULT NULL;
ALTER TABLE notes ADD COLUMN is_encrypted BOOLEAN NOT NULL DEFAULT FALSE;

-- Сессии
CREATE TABLE user_sessions (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    refresh_token_id BIGINT UNSIGNED DEFAULT NULL,
    device_name     VARCHAR(255) DEFAULT NULL,
    device_type     ENUM('desktop','mobile','tablet','unknown') DEFAULT 'unknown',
    os              VARCHAR(64) DEFAULT NULL,
    browser         VARCHAR(64) DEFAULT NULL,
    ip_address      VARCHAR(45) DEFAULT NULL,
    location        VARCHAR(128) DEFAULT NULL,
    last_active_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current      BOOLEAN NOT NULL DEFAULT FALSE,
    INDEX idx_us_user (user_id),
    CONSTRAINT fk_us_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
