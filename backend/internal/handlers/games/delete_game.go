package games

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/hub"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/games"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/participants"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type DeleteGameHandler struct {
	gameService        games.GameService
	participantService participants.ParticipantService
	Hub                hub.Hub
}

func NewDeleteGameHandler(gameService games.GameService, participantService participants.ParticipantService) *DeleteGameHandler {
	return &DeleteGameHandler{gameService: gameService, participantService: participantService}
}

func (h *DeleteGameHandler) Handle(c *fiber.Ctx) error {
	room_id := c.Locals("room_id").(string)
	game_id := c.Params("game_id")

	game, err := h.gameService.Get(c.Context(), game_id)
	if err != nil {
		logger.Errorf(c.Context(), "DeleteGame Handle Get game error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get game"},
		)
	}

	if game.RoomID != room_id {
		return c.Status(fiber.StatusForbidden).JSON(
			fiber.Map{"error": "You are not allowed to delete this game from another room"},
		)
	}

	if err := h.gameService.Delete(c.Context(), game_id, room_id); err != nil {
		logger.Errorf(c.Context(), "DeleteGame Handle Delete error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to delete game"},
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
