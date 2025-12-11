# Модель данных

Ниже описание схемы БД согласно миграциям `000001_init.up.sql` и `000002_add_entities.up.sql`.

## Таблицы

### users
| Поле | Тип | Ограничения |
| --- | --- | --- |
| id | UUID | PK |
| name | VARCHAR(50) | NOT NULL, UNIQUE |
| password_hash | VARCHAR(255) | NOT NULL |
| created_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP |

### refresh_tokens
| Поле | Тип | Ограничения |
| --- | --- | --- |
| token | VARCHAR(255) | PK |
| user_id | UUID | NOT NULL, FK → users(id), ON DELETE CASCADE |
| expires_at | TIMESTAMPTZ | NOT NULL |
| created_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP |

### rooms
| Поле | Тип | Ограничения |
| --- | --- | --- |
| id | UUID | PK |
| name | TEXT | NOT NULL |
| owner_id | UUID | NOT NULL, FK → users(id), ON DELETE CASCADE |
| created_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP |

### room_participants
| Поле | Тип | Ограничения |
| --- | --- | --- |
| id | UUID | PK |
| room_id | UUID | NOT NULL, FK → rooms(id), ON DELETE CASCADE |
| user_id | UUID | NOT NULL, FK → users(id), ON DELETE CASCADE |
| role | VARCHAR(20) | NOT NULL, DEFAULT 'member', CHECK role IN ('owner','member') |
| created_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP |
| (room_id, user_id) | — | UNIQUE (участник один раз в комнате) |

### games
| Поле | Тип | Ограничения |
| --- | --- | --- |
| id | UUID | PK |
| room_id | UUID | NOT NULL, FK → rooms(id), ON DELETE CASCADE |
| title | TEXT | NOT NULL |
| created_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP |

### votes
| Поле | Тип | Ограничения |
| --- | --- | --- |
| id | UUID | PK |
| room_id | UUID | NOT NULL, FK → rooms(id), ON DELETE CASCADE |
| game_id | UUID | NOT NULL, FK → games(id), ON DELETE CASCADE |
| user_id | UUID | NOT NULL, FK → users(id), ON DELETE CASCADE |
| created_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP |
| (room_id, game_id, user_id) | — | UNIQUE (пользователь голосует один раз за игру в комнате) |

### random_results
| Поле | Тип | Ограничения |
| --- | --- | --- |
| id | UUID | PK |
| room_id | UUID | NOT NULL, FK → rooms(id), ON DELETE CASCADE |
| game_id | UUID | NOT NULL, FK → games(id) |
| chosen_by | UUID | NOT NULL, FK → users(id) |
| created_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP |

## Связи
- `users` 1—N `refresh_tokens` (каскадное удаление токенов при удалении пользователя).
- `users` 1—N `rooms` через `owner_id` (комнаты удаляются при удалении владельца).
- `users` 1—N `room_participants`, `rooms` 1—N `room_participants`; уникальность пары ограничивает дубликаты.
- `rooms` 1—N `games`; при удалении комнаты удаляются игры и каскадно связанные голоса.
- `games` 1—N `votes`; `users` 1—N `votes`; уникальный состав (room, game, user) предотвращает повторные голоса.
- `rooms` 1—N `votes` (через room_id) — голос принадлежит конкретной комнате.
- `rooms` 1—N `random_results`; `games` 1—N `random_results`; `users` 1—N `random_results` (кто выбрал).

## Ключевые инварианты
- Комната принадлежит владельцу (`owner_id`) и исчезает при удалении владельца.
- Участник не может быть добавлен в одну комнату дважды.
- Голос уникален для сочетания комната+игра+пользователь.
- Все сущности, связанные с комнатой, удаляются каскадно при удалении комнаты (участники, игры, голоса, результаты выбора).
- Токены и связанные сущности пользователей удаляются каскадно при удалении пользователя.
