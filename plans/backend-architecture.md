# Sonet — Архитектура бэкенда (Go)

## Структура проекта

```
sonet-api/
├── cmd/
│   └── sonet/
│       └── main.go                  # точка входа, загрузка конфига, запуск сервера
├── internal/
│   ├── config/
│   │   └── config.go                # структура конфига + загрузка из env/yaml
│   ├── database/
│   │   ├── mariadb.go               # пул соединений
│   │   └── migrations/
│   │       ├── 001_create_users.sql
│   │       ├── 002_create_refresh_tokens.sql
│   │       ├── 003_create_workspaces.sql
│   │       ├── 004_create_workspace_members.sql
│   │       ├── 005_create_folders.sql
│   │       ├── 006_create_notes.sql
│   │       ├── 007_create_note_links.sql
│   │       ├── 008_create_tags.sql
│   │       └── 009_create_note_tags.sql
│   ├── model/                        # структуры данных
│   │   ├── user.go
│   │   ├── workspace.go
│   │   ├── folder.go
│   │   ├── note.go
│   │   ├── note_link.go
│   │   ├── tag.go
│   │   └── refresh_token.go
│   ├── repository/                   # слой доступа к данным (sqlx)
│   │   ├── user_repo.go
│   │   ├── workspace_repo.go
│   │   ├── folder_repo.go
│   │   ├── note_repo.go
│   │   ├── note_link_repo.go
│   │   ├── tag_repo.go
│   │   └── token_repo.go
│   ├── service/                      # бизнес-логика
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── workspace_service.go
│   │   ├── folder_service.go
│   │   ├── note_service.go
│   │   ├── link_service.go
│   │   ├── tag_service.go
│   │   └── search_service.go
│   ├── handler/                      # HTTP-обработчики
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── workspace_handler.go
│   │   ├── folder_handler.go
│   │   ├── note_handler.go
│   │   ├── tag_handler.go
│   │   ├── search_handler.go
│   │   └── graph_handler.go
│   ├── middleware/
│   │   ├── auth.go                   # валидация JWT, инъекция user в context
│   │   ├── cors.go
│   │   ├── rate_limit.go
│   │   ├── logging.go                # структурированное логирование запросов
│   │   └── workspace_access.go       # проверка роли в workspace
│   ├── router/
│   │   └── router.go                 # регистрация маршрутов
│   └── pkg/
│       ├── jwt/
│       │   └── jwt.go                # генерация + валидация токенов
│       ├── hash/
│       │   └── password.go           # bcrypt обёртка
│       ├── validator/
│       │   └── validator.go          # валидация входных данных
│       ├── linkparser/
│       │   └── parser.go             # извлечение [[wiki-links]] из markdown
│       └── response/
│           └── response.go           # стандартные JSON-ответы
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── docker-compose.yml
```

## Библиотеки

| Назначение | Библиотека |
|-----------|-----------|
| HTTP роутер | `chi` (go-chi/chi/v5) |
| SQL | `sqlx` + `go-sql-driver/mysql` |
| JWT | `golang-jwt/jwt/v5` |
| Хэширование | `golang.org/x/crypto/bcrypt` |
| Валидация | `go-playground/validator/v10` |
| Логирование | `slog` (stdlib Go 1.21+) |
| Конфиг | `caarlos0/env` |
| Миграции | `golang-migrate/migrate/v4` |
| Rate limiting | `go-chi/httprate` |

## Цепочка middleware

```
Request
  → CORS
  → Request ID (X-Request-ID)
  → Structured Logger (slog)
  → Rate Limiter (100 req/min публичные, 300 для авторизованных)
  → [Защищённые роуты] JWT Auth → user_id в context
  → [Workspace роуты] Workspace Access Check (роль >= требуемой)
  → Handler
```

## API эндпоинты (префикс `/api/v1`)

