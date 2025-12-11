package accounts

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/tokens"
	servicetokens "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/tokens"
	serviceusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/users"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type RefreshHandler struct {
	tokenService servicetokens.TokenService
	userService  serviceusers.UserService
}

func NewRefreshHandler(tokensService tokens.TokenService, service serviceusers.UserService) *RefreshHandler {
	return &RefreshHandler{userService: service, tokenService: tokensService}
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *RefreshHandler) HandleRefresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	userID, err := h.tokenService.ValidateRefreshToken(c.Context(), refreshToken)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to validate refresh token: %v", err)

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to get user: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	if user.ID == "" {
		logger.Errorf(c.Context(), "User not found")

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	newAccessToken, err := h.tokenService.GenerateJWTToken(c.Context(), user.ID)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to generate new access token: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate new access token",
		})
	}

	newRefreshToken, err := h.tokenService.CreateRefreshToken(c.Context(), user.ID)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to generate new refresh token: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate new refresh token",
		})
	}

	cookie := utils.CreateRefreshTokenCookie(newRefreshToken)
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(RefreshResponse{
		AccessToken: newAccessToken,
	})
}
