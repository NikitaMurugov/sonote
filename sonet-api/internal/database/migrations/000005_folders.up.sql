CREATE TABLE folders (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    workspace_id    BIGINT UNSIGNED NOT NULL,
    parent_id       BIGINT UNSIGNED DEFAULT NULL,
    name            VARCHAR(255)    NOT NULL,
    sort_order      INT             NOT NULL DEFAULT 0,
    created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_f_ws (workspace_id),
    INDEX idx_f_parent (parent_id),
    UNIQUE KEY uk_folder_name (workspace_id, parent_id, name),
    CONSTRAINT fk_f_ws     FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
    CONSTRAINT fk_f_parent FOREIGN KEY (parent_id)    REFERENCES folders(id)    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
