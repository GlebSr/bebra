package participants

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/participants"
	serviceusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/users"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InviteHandler struct {
	participantService participants.ParticipantService
	userService        serviceusers.UserService
}

func NewInviteHandler(participantService participants.ParticipantService, userService serviceusers.UserService) *InviteHandler {
	return &InviteHandler{participantService: participantService, userService: userService}
}

type InviteRequest struct {
	Name string `json:"name"`
}

func (h *InviteHandler) Handle(c *fiber.Ctx) error {
	roomID := c.Locals("room_id").(string)
	var req InviteRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Errorf(c.Context(), "Invite Handle BodyParser error: %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "Invalid request body"},
		)
	}

	user, err := h.userService.GetByName(c.Context(), req.Name)
	if err != nil {
		logger.Errorf(c.Context(), "Invite Handle GetByName error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to get user by name"},
		)
	}

	participant, err := h.participantService.Get(c.Context(), roomID, user.ID)
	if err != nil {
		logger.Errorf(c.Context(), "Invite Handle GetByRoomAndUserID error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to check existing participant"},
		)
	}

	if participant.ID != "" {
		return c.Status(fiber.StatusConflict).JSON(
			fiber.Map{"error": "Participant already exists in the room"},
		)
	}

	newParticipant, err := h.participantService.Add(c.Context(), rooms.RoomParticipant{
		ID:     uuid.New().String(),
		RoomID: roomID,
		UserID: user.ID,
		Role:   "member",
	})
	if err != nil {
		logger.Errorf(c.Context(), "Invite Handle Add error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to invite participant"},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(newParticipant)
}
