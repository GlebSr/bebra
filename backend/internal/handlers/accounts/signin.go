package accounts

import (
	servicetokens "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/tokens"
	serviceusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/users"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type SignInHandler struct {
	userService  serviceusers.UserService
	tokenService servicetokens.TokenService
}

func NewSignInHandler(tokenService servicetokens.TokenService, userService serviceusers.UserService) *SignInHandler {
	return &SignInHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

type SignInRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInResponse struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

func (h *SignInHandler) HandleSignIn(c *fiber.Ctx) error {
	var body SignInRequest
	if err := c.BodyParser(&body); err != nil {
		logger.Errorf(c.Context(), "Failed to parse request body: %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.GetByName(c.Context(), body.Name)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to get user by name and discriminator: %v", err)

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid name, discriminator, or password",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
	if err != nil {
		logger.Errorf(c.Context(), "Failed to compare password: %v", err)

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	accessToken, err := h.tokenService.GenerateJWTToken(c.Context(), user.ID)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to generate JWT token: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	refreshToken, err := h.tokenService.CreateRefreshToken(c.Context(), user.ID)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to generate refresh token: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	cookie := utils.CreateRefreshTokenCookie(refreshToken)
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(SignInResponse{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
