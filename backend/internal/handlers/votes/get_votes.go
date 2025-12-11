package votes

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/votes"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type GetVotesHandler struct {
	voteService votes.VoteService
}

func NewGetVotesHandler(voteService votes.VoteService) *GetVotesHandler {
	return &GetVotesHandler{voteService: voteService}
}

func (h *GetVotesHandler) Handle(c *fiber.Ctx) error {
	room_id := c.Locals("room_id").(string)
	votesList, err := h.voteService.GetForRoom(c.Context(), room_id)
	if err != nil {
		logger.Errorf(c.Context(), "GetVotes Handle GetForRoom error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get votes"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(votesList)
}
