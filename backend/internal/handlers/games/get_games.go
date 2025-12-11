package games

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/games"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type GetGamesHandler struct {
	gameService games.GameService
}

func NewGetGamesHandler(gameService games.GameService) *GetGamesHandler {
	return &GetGamesHandler{gameService: gameService}
}

func (h *GetGamesHandler) Handle(c *fiber.Ctx) error {
	room_id := c.Locals("room_id").(string)
	games, err := h.gameService.GetAllRoomGames(c.Context(), room_id)
	if err != nil {
		logger.Errorf(c.Context(), "GetGames Handle GetAllRoomGames error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get games"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(games)
}
