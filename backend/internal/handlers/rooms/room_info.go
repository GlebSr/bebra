package rooms

import (
	"context"

	servicerooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type GetRoomInfoHandler struct {
	roomService servicerooms.RoomService
}

func NewGetRoomInfoHandler(roomService servicerooms.RoomService) *GetRoomInfoHandler {
	return &GetRoomInfoHandler{roomService: roomService}
}

type GetRoomInfoResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

func (h *GetRoomInfoHandler) HandleGetRoomInfo(c *fiber.Ctx) error {
	roomID := c.Locals("room_id").(string)

	room, err := h.roomService.GetByID(context.Background(), roomID)
	if err != nil {
		logger.Errorf(c.Context(), "GetRoomInfo Handle GetByID error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get room info"},
		)
	}

	return c.JSON(GetRoomInfoResponse{
		ID:      room.ID,
		Name:    room.Name,
		OwnerID: room.OwnerID,
	})
}
