package accounts

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/profile"
	servicetokens "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/tokens"
	serviceusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/users"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignUpHandler struct {
	tokenService servicetokens.TokenService
	userService  serviceusers.UserService
}

func NewSignupHandler(tokenService servicetokens.TokenService, userService serviceusers.UserService) *SignUpHandler {
	return &SignUpHandler{
		tokenService: tokenService,
		userService:  userService,
	}
}

type SignupRequest struct {
	Password string `json:"password"`
	Name     string `json:"name"`
}

type SignupResponse struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

func (h *SignUpHandler) HandleSignup(c *fiber.Ctx) error {
	var body SignupRequest
	if err := c.BodyParser(&body); err != nil {
		logger.Errorf(c.Context(), "Failed to parse request body: %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.GetByName(c.Context(), body.Name)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to check existing user: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check existing user",
		})
	}

	if user.ID != "" {
		logger.Infof(c.Context(), "User with name %s already exists id: %s", body.Name, user.ID)

		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this name already exists",
		})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to hash password: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	user = profile.User{
		ID:           uuid.New().String(),
		PasswordHash: string(passwordHash),
		Name:         body.Name,
	}

	createdUser, err := h.userService.Add(c.Context(), user)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to create user: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	accessToken, err := h.tokenService.GenerateJWTToken(c.Context(), createdUser.ID)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to generate JWT token: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	refreshToken, err := h.tokenService.CreateRefreshToken(c.Context(), createdUser.ID)
	if err != nil {
		logger.Errorf(c.Context(), "Failed to generate refresh token: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	cookie := utils.CreateRefreshTokenCookie(refreshToken)
	c.Cookie(&cookie)

	return c.Status(fiber.StatusCreated).JSON(SignupResponse{
		UserID:      createdUser.ID,
		AccessToken: accessToken,
	})
}
