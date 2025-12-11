# main - Основной пакет приложения

Этот репозиторий представляет собой шаблон (boilerplate) для создания серверных приложений на Go. Он использует веб-фреймворк Fiber, работу с базой данных через `sqlx` и управление миграциями с помощью `golang-migrate`.
https://github.com/GlebSr/ServerTemplate
## Особенности

- **REST API** с аутентификацией JWT
- **WebSocket** для обновлений в реальном времени
- **PostgreSQL** база данных с миграциями
- **Централизованный Hub** для broadcast-событий по комнатам

## WebSocket Real-Time

Приложение поддерживает WebSocket-соединения для мониторинга изменений в комнатах в реальном времени:

- Подключение: `ws://localhost:8080/api/v1/rooms/{room_id}/ws`
- События: добавление/удаление игр, голосов, участников, обновление комнаты
- Формат: JSON с типом события, payload, и timestamp
- Auth: требуется Bearer token в заголовках при upgrade

Подробнее в `API.md` → раздел **WebSocket Real-Time Updates**.

## Установка и запуск

# Задаём миграции
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Создание новых миграций
migrate create -ext sql -dir migrations -seq migration_name