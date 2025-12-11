package accounts

import (
	serviceusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/users"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type GetByNameHandler struct {
	userService serviceusers.UserService
}

func NewGetByNameHandler(userService serviceusers.UserService) *GetByNameHandler {
	return &GetByNameHandler{
		userService: userService,
	}
}

type GetByNameResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *GetByNameHandler) HandleGetByName(c *fiber.Ctx) error {
	userName := c.Query("name")
	user, err := h.userService.GetByName(c.Context(), userName)
	if err != nil {
		logger.Errorf(c.Context(), "HandleGetByName failed to get user by name: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	if user.ID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetByNameResponse{
		ID:   user.ID,
		Name: user.Name,
	})
}
