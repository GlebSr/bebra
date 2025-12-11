package games

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/games"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/participants"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddGameHandler struct {
	gameService        games.GameService
	participantService participants.ParticipantService
}

func NewAddGameHandler(gameService games.GameService, participantService participants.ParticipantService) *AddGameHandler {
	return &AddGameHandler{gameService: gameService, participantService: participantService}
}

type AddGameRequest struct {
	Title string `json:"title"`
}

func (h *AddGameHandler) Handle(c *fiber.Ctx) error {
	var req AddGameRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Errorf(c.Context(), "AddGame Handle BodyParser error: %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "Invalid request body"},
		)
	}

	room_id := c.Locals("room_id").(string)

	game, err := h.gameService.Add(c.Context(), rooms.Game{
		ID:     uuid.New().String(),
		RoomID: room_id,
		Title:  req.Title,
	})
	if err != nil {
		logger.Errorf(c.Context(), "AddGame Handle Add error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to add game"},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(game)
}
