# API Documentation

## Базовый URL
```
/
```

---

## Публичные эндпоинты

### 1. Регистрация пользователя
**POST** `/api/auth/signup`

Создает новый аккаунт пользователя.

**Request Body:**
```json
{
  "name": "string",
  "password": "string"
}
```

**Response (201 Created):**
```json
{
  "user_id": "uuid",
  "access_token": "jwt_token"
}
```

**Cookies:**
- `refresh_token` - HTTP-only cookie для обновления токена

**Errors:**
- `400` - Неверный формат запроса
- `409` - Пользователь с таким именем уже существует
- `500` - Внутренняя ошибка сервера

---

### 2. Вход в систему
**POST** `/api/auth/signin`

Авторизует пользователя и возвращает токены доступа.

**Request Body:**
```json
{
  "name": "string",
  "password": "string"
}
```

**Response (200 OK):**
```json
{
  "user_id": "uuid",
  "access_token": "jwt_token"
}
```

**Cookies:**
- `refresh_token` - HTTP-only cookie для обновления токена

**Errors:**
- `400` - Неверный формат запроса
- `401` - Неверное имя пользователя или пароль
- `500` - Внутренняя ошибка сервера

---

### 3. Обновление токена доступа
**POST** `/api/auth/refresh`

Обновляет access token используя refresh token из cookies.

**Request:**
- Требуется cookie `refresh_token`

**Response (200 OK):**
```json
{
  "access_token": "jwt_token"
}
```

**Cookies:**
- `refresh_token` - Новый HTTP-only cookie

**Errors:**
- `401` - Неверный или истекший refresh token
- `500` - Внутренняя ошибка сервера

---

#### 4. Выход из системы
**POST** `/api/auth/logout`

Удаляет refresh token пользователя из базы данных.

**Request:**
- Требуется cookie `refresh_token`

**Response (204 No Content)**

**Errors:**
- `401` - Не авторизован (невалидный или отсутствующий refresh token)
- `500` - Внутренняя ошибка сервера

---

## Аутентифицированные эндпоинты

> **Требуется:** Header `Authorization: Bearer <access_token>`

### Управление аккаунтом

#### 5. Получить информацию о текущем пользователе
**GET** `/api/v1/user`

Возвращает информацию о текущем авторизованном пользователе.

**Response (200 OK):**
```json
{
  "id": "uuid",
  "name": "string"
}
```

**Errors:**
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

---

#### 6. Поиск пользователя по имени
**GET** `/api/v1/user/by-name?name=username`

Находит пользователя по имени.

**Query Parameters:**
- `name` (string, required) - Имя пользователя для поиска

**Response (200 OK):**
```json
{
  "id": "uuid",
  "name": "string"
}
```

**Errors:**
- `401` - Не авторизован
- `404` - Пользователь не найден
- `500` - Внутренняя ошибка сервера

---

#### 7. Обновить данные пользователя
**PUT** `/api/v1/user`

Обновляет имя и/или пароль текущего пользователя.

**Request Body:**
```json
{
  "name": "string (optional)",
  "password": "string (optional)"
}
```

**Response (200 OK)**

**Errors:**
- `400` - Неверный формат запроса
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

---

### Управление комнатами

#### 8. Создать комнату
**POST** `/api/v1/rooms`

Создает новую комнату для голосования.

**Request Body:**
```json
{
  "name": "string"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "name": "string",
  "owner_id": "uuid"
}
```

**Errors:**
- `400` - Неверный формат запроса
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

---

#### 9. Получить все комнаты пользователя
**GET** `/api/v1/rooms`

Возвращает список всех комнат, в которых участвует пользователь.

**Response (200 OK):**
```json
{
  "rooms": [
    {
      "id": "uuid",
      "name": "string",
      "owner_id": "uuid"
    }
  ]
}
```

**Errors:**
- `401` - Не авторизован
- `500` - Внутренняя ошибка сервера

---

### Эндпоинты конкретной комнаты

> **Требуется:** Участие в комнате (проверяется middleware)

#### 10. Получить информацию о комнате
**GET** `/api/v1/rooms/:room_id`

Возвращает детальную информацию о комнате.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Response (200 OK):**
```json
{
  "id": "uuid",
  "name": "string",
  "owner_id": "uuid"
}
```

**Errors:**
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

#### 11. Обновить комнату
**PUT** `/api/v1/rooms/:room_id`

Обновляет название комнаты. Только владелец может обновлять комнату.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Request Body:**
```json
{
  "name": "string"
}
```

