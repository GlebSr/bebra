package utils

import (
	"github.com/gofiber/fiber/v2"
)

func AddLoggData(ctx *fiber.Ctx, key string, value interface{}) {
	fieldsToLog, ok := ctx.Locals("fields_to_log").(map[string]interface{})
	if !ok {
		fieldsToLog = make(map[string]interface{})
	}

	fieldsToLog[key] = value
	ctx.Locals("fields_to_log", fieldsToLog)
}
