package rooms

import (
	"context"
	"time"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/hub"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/tokens"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WSRoomHandler struct {
	Hub          hub.Hub
	tokenService tokens.TokenService
}

func NewWSRoomHandler(h hub.Hub, ts tokens.TokenService) *WSRoomHandler {
	return &WSRoomHandler{Hub: h, tokenService: ts}
}

// WebSocket upgrade middleware and authentication
func (h *WSRoomHandler) Handle(c *fiber.Ctx) error {
	token := c.Query("token")
	userID, err := h.tokenService.ValidateJWTToken(c.Context(), token)
	if err != nil {
		logger.Warnf(c.Context(), "Invalid token: %v", err)

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	c.Locals("user_id", userID)
	logger.Infof(c.Context(), "Authenticated user with ID: %s", userID)

	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// Conn handles the actual websocket connection after upgrade.
func (h *WSRoomHandler) Conn(c *websocket.Conn) {
	roomID := c.Params("room_id")
	userID := c.Locals("user_id").(string)

	cl := h.Hub.Subscribe(roomID, &hub.Client{
		Conn:   c,
		UserID: userID,
	})
	defer func() {
		h.Hub.Unsubscribe(roomID, cl)
	}()

	// Initial ping to confirm connection
	err := c.WriteMessage(websocket.TextMessage, []byte(`{"type":"connected","room_id":"`+roomID+`"}`))
	if err != nil {
		logger.Warnf(context.Background(), "WebSocket write error: %v", err)

		return
	}

	c.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.SetPongHandler(func(appData string) error {
		c.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Read loop (client messages are ignored except for ping/pong)
	for {
		mt, _, err := c.ReadMessage()
		if err != nil {
			logger.Warnf(context.Background(), "WebSocket read error: %v", err)
			// WebSocket closed, no context needed
			break
		}
		if mt == websocket.PingMessage {
			_ = c.WriteMessage(websocket.PongMessage, nil)
		}
	}
}
