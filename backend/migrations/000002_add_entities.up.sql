-- ROOMS
CREATE TABLE rooms (
  id           UUID PRIMARY KEY,
  name         TEXT NOT NULL,
  owner_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- PARTICIPANTS
CREATE TABLE room_participants (
  id         UUID PRIMARY KEY,
  room_id    UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role       VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'member')),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(room_id, user_id)
);

-- GAMES
CREATE TABLE games (
  id        UUID PRIMARY KEY,
  room_id   UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  title     TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- VOTES
CREATE TABLE votes (
  id        UUID PRIMARY KEY,
  room_id   UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  game_id   UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(room_id, game_id, user_id)  -- чтобы один пользователь не голосовал дважды за ту же игру
);

-- RANDOM RESULTS
CREATE TABLE random_results (
  id        UUID PRIMARY KEY,
  room_id   UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  game_id   UUID NOT NULL REFERENCES games(id),
  chosen_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);