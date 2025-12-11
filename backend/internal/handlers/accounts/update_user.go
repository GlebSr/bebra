package accounts

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/profile"
	serviceusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/users"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserHandler struct {
	userService serviceusers.UserService
}

func NewUpdateUserHandler(userService serviceusers.UserService) *UpdateUserHandler {
	return &UpdateUserHandler{
		userService: userService,
	}
}

type UpdateUserRequest struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

func (h *UpdateUserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Errorf(c.Context(), "HandleUpdateUser failed to parse body: %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var passwordHash string
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Errorf(c.Context(), "HandleUpdateUser failed to hash password: %v", err)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}

		passwordHash = string(hash)
	}

	_, err := h.userService.Update(c.Context(), profile.User{
		ID:           userID,
		Name:         req.Name,
		PasswordHash: passwordHash,
	})
	if err != nil {
		logger.Errorf(c.Context(), "HandleUpdateUser failed to update user: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
