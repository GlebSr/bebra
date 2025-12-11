package rooms

import (
	"context"

	servicerooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type UpdateRoomHandler struct {
	roomService servicerooms.RoomService
}

func NewUpdateRoomHandler(roomService servicerooms.RoomService) *UpdateRoomHandler {
	return &UpdateRoomHandler{roomService: roomService}
}

type UpdateRoomRequest struct {
	Name string `json:"name"`
}

type UpdateRoomResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

func (h *UpdateRoomHandler) HandleUpdateRoom(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	roomID := c.Locals("room_id").(string)

	var req UpdateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Errorf(c.Context(), "UpdateRoom Handle BodyParser error: %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "Invalid request body"},
		)
	}

	room, err := h.roomService.GetByID(context.Background(), roomID)
	if err != nil {
		logger.Errorf(c.Context(), "UpdateRoom Handle GetByID error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get room"},
		)
	}

	if room.OwnerID != userID {
		logger.Errorf(c.Context(), "UpdateRoom Handle unauthorized user: %v", userID)

		return c.Status(fiber.StatusForbidden).JSON(
			fiber.Map{"error": "You are not the owner of this room"},
		)
	}

	room.Name = req.Name

	updatedRoom, err := h.roomService.Update(c.Context(), room)
	if err != nil {
		logger.Errorf(c.Context(), "UpdateRoom Handle Update error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to update room"},
		)
	}

	response := UpdateRoomResponse{
		ID:      updatedRoom.ID,
		Name:    updatedRoom.Name,
		OwnerID: updatedRoom.OwnerID,
	}

	return c.JSON(response)
}
