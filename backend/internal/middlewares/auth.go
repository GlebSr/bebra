package middlewares

import (
	"strings"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/tokens"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	tokenService tokens.TokenService
}

func NewAuthMiddleware(tokenService tokens.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: tokenService}
}
func (m *AuthMiddleware) AuthRequired(c *fiber.Ctx) error {
	tokens := strings.Split(c.Get("Authorization"), "Bearer ")
	if len(tokens) != 2 || tokens[1] == "" {
		logger.Warnf(c.Context(), "Missing token in request")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := m.tokenService.ValidateJWTToken(c.Context(), tokens[1])
	if err != nil {
		logger.Warnf(c.Context(), "Invalid token: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	c.Locals("user_id", userID)
	logger.Infof(c.Context(), "Authenticated user with ID: %s", userID)

	return c.Next()
}
