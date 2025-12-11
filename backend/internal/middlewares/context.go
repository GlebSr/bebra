package middlewares

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func LogFieldsMiddleware(ctx *fiber.Ctx) error {
	fieldsToLog := make(map[string]interface{})

	fieldsToLog["method"] = ctx.Method()
	fieldsToLog["path"] = ctx.Path()
	fieldsToLog["ip"] = ctx.IP()
	fieldsToLog["hostname"] = ctx.Hostname()
	fieldsToLog["protocol"] = ctx.Protocol()
	fieldsToLog["query"] = ctx.Queries()
	if !websocket.IsWebSocketUpgrade(ctx) {
		fieldsToLog["body"] = string(ctx.Body())
		fieldsToLog["headers"] = ctx.GetReqHeaders()
	}
	fieldsToLog["params"] = ctx.AllParams()
	fieldsToLog["trace_id"] = uuid.New().String()

	ctx.Locals("fields_to_log", fieldsToLog)

	logger.Infof(ctx.Context(), "Got request to %s", ctx.Path())
	return ctx.Next()
}
