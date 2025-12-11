package utils

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/profile"
	"github.com/gofiber/fiber/v2"
)

func CreateRefreshTokenCookie(token profile.RefreshToken) fiber.Cookie {
	cookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    token.Token,
		Expires:  token.ExpiresAt,
		HTTPOnly: true,
		SameSite: "Lax",
	}
	return cookie
}
