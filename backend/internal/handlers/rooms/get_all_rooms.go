package rooms

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	servicerooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type GetAllRoomsHandler struct {
	roomService servicerooms.RoomService
}

func NewGetAllRoomsHandler(roomService servicerooms.RoomService) *GetAllRoomsHandler {
	return &GetAllRoomsHandler{roomService: roomService}
}

type GetAllRoomsResponse struct {
	Rooms []rooms.Room `json:"rooms"`
}

func (h *GetAllRoomsHandler) HandleGetAllRooms(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	rooms, err := h.roomService.GetAllForUser(c.Context(), userID)
	if err != nil {
		logger.Errorf(c.Context(), "GetAllRooms Handle GetAllForUser error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get rooms"},
		)
	}

	response := GetAllRoomsResponse{
		Rooms: rooms,
	}

	return c.JSON(response)
}
