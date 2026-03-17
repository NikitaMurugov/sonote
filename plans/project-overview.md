# Sonet — Обзор проекта

## Концепция

**Sonet** — легковесное веб-приложение для личных заметок, вдохновлённое Obsidian. Название отсылает к сонету — лёгкой, простой и функциональной форме.

### Ключевые возможности
- Markdown-заметки с rich-text редактором (TipTap)
- Bi-directional links (`[[ссылки между заметками]]`)
- Теги для гибкой группировки
- Иерархические папки
- Граф связей между заметками
- Рабочие пространства (Workspaces) с шарингом и ролями
- Полнотекстовый поиск
- Светлая и тёмная темы

### Будущее
- Нативное macOS приложение через Tauri
- Нативное iOS приложение через Capacitor
- Общая REST API для всех клиентов

## Стек технологий

| Слой | Технология |
|------|-----------|
| Backend | Go, Chi (router), sqlx, MariaDB |
| Auth | JWT (access 15мин + refresh 30 дней), bcrypt |
| Frontend | Vue 3, TypeScript, Vite, Tailwind CSS |
| Редактор | TipTap (ProseMirror) |
| Граф | d3-force (SVG) |
| DevOps | Docker Compose (MariaDB + API) |
