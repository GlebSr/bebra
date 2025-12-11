package hub

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

// RoomEventType defines types of realtime events in a room.
type RoomEventType string

const (
	EventRoomUpdated      RoomEventType = "room.updated"
	EventParticipantAdded RoomEventType = "participant.added"
	EventParticipantLeft  RoomEventType = "participant.left"
	EventGameAdded        RoomEventType = "game.added"
	EventGameDeleted      RoomEventType = "game.deleted"
	EventVoteAdded        RoomEventType = "vote.added"
	EventVoteDeleted      RoomEventType = "vote.deleted"
	EventResultsUpdated   RoomEventType = "results.updated"
)

// RoomEvent is a generic broadcast payload.
type RoomEvent struct {
	Type    RoomEventType `json:"type"`
	RoomID  string        `json:"room_id"`
	Payload any           `json:"payload"`
	Ts      int64         `json:"ts"`
}

type Client struct {
	Conn   *websocket.Conn
	UserID string
}

type Hub interface {
	Subscribe(roomID string, cl *Client) *Client
	Unsubscribe(roomID string, cl *Client)
	Broadcast(roomID string, evt RoomEvent)
}

// Hub manages per-room websocket clients and broadcasts events.
type HubWS struct {
	mu        sync.RWMutex
	rooms     map[string]map[*Client]struct{}
	writeWait time.Duration
}

func NewHubWS() *HubWS {
	return &HubWS{
		rooms:     make(map[string]map[*Client]struct{}),
		writeWait: 10 * time.Second,
	}
}

// Subscribe adds a client to a room.
func (h *HubWS) Subscribe(roomID string, cl *Client) *Client {
	h.mu.Lock()
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[*Client]struct{})
	}
	h.rooms[roomID][cl] = struct{}{}
	h.mu.Unlock()
	return cl
}

// Unsubscribe removes a client from a room.
func (h *HubWS) Unsubscribe(roomID string, cl *Client) {
	h.mu.Lock()
	if clients, ok := h.rooms[roomID]; ok {
		delete(clients, cl)
		if len(clients) == 0 {
			delete(h.rooms, roomID)
		}
	}
	h.mu.Unlock()
}

// Broadcast sends an event to all clients in a room.
func (h *HubWS) Broadcast(roomID string, evt RoomEvent) {
	h.mu.RLock()
	clients := h.rooms[roomID]
	h.mu.RUnlock()
	if len(clients) == 0 {
		return
	}

	if evt.Ts == 0 {
		evt.Ts = time.Now().UnixMilli()
	}
	msg, _ := json.Marshal(evt)

	for cl := range clients {
		cl.Conn.SetWriteDeadline(time.Now().Add(h.writeWait))
		if err := cl.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			cl.Conn.Close()
			h.Unsubscribe(roomID, cl)
		}
	}
}
