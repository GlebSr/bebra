package participants

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/participants"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type DeleteParticipantHandler struct {
	participantService participants.ParticipantService
}

func NewDeleteParticipantHandler(participantService participants.ParticipantService) *DeleteParticipantHandler {
	return &DeleteParticipantHandler{participantService: participantService}
}

func (h *DeleteParticipantHandler) Handle(c *fiber.Ctx) error {
	roomID := c.Locals("room_id").(string)
	userID := c.Locals("user_id").(string)
	if err := h.participantService.Delete(c.Context(), roomID, userID); err != nil {
		logger.Errorf(c.Context(), "DeleteParticipant Handle Delete error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to delete participant"},
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
