package random

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/results"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetRandomHandler struct {
	resultService results.ResultService
}

func NewGetRandomHandler(resultService results.ResultService) *GetRandomHandler {
	return &GetRandomHandler{resultService: resultService}
}

func (h *GetRandomHandler) Handle(c *fiber.Ctx) error {
	room_id := c.Locals("room_id").(string)
	user_id := c.Locals("user_id").(string)
	randomResult, err := h.resultService.PickResult(c.Context(), room_id)
	if err != nil {
		logger.Errorf(c.Context(), "GetRandom Handle GenerateRandomResult error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get random result"},
		)
	}

	_, err = h.resultService.Add(c.Context(), rooms.Result{
		ID:       uuid.New().String(),
		RoomID:   room_id,
		GameID:   randomResult,
		ChosenBy: user_id,
	})
	if err != nil {
		logger.Errorf(c.Context(), "GetRandom Handle Add Result error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to save random result"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(randomResult)
}
