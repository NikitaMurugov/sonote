# Sonet — Схема базы данных (MariaDB)

## Диаграмма связей

```
User ──owns──> Workspace (is_personal=true для дефолтного)
User ──member──> Workspace (через workspace_members с ролью)

Workspace ──has──> Folder[] (иерархия через parent_id)
Workspace ──has──> Note[] (могут быть без папки)
Workspace ──has──> Tag[] (теги привязаны к workspace)

Folder ──contains──> Note[]
Folder ──parent──> Folder (вложенность)

Note ──links──> Note[] (через note_links, source → target)
Note ──tagged──> Tag[] (через note_tags)
Note ──authored──> User
```

## Таблицы

### users
```sql
CREATE TABLE users (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email           VARCHAR(255)    NOT NULL UNIQUE,
    username        VARCHAR(64)     NOT NULL UNIQUE,
    display_name    VARCHAR(128)    NOT NULL DEFAULT '',
    password_hash   VARCHAR(255)    NOT NULL,
    avatar_url      VARCHAR(512)    DEFAULT NULL,
    email_verified  BOOLEAN         NOT NULL DEFAULT FALSE,
    created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_users_email (email),
    INDEX idx_users_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### refresh_tokens
```sql
CREATE TABLE refresh_tokens (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    token_hash      VARCHAR(255)    NOT NULL UNIQUE,
    device_info     VARCHAR(255)    DEFAULT NULL,
    ip_address      VARCHAR(45)     DEFAULT NULL,
    expires_at      DATETIME        NOT NULL,
    revoked         BOOLEAN         NOT NULL DEFAULT FALSE,
    created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_rt_user (user_id),
    INDEX idx_rt_token (token_hash),
    INDEX idx_rt_expires (expires_at),
    CONSTRAINT fk_rt_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### workspaces
```sql
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
```

### workspace_members
```sql
CREATE TABLE workspace_members (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    workspace_id    BIGINT UNSIGNED NOT NULL,
    user_id         BIGINT UNSIGNED NOT NULL,
    role            ENUM('viewer','editor','admin') NOT NULL DEFAULT 'viewer',
    invited_by      BIGINT UNSIGNED DEFAULT NULL,
    joined_at       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_ws_member (workspace_id, user_id),
    INDEX idx_wm_user (user_id),
    CONSTRAINT fk_wm_ws   FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
    CONSTRAINT fk_wm_user FOREIGN KEY (user_id)      REFERENCES users(id)      ON DELETE CASCADE,
    CONSTRAINT fk_wm_inv  FOREIGN KEY (invited_by)    REFERENCES users(id)      ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### folders
```sql
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
```

### notes
```sql
CREATE TABLE notes (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    workspace_id    BIGINT UNSIGNED NOT NULL,
    folder_id       BIGINT UNSIGNED DEFAULT NULL,
    title           VARCHAR(512)    NOT NULL,
    slug            VARCHAR(512)    NOT NULL,
    content_md      MEDIUMTEXT      NOT NULL DEFAULT '',
    content_html    MEDIUMTEXT      NOT NULL DEFAULT '',
    content_json    JSON            DEFAULT NULL,
    author_id       BIGINT UNSIGNED NOT NULL,
    is_pinned       BOOLEAN         NOT NULL DEFAULT FALSE,
    is_archived     BOOLEAN         NOT NULL DEFAULT FALSE,
    word_count      INT UNSIGNED    NOT NULL DEFAULT 0,
    created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_n_ws (workspace_id),
    INDEX idx_n_folder (folder_id),
    INDEX idx_n_author (author_id),
    INDEX idx_n_slug (workspace_id, slug),
    INDEX idx_n_updated (updated_at),
    FULLTEXT INDEX ft_n_search (title, content_md),
    CONSTRAINT fk_n_ws     FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
    CONSTRAINT fk_n_folder FOREIGN KEY (folder_id)    REFERENCES folders(id)    ON DELETE SET NULL,
    CONSTRAINT fk_n_author FOREIGN KEY (author_id)    REFERENCES users(id)      ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### note_links
```sql
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
```

### tags
```sql
CREATE TABLE tags (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    workspace_id    BIGINT UNSIGNED NOT NULL,
    name            VARCHAR(128)    NOT NULL,
    color           VARCHAR(7)      DEFAULT NULL,
    created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_tag (workspace_id, name),
    CONSTRAINT fk_t_ws FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### note_tags
```sql
CREATE TABLE note_tags (
    note_id         BIGINT UNSIGNED NOT NULL,
    tag_id          BIGINT UNSIGNED NOT NULL,
    created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (note_id, tag_id),
    INDEX idx_nt_tag (tag_id),
    CONSTRAINT fk_nt_note FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
    CONSTRAINT fk_nt_tag  FOREIGN KEY (tag_id)  REFERENCES tags(id)  ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## Ключевые решения
- **BIGINT UNSIGNED** для всех ID — запас на рост
- **content_json** хранит TipTap ProseMirror JSON, **content_md** — markdown для поиска и экспорта
- **FULLTEXT INDEX** на `notes(title, content_md)` для полнотекстового поиска
- **note_links** — направленные связи: source → target. Backlinks — обратный запрос по target_note_id
- **Теги привязаны к workspace** — имена уникальны в рамках workspace
- **slug** на заметках для разрешения `[[wiki-links]]`
