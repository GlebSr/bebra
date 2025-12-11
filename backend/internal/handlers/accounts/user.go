package accounts

import (
	serviceusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/users"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService serviceusers.UserService
}

func NewUserHandler(userService serviceusers.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(UserResponse{
		ID:   user.ID,
		Name: user.Name,
	})
}
