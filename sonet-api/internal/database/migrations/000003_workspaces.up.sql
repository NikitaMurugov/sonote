CREATE TABLE workspaces (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(128)    NOT NULL,
    slug            VARCHAR(128)    NOT NULL UNIQUE,
    description     TEXT            DEFAULT NULL,
    owner_id        BIGINT UNSIGNED NOT NULL,
    is_personal     BOOLEAN         NOT NULL DEFAULT FALSE,
    icon            VARCHAR(32)     DEFAULT NULL,
    created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ws_owner (owner_id),
    INDEX idx_ws_slug (slug),
    CONSTRAINT fk_ws_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
