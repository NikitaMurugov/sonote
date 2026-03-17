CREATE TABLE note_links (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    source_note_id  BIGINT UNSIGNED NOT NULL,
    target_note_id  BIGINT UNSIGNED NOT NULL,
    context_snippet VARCHAR(512)    DEFAULT NULL,
    created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_link (source_note_id, target_note_id),
    INDEX idx_nl_target (target_note_id),
    CONSTRAINT fk_nl_source FOREIGN KEY (source_note_id) REFERENCES notes(id) ON DELETE CASCADE,
    CONSTRAINT fk_nl_target FOREIGN KEY (target_note_id) REFERENCES notes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
