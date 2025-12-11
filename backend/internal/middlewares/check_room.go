package middlewares

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/participants"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type CheckRoomMiddleware struct {
	roomService        rooms.RoomService
	participantService participants.ParticipantService
}

func NewCheckRoomMiddleware(roomService rooms.RoomService, participantService participants.ParticipantService) *CheckRoomMiddleware {
	return &CheckRoomMiddleware{roomService: roomService, participantService: participantService}
}

func (m *CheckRoomMiddleware) Handle(c *fiber.Ctx) error {
	room_id := c.Params("room_id")
	user_id := c.Locals("user_id").(string)
	room, err := m.roomService.GetByID(c.Context(), room_id)
	if err != nil {
		logger.Errorf(c.Context(), "CheckRoomMiddleware Handle GetRoom error: %v", err)

		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get room")
	}

	if room.ID == "" {
		return c.Status(fiber.StatusNotFound).SendString("Room not found")
	}

	participant, err := m.participantService.Get(c.Context(), room_id, user_id)
	if err != nil {
		logger.Errorf(c.Context(), "CheckRoomMiddleware Handle GetParticipant error: %v", err)

		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get participant")
	}

	if participant.ID == "" {
		return c.Status(fiber.StatusForbidden).SendString("You are not a participant of this room")
	}

	c.Locals("room_id", room_id)
	c.Locals("participant_id", participant.ID)

	return c.Next()
}
