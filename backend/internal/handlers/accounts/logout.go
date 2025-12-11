package accounts

import (
	servicetokens "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/tokens"
	serviceusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/users"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type LogoutHandler struct {
	tokenService servicetokens.TokenService
	userService  serviceusers.UserService
}

func NewLogoutHandler(tokenService servicetokens.TokenService, userService serviceusers.UserService) *LogoutHandler {
	return &LogoutHandler{
		tokenService: tokenService,
		userService:  userService,
	}
}

func (h *LogoutHandler) HandleLogout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	err := h.tokenService.DeleteRefreshToken(c.Context(), refreshToken)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to delete refresh token: %v", err)
	}

	c.ClearCookie("refresh_token")

	return c.SendStatus(fiber.StatusNoContent)
}
