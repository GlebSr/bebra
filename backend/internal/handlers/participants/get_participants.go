package participants

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/participants"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type GetParticipantsHandler struct {
	participantService participants.ParticipantService
}

func NewGetParticipantsHandler(participantService participants.ParticipantService) *GetParticipantsHandler {
	return &GetParticipantsHandler{participantService: participantService}
}

func (h *GetParticipantsHandler) Handle(c *fiber.Ctx) error {
	roomID := c.Locals("room_id").(string)
	users, roles, err := h.participantService.GetAllParticipants(c.Context(), roomID)
	if err != nil {
		logger.Errorf(c.Context(), "GetParticipants Handle GetAllParticipants error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get participants"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
		"roles": roles,
	})
}
