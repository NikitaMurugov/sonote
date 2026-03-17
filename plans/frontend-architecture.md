# Sonet — Архитектура фронтенда (Vue 3)

## Структура проекта

```
sonet-web/
├── index.html
├── vite.config.ts
├── tailwind.config.ts
├── tsconfig.json
├── package.json
├── public/
│   ├── favicon.svg
│   └── fonts/
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── router/
│   │   └── index.ts                     # vue-router конфигурация
│   ├── stores/                           # Pinia stores
│   │   ├── auth.ts                       # user, токены, login/logout
│   │   ├── workspace.ts                  # текущий workspace, список
│   │   ├── folder.ts                     # дерево папок
│   │   ├── note.ts                       # список заметок, текущая заметка, CRUD
│   │   ├── tag.ts                        # теги workspace
│   │   ├── search.ts                     # поисковый запрос + результаты
│   │   └── ui.ts                         # sidebar, тема, модалки
│   ├── composables/                      # переиспользуемая логика
│   │   ├── useApi.ts                     # axios с interceptors (refresh токена)
│   │   ├── useAuth.ts                    # обёртки login, register, logout
│   │   ├── useDebounce.ts
│   │   ├── useTheme.ts                   # light/dark, сохранение в localStorage
│   │   └── useGraph.ts                   # трансформация данных для d3
│   ├── api/                              # типизированные API-функции
│   │   ├── auth.ts
│   │   ├── workspaces.ts
│   │   ├── folders.ts
│   │   ├── notes.ts
│   │   ├── tags.ts
│   │   ├── search.ts
│   │   └── graph.ts
│   ├── types/                            # TypeScript интерфейсы
│   │   ├── user.ts
│   │   ├── workspace.ts
│   │   ├── folder.ts
│   │   ├── note.ts
│   │   ├── tag.ts
│   │   └── api.ts                        # ApiResponse<T>, PaginatedResponse
│   ├── views/
│   │   ├── auth/
│   │   │   ├── LoginView.vue
│   │   │   └── RegisterView.vue
│   │   ├── app/
│   │   │   ├── DashboardView.vue         # домашняя: недавние, закреплённые
│   │   │   ├── NoteEditorView.vue        # основной вид редактирования
│   │   │   ├── GraphView.vue             # полноэкранный граф
│   │   │   ├── SearchView.vue            # результаты поиска
│   │   │   ├── SettingsView.vue          # настройки workspace, участники
│   │   │   └── ProfileView.vue           # профиль пользователя
│   │   └── NotFoundView.vue
│   ├── components/
│   │   ├── layout/
│   │   │   ├── AppShell.vue              # sidebar + main content обёртка
│   │   │   ├── Sidebar.vue               # workspace switcher, папки, быстрые действия
│   │   │   ├── SidebarFolderTree.vue     # рекурсивное дерево папок
│   │   │   ├── SidebarNoteList.vue       # заметки в выбранной папке
│   │   │   └── TopBar.vue               # breadcrumb, поиск, аватар
│   │   ├── editor/
│   │   │   ├── NoteEditor.vue            # обёртка TipTap
│   │   │   ├── EditorToolbar.vue         # панель форматирования
│   │   │   ├── WikiLinkSuggestion.vue    # автокомплит для [[ ссылок
│   │   │   └── BacklinksPanel.vue        # панель входящих ссылок
│   │   ├── notes/
│   │   │   ├── NoteCard.vue              # карточка заметки
│   │   │   └── NoteMetaBar.vue           # теги, дата, кол-во слов
│   │   ├── tags/
│   │   │   ├── TagPill.vue               # цветной бейдж тега
│   │   │   ├── TagSelector.vue           # мульти-выбор тегов
│   │   │   └── TagManager.vue            # CRUD тегов workspace
│   │   ├── graph/
│   │   │   └── NoteGraph.vue             # d3-force граф
│   │   ├── workspace/
│   │   │   ├── WorkspaceSwitcher.vue      # переключатель workspace
│   │   │   ├── WorkspaceSettings.vue
│   │   │   └── MemberManager.vue          # приглашение, роли, удаление
│   │   ├── shared/
│   │   │   ├── SonetButton.vue
│   │   │   ├── SonetInput.vue
│   │   │   ├── SonetModal.vue
│   │   │   ├── SonetDropdown.vue
│   │   │   ├── SonetToast.vue
│   │   │   └── SonetAvatar.vue
│   │   └── auth/
│   │       ├── LoginForm.vue
│   │       └── RegisterForm.vue
│   ├── plugins/
│   │   └── tiptap.ts                     # регистрация расширений TipTap
│   ├── assets/
│   │   ├── styles/
│   │   │   ├── base.css                  # CSS custom properties для тем
│   │   │   └── tailwind.css              # @tailwind директивы
│   │   └── icons/
│   └── utils/
│       ├── date.ts
│       ├── slug.ts
│       └── markdown.ts
```

## Маршрутизация

```typescript
// Публичные (без авторизации)
/login                          → LoginView
/register                       → RegisterView

// Авторизованные (обёрнуты в AppShell)
/                               → redirect → /w/:defaultWorkspaceSlug
/w/:wsSlug                      → DashboardView
/w/:wsSlug/note/:noteSlug       → NoteEditorView
/w/:wsSlug/graph                → GraphView
/w/:wsSlug/search               → SearchView
/w/:wsSlug/settings             → SettingsView
/profile                        → ProfileView

// Catch-all
/:pathMatch(.*)*                → NotFoundView
```

## State Management (Pinia)

### authStore
- Состояние: `user`, `accessToken`, `refreshToken`
- Действия: `login()`, `register()`, `logout()`, `refreshTokens()`
- Персистенция: токены в `localStorage`

### workspaceStore
- Состояние: `workspaces[]`, `currentWorkspace`
- Загружается при инициализации после авторизации
- Смена workspace → перезагрузка folder/note/tag stores

### folderStore
- Состояние: `folderTree` (вложенный массив), `selectedFolderId`
- Загружается при смене workspace

### noteStore
- Состояние: `notes[]` (для текущей папки/вида), `currentNote` (полное содержание)
- Загружается при смене папки или при поиске

### tagStore
- Состояние: `tags[]` для текущего workspace
- Загружается один раз при смене workspace

### uiStore
- Состояние: `sidebarCollapsed`, `theme` ('light'|'dark'), `activeModal`
- Персистенция: `localStorage`

## API-клиент (useApi)
- Обёртка над axios
- Автоматический заголовок `Authorization: Bearer` из `authStore.accessToken`
- Interceptor на 401: вызов `authStore.refreshTokens()`, очередь отложенных запросов, retry
- Если refresh тоже 401 → redirect на `/login`
