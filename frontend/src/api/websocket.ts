import { tokenStorage } from './tokenStorage';
import type {
  WSEvent,
  WSEventType,
  WSConnectedPayload,
  WSRoomUpdatedPayload,
  WSParticipantPayload,
  WSGamePayload,
  WSVotePayload,
  WSResultsPayload,
} from './types';

type EventHandler = (event: WSEvent) => void;

export class RoomWebSocket {
  private ws: WebSocket | null = null;
  private roomId: string;
  private handlers: Map<WSEventType, EventHandler> = new Map();
  private pingInterval: number | null = null;
  private reconnectTimeout: number | null = null;
  private shouldReconnect = true;

  constructor(roomId: string) {
    this.roomId = roomId;
  }
  // Коннект к WebSocket серверу
  connect(): void {
    const token = tokenStorage.getToken();
    if (!token) {
      console.error('WebSocket: No access token');
      return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/v1/rooms/${this.roomId}/ws?token=${token}`;

    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.startPing();
    };
    // Обработка входящих сообщений
    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as WSEvent;
        this.handleEvent(data);
      } catch (err) {
        console.error('WebSocket: Failed to parse message', err);
      }
    };
    // Обработка отключения
    this.ws.onclose = () => {
      console.log('WebSocket disconnected');
      this.stopPing();
      if (this.shouldReconnect) {
        this.scheduleReconnect();
      }
    };
    // Обработка ошибок
    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }
  // Отключение от WebSocket сервера
  disconnect(): void {
    this.shouldReconnect = false;
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    this.stopPing();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
  // Регистрация обработчиков событий
  on(eventType: WSEventType, handler: EventHandler): void {
    this.handlers.set(eventType, handler);
  }
  // Удаление обработчиков событий
  off(eventType: WSEventType): void {
    this.handlers.delete(eventType);
  }
  // Обработка события
  private handleEvent(event: WSEvent): void {
    const handler = this.handlers.get(event.type);
    if (handler) {
      handler(event);
    }
  }
  // Пинг для поддержания соединения
  private startPing(): void {
    this.pingInterval = window.setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send('ping');
      }
    }, 30000);
  }
  // Остановка пинга
  private stopPing(): void {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }
  // Планирование повторного подключения
  private scheduleReconnect(): void {
    this.reconnectTimeout = window.setTimeout(() => {
      console.log('WebSocket: Attempting to reconnect...');
      this.connect();
    }, 3000);
  }
}
// Фабричная функция для создания экземпляра RoomWebSocket
export function createRoomWebSocket(roomId: string): RoomWebSocket {
  return new RoomWebSocket(roomId);
}
