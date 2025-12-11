package votes

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/votes"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddVoteHandler struct {
	voteService votes.VoteService
}

func NewAddVoteHandler(voteService votes.VoteService) *AddVoteHandler {
	return &AddVoteHandler{voteService: voteService}
}

type AddVoteRequest struct {
	GameID string `json:"game_id"`
}

func (h *AddVoteHandler) Handle(c *fiber.Ctx) error {
	room_id := c.Locals("room_id").(string)

	var req AddVoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	vote, err := h.voteService.Add(c.Context(), rooms.Vote{
		ID:     uuid.New().String(),
		RoomID: room_id,
		GameID: req.GameID,
		UserID: c.Locals("user_id").(string),
	})
	if err != nil {
		logger.Errorf(c.Context(), "AddVote Handle Add error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to add vote"},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(vote)
}