**Response (200 OK):**
```json
{
  "id": "uuid",
  "name": "string",
  "owner_id": "uuid"
}
```

**Errors:**
- `400` - Неверный формат запроса
- `401` - Не авторизован
- `403` - Не владелец комнаты
- `500` - Внутренняя ошибка сервера

---

#### 12. Удалить комнату
**DELETE** `/api/v1/rooms/:room_id`

Удаляет комнату. Только владелец может удалить комнату.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Response (204 No Content)**

**Errors:**
- `401` - Не авторизован
- `403` - Не владелец комнаты
- `500` - Внутренняя ошибка сервера

---

### Управление играми

#### 13. Добавить игру в комнату
**POST** `/api/v1/rooms/:room_id/games`

Добавляет новую игру в список для голосования.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Request Body:**
```json
{
  "title": "string"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "room_id": "uuid",
  "title": "string"
}
```

**Errors:**
- `400` - Неверный формат запроса
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

#### 14. Получить список игр комнаты
**GET** `/api/v1/rooms/:room_id/games`

Возвращает все игры в комнате.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Response (200 OK):**
```json
[
  {
    "id": "uuid",
    "room_id": "uuid",
    "title": "string"
  }
]
```

**Errors:**
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

#### 15. Удалить игру
**DELETE** `/api/v1/rooms/:room_id/games/:game_id`

Удаляет игру из комнаты.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты
- `game_id` (uuid) - ID игры

**Response (204 No Content)**

**Errors:**
- `401` - Не авторизован
- `403` - Игра не принадлежит указанной комнате
- `500` - Внутренняя ошибка сервера

---

### Управление участниками

#### 16. Пригласить участника
**POST** `/api/v1/rooms/:room_id/participants`

Добавляет нового участника в комнату.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Request Body:**
```json
{
  "user_id": "uuid"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "room_id": "uuid",
  "user_id": "uuid",
  "role": "member"
}
```

**Errors:**
- `400` - Неверный формат запроса
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `409` - Участник уже в комнате
- `500` - Внутренняя ошибка сервера

---

#### 17. Получить список участников
**GET** `/api/v1/rooms/:room_id/participants`

Возвращает всех участников комнаты с их ролями.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Response (200 OK):**
```json
{
  "users": [
    {
      "id": "uuid",
      "name": "string"
    }
  ],
  "roles": [
    "owner" | "member"
  ]
}
```

**Errors:**
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

#### 18. Удалить участника (покинуть комнату)
**DELETE** `/api/v1/rooms/:room_id/participants`

Удаляет текущего пользователя из участников комнаты.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Response (204 No Content)**

**Errors:**
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

### Управление голосами

#### 19. Добавить голос за игру
**POST** `/api/v1/rooms/:room_id/votes`

Голосует за конкретную игру.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Request Body:**
```json
{
  "game_id": "uuid"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "room_id": "uuid",
  "game_id": "uuid",
  "user_id": "uuid"
}
```

**Errors:**
- `400` - Неверный формат запроса
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

#### 20. Получить все голоса комнаты
**GET** `/api/v1/rooms/:room_id/votes`

Возвращает все голоса в комнате.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Response (200 OK):**
```json
[
  {
    "id": "uuid",
    "room_id": "uuid",
    "game_id": "uuid",
    "user_id": "uuid"
  }
]
```

**Errors:**
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

#### 21. Удалить свой голос
**DELETE** `/api/v1/rooms/:room_id/votes/:vote_id`

Удаляет голос. Пользователь может удалить только свой собственный голос.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты
- `vote_id` (uuid) - ID голоса

**Response (204 No Content)**

**Errors:**
- `401` - Не авторизован
- `403` - Попытка удалить чужой голос или нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

### Случайный выбор игры

#### 22. Получить случайный результат
**GET** `/api/v1/rooms/:room_id/random`

Генерирует новый случайный выбор игры на основе голосов и сохраняет результат.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Response (200 OK):**
```json
"uuid"
```

Возвращает `game_id` выбранной игры как строку UUID.

**Errors:**
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

#### 23. Получить последний результат
**GET** `/api/v1/rooms/:room_id/random/last`

Возвращает последний выбранный результат без генерации нового.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Response (200 OK):**
```json
"uuid"
```

Возвращает `game_id` последнего результата как строку UUID.

**Errors:**
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

#### 24. Получить историю результатов
**GET** `/api/v1/rooms/:room_id/random/history`

