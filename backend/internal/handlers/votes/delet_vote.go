package votes

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/votes"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type DeleteVoteHandler struct {
	voteService votes.VoteService
}

func NewDeleteVoteHandler(voteService votes.VoteService) *DeleteVoteHandler {
	return &DeleteVoteHandler{voteService: voteService}
}

func (h *DeleteVoteHandler) Handle(c *fiber.Ctx) error {
	voteID := c.Params("vote_id")
	vote, err := h.voteService.Get(c.Context(), voteID)
	if err != nil {
		logger.Errorf(c.Context(), "DeleteVote Handle Get error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get vote"},
		)
	}

	if vote.UserID != c.Locals("user_id").(string) {
		return c.Status(fiber.StatusForbidden).JSON(
			fiber.Map{"error": "You are not allowed to delete this vote"},
		)
	}

	roomID := c.Locals("room_id").(string)
	if err := h.voteService.Delete(c.Context(), voteID, roomID); err != nil {
		logger.Errorf(c.Context(), "DeleteVote Handle Delete error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to delete vote"},
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
