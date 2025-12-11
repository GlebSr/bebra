// Пользователи
export interface User {
  id: string;
  name: string;
}

export interface AuthResponse {
  user_id: string;
  access_token: string;
}

export interface SignUpRequest {
  name: string;
  password: string;
}

export interface SignInRequest {
  name: string;
  password: string;
}

export interface UpdateUserRequest {
  name?: string;
  password?: string;
}

// Комнаты
export interface Room {
  id: string;
  name: string;
  owner_id: string;
}

export interface RoomsResponse {
  rooms: Room[];
}

export interface CreateRoomRequest {
  name: string;
}

export interface UpdateRoomRequest {
  name: string;
}

// Игры
export interface Game {
  id: string;
  room_id: string;
  title: string;
  created_at?: string;
}

export type GamesResponse = Game[];

export interface CreateGameRequest {
  title: string;
}

// Участники
export type ParticipantRole = 'owner' | 'member';

export interface ParticipantsResponse {
  users: User[];
  roles: ParticipantRole[];
}

export interface InviteParticipantRequest {
  name: string;
}

// Голоса
export interface Vote {
  id: string;
  room_id: string;
  game_id: string;
  user_id: string;
}

export interface VotesResponse {
  votes: Vote[];
}

export interface AddVoteRequest {
  game_id: string;
}

// Случайный выбор
// Возвращает game_id как UUID строку
export type RandomResult = string;

// История - массив game_id
export type RandomHistoryResponse = string[];

// WebSocket
export type WSEventType = 
  | 'connected'
  | 'room.updated'
  | 'participant.added'
  | 'participant.left'
  | 'game.added'
  | 'game.deleted'
  | 'vote.added'
  | 'vote.deleted'
  | 'results.updated';

export interface WSEvent<T = unknown> {
  type: WSEventType;
  room_id: string;
  payload: T;
  ts: number;
}

export interface WSConnectedPayload {
  user_id: string;
}

export interface WSRoomUpdatedPayload {
  name: string;
}

export interface WSParticipantPayload {
  user_id: string;
}

export interface WSGamePayload {
  game_id: string;
  title?: string;
}

export interface WSVotePayload {
  vote_id: string;
  game_id: string;
  user_id: string;
}

export interface WSResultsPayload {
  game_id: string;
}

// Ошибки
export interface ApiError {
  error: string;
}