Возвращает полную историю всех случайных выборов в комнате.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Response (200 OK):**
```json
[
  "uuid"
]
```

Возвращает массив `game_id` (строки UUID) по истории.

**Errors:**
- `401` - Не авторизован
- `403` - Нет доступа к комнате
- `500` - Внутренняя ошибка сервера

---

## WebSocket Real-Time Updates

### WebSocket Connection
**WS** `/api/v1/rooms/:room_id/ws`

Устанавливает WebSocket-соединение для получения обновлений в реальном времени.

**URL Parameters:**
- `room_id` (uuid) - ID комнаты

**Query Parameters:**
- `token` (string, required) - JWT токен доступа

**Headers:**
- `Upgrade: websocket`
- `Connection: Upgrade`

**Connection Flow:**
1. Клиент отправляет HTTP GET с query-параметром `token=<jwt>` и заголовками для апгрейда
2. Сервер валидирует JWT токен из query-параметра
3. Соединение апгрейдится до WebSocket
4. Сервер отправляет приветственное сообщение: `{"type":"connected","room_id":"<uuid>"}`
5. Клиент получает все события комнаты в формате JSON

**Heartbeat:**
- Клиент должен отправлять Ping фреймы
- Сервер отвечает Pong фреймами
- Таймаут чтения: 60 секунд
- Таймаут записи: 10 секунд

### Event Types

Все события имеют структуру:
```json
{
  "type": "event_type",
  "room_id": "uuid",
  "payload": {...},
  "ts": 1733090000000
}
```

#### 1. Room Updated
**Type:** `room.updated`

Отправляется при изменении настроек комнаты.

**Payload:**
```json
{
  "id": "uuid",
  "name": "string"
}
```

#### 2. Participant Added
**Type:** `participant.added`

Отправляется при добавлении участника в комнату.

**Payload:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "role": "member"
}
```

#### 3. Participant Left
**Type:** `participant.left`

Отправляется при удалении участника из комнаты.

**Payload:**
```json
{
  "user_id": "uuid"
}
```

#### 4. Game Added
**Type:** `game.added`

Отправляется при добавлении игры в комнату.

**Payload:**
```json
{
  "id": "uuid",
  "title": "string"
}
```

#### 5. Game Deleted
**Type:** `game.deleted`

Отправляется при удалении игры из комнаты.

**Payload:**
```json
{
  "id": "uuid"
}
```

#### 6. Vote Added
**Type:** `vote.added`

Отправляется при добавлении голоса за игру.

**Payload:**
```json
{
  "id": "uuid",
  "game_id": "uuid",
  "user_id": "uuid"
}
```

#### 7. Vote Deleted
**Type:** `vote.deleted`

Отправляется при удалении голоса.

**Payload:**
```json
{
  "id": "uuid"
}
```

#### 8. Results Updated
**Type:** `results.updated`

Отправляется при изменении результатов случайного выбора.

**Payload:**
```json
{
  "result_id": "uuid"
}
```

**Errors:**
- `401` - Не авторизован (токен невалиден или отсутствует в query)
- `426` - Upgrade Required (отсутствуют заголовки WebSocket)

---

## Коды ошибок

| Код | Описание |
|-----|----------|
| 200 | OK - Успешный запрос |
| 201 | Created - Ресурс успешно создан |
| 204 | No Content - Успешно, нет содержимого для возврата |
| 400 | Bad Request - Неверный формат запроса |
| 401 | Unauthorized - Требуется аутентификация |
| 403 | Forbidden - Доступ запрещен |
| 404 | Not Found - Ресурс не найден |
| 409 | Conflict - Конфликт (например, ресурс уже существует) |
| 500 | Internal Server Error - Внутренняя ошибка сервера |

## Общий формат ошибок

```json
{
  "error": "Описание ошибки"
}
```

## Аутентификация

API использует два типа токенов:
1. **Access Token (JWT)** - отправляется в заголовке `Authorization: Bearer <token>`
   - Используется для авторизации запросов
   - Короткий срок действия

2. **Refresh Token** - хранится в HTTP-only cookie
   - Используется для получения нового access token
   - Длительный срок действия

## Middleware

### AuthMiddleware
Проверяет наличие и валидность JWT токена. Применяется ко всем эндпоинтам под `/api/v1`.

### CheckRoomMiddleware
Проверяет, является ли пользователь участником комнаты. Применяется ко всем эндпоинтам под `/api/v1/rooms/:room_id`.