### Auth
```
POST   /auth/register              # { email, username, password, display_name }
POST   /auth/login                 # { email, password } → { access_token, refresh_token }
POST   /auth/refresh               # { refresh_token } → { access_token, refresh_token }
POST   /auth/logout                # отзыв refresh token
```

### Users
```
GET    /users/me                   # профиль текущего пользователя
PATCH  /users/me                   # обновить display_name, avatar_url
PUT    /users/me/password          # сменить пароль
```

### Workspaces
```
GET    /workspaces                          # список workspace пользователя
POST   /workspaces                          # создать workspace
GET    /workspaces/:wsId                    # детали workspace
PATCH  /workspaces/:wsId                    # обновить (admin)
DELETE /workspaces/:wsId                    # удалить (owner)
GET    /workspaces/:wsId/members            # список участников
POST   /workspaces/:wsId/members            # пригласить { email, role }
PATCH  /workspaces/:wsId/members/:userId    # изменить роль (admin)
DELETE /workspaces/:wsId/members/:userId    # удалить участника (admin)
```

### Folders
```
GET    /workspaces/:wsId/folders            # дерево папок
POST   /workspaces/:wsId/folders            # создать папку
PATCH  /workspaces/:wsId/folders/:fId       # переименовать / переместить
DELETE /workspaces/:wsId/folders/:fId       # удалить (заметки → родитель или корень)
```

### Notes
```
GET    /workspaces/:wsId/notes              # список (фильтры: folder_id, tag, archived, pinned)
POST   /workspaces/:wsId/notes              # создать заметку
GET    /workspaces/:wsId/notes/:nId         # полное содержание заметки
PATCH  /workspaces/:wsId/notes/:nId         # обновить
DELETE /workspaces/:wsId/notes/:nId         # удалить
GET    /workspaces/:wsId/notes/:nId/links      # исходящие ссылки
GET    /workspaces/:wsId/notes/:nId/backlinks  # входящие ссылки (backlinks)
```

### Tags
```
GET    /workspaces/:wsId/tags               # все теги workspace
POST   /workspaces/:wsId/tags               # создать тег
PATCH  /workspaces/:wsId/tags/:tId          # переименовать / изменить цвет
DELETE /workspaces/:wsId/tags/:tId          # удалить тег
POST   /workspaces/:wsId/notes/:nId/tags    # привязать теги { tag_ids: [] }
DELETE /workspaces/:wsId/notes/:nId/tags/:tId  # отвязать тег
```

### Search
```
GET    /workspaces/:wsId/search?q=...&tags=...&folder=...  # полнотекстовый поиск
```

### Graph
```
GET    /workspaces/:wsId/graph              # { nodes: [...], edges: [...] }
```

## Поток авторизации

1. **Register** → валидация, bcrypt (cost 12), создание user, автоматическое создание personal workspace, выдача токенов
2. **Login** → поиск по email, проверка bcrypt, генерация access token (JWT, 15 мин, содержит user_id + email), генерация refresh token (256-bit random, SHA-256 хэш в БД, 30 дней), возврат обоих
3. **Refresh** → принять refresh token, SHA-256 хэш, поиск в БД, проверка revoked/expired, отзыв старого, выпуск новых (ротация токенов)
4. **Logout** → отзыв refresh token в БД, клиент удаляет access token

## Правила авторизации по ролям

| Ресурс | Viewer | Editor | Admin | Owner |
|--------|--------|--------|-------|-------|
| Чтение заметок | ✓ | ✓ | ✓ | ✓ |
| Создание/редактирование | ✗ | ✓ | ✓ | ✓ |
| Управление папками | ✗ | ✓ | ✓ | ✓ |
| Управление тегами | ✗ | ✓ | ✓ | ✓ |
| Приглашение участников | ✗ | ✗ | ✓ | ✓ |
| Изменение ролей | ✗ | ✗ | ✓ | ✓ |
| Настройки workspace | ✗ | ✗ | ✓ | ✓ |
| Удаление workspace | ✗ | ✗ | ✗ | ✓ |
