package random

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/results"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type GetHistoryHandler struct {
	resultService results.ResultService
}

func NewGetHistoryHandler(resultService results.ResultService) *GetHistoryHandler {
	return &GetHistoryHandler{resultService: resultService}
}

func (h *GetHistoryHandler) Handle(c *fiber.Ctx) error {
	room_id := c.Locals("room_id").(string)
	history, err := h.resultService.GetAllResults(c.Context(), room_id)
	if err != nil {
		logger.Errorf(c.Context(), "GetHistory Handle GetAllResults error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get history"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(history)
}
