package random

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/results"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type GetLastHandler struct {
	resultService results.ResultService
}

func NewGetLastHandler(resultService results.ResultService) *GetLastHandler {
	return &GetLastHandler{resultService: resultService}
}

func (h *GetLastHandler) Handle(c *fiber.Ctx) error {
	room_id := c.Locals("room_id").(string)
	lastResult, err := h.resultService.GetLastResult(c.Context(), room_id)
	if err != nil {
		logger.Errorf(c.Context(), "GetLast Handle GetLastResult error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get last result"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(lastResult)
}
