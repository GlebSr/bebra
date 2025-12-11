package rooms

import (
	servicerooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type DeleteRoomHandler struct {
	roomService servicerooms.RoomService
}

func NewDeleteRoomHandler(roomService servicerooms.RoomService) *DeleteRoomHandler {
	return &DeleteRoomHandler{roomService: roomService}
}

func (h *DeleteRoomHandler) HandleDeleteRoom(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	roomID := c.Locals("room_id").(string)

	room, err := h.roomService.GetByID(c.Context(), roomID)
	if err != nil {
		logger.Errorf(c.Context(), "DeleteRoom Handle GetByID error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get room"},
		)
	}

	if room.OwnerID != userID {
		logger.Errorf(c.Context(), "DeleteRoom Handle unauthorized user: %v", userID)

		return c.Status(fiber.StatusForbidden).JSON(
			fiber.Map{"error": "You are not the owner of this room"},
		)
	}

	if err := h.roomService.Delete(c.Context(), roomID); err != nil {
		logger.Errorf(c.Context(), "DeleteRoom Handle Delete error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to delete room"},
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
