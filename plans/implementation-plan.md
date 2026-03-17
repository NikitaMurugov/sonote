# Sonet — План реализации

## Фаза 1: Фундамент

### Backend
1. Инициализация Go-проекта (`go mod init`), структура директорий
2. `docker-compose.yml` с MariaDB
3. Конфиг (`internal/config`) — загрузка из env
4. Подключение к БД (`internal/database/mariadb.go`)
5. Миграции (9 файлов) — создание всех таблиц
6. Пакеты: `pkg/jwt`, `pkg/hash`, `pkg/response`, `pkg/validator`
7. `model/user.go`, `repository/user_repo.go`
8. `service/auth_service.go` — register, login, refresh, logout
9. `handler/auth_handler.go` — HTTP-эндпоинты авторизации
10. `middleware/auth.go` — JWT валидация
11. `middleware/cors.go` — CORS для фронтенда
12. `router/router.go` — сборка маршрутов

### Frontend
1. Инициализация Vue проекта (Vite + TypeScript)
2. Настройка Tailwind CSS + CSS custom properties (темы)
3. Установка шрифтов (Inter, Playfair Display, JetBrains Mono)
4. Базовые shared-компоненты: SonetButton, SonetInput, SonetModal
5. `stores/auth.ts` + `stores/ui.ts`
6. `composables/useApi.ts` — axios + interceptors
7. `composables/useTheme.ts` — light/dark
8. Страницы LoginView, RegisterView
9. Маршрутизация + navigation guards

## Фаза 2: Основной CRUD

### Backend
1. Workspace: model, repo, service, handler
2. Folder: model, repo, service, handler (дерево)
3. Note: model, repo, service, handler (CRUD + фильтрация)
4. `middleware/workspace_access.go` — проверка ролей

### Frontend
1. `AppShell.vue` — sidebar + main content layout
2. `Sidebar.vue`, `WorkspaceSwitcher.vue`
3. `SidebarFolderTree.vue` — рекурсивное дерево
4. `SidebarNoteList.vue` — список заметок
5. `TopBar.vue` — breadcrumb, поиск, аватар
6. `DashboardView.vue` — домашняя workspace
7. `NoteEditorView.vue` — интеграция TipTap
8. `NoteEditor.vue`, `EditorToolbar.vue` — базовый редактор
9. `NoteCard.vue`, `NoteMetaBar.vue`
10. `stores/workspace.ts`, `stores/folder.ts`, `stores/note.ts`

## Фаза 3: Расширенные фичи

### Backend
1. Tags: model, repo, service, handler
2. `pkg/linkparser/parser.go` — извлечение `[[wiki-links]]`
3. `service/link_service.go` — при сохранении заметки парсить ссылки, diff с note_links
4. Эндпоинты links/backlinks для заметок
5. `service/search_service.go` + `handler/search_handler.go` — FULLTEXT поиск

### Frontend
1. `TagPill.vue`, `TagSelector.vue`, `TagManager.vue`
2. `stores/tag.ts`
3. TipTap расширение для `[[wiki-links]]` — suggestion plugin
4. `WikiLinkSuggestion.vue` — автокомплит
5. `BacklinksPanel.vue` — панель backlinks
6. `SearchView.vue` + `stores/search.ts`

## Фаза 4: Коллаборация и граф

### Backend
1. Workspace members: handler для invite, role change, remove
2. `handler/graph_handler.go` — возврат nodes/edges для графа

### Frontend
1. `MemberManager.vue` — приглашение, управление ролями
2. `SettingsView.vue` — вкладки General/Members/Tags
3. `NoteGraph.vue` — d3-force визуализация
4. `GraphView.vue` — полноэкранный граф
5. `ProfileView.vue` — настройки профиля
6. Тёмная тема — проверка и доводка всех компонентов
7. Адаптивный дизайн: сворачивание sidebar
8. Toast-уведомления, loading states, error handling

## Фаза 5: Нативные приложения (будущее)

1. Tauri обёртка → macOS desktop app
2. Capacitor обёртка → iOS app
3. Общий REST API, без изменений в бэкенде

---

## Порядок запуска разработки

```
1. docker-compose up -d        # MariaDB
2. cd sonet-api && go run ./cmd/sonet   # Backend на :8080
3. cd sonet-web && npm run dev           # Frontend на :5173
```

## Верификация

- **Auth**: регистрация → логин → получение токенов → refresh → доступ к защищённым роутам
- **CRUD**: создание workspace → папка → заметка → редактирование → удаление
- **Links**: вставка `[[ссылки]]` → проверка note_links в БД → отображение backlinks
- **Search**: полнотекстовый поиск по заметкам
- **Sharing**: приглашение пользователя → проверка доступа по ролям
- **Graph**: визуализация связей, клик → навигация
- **Темы**: переключение light/dark, проверка всех компонентов
