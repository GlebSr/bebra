package rooms

import (
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	serviceparticipants "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/participants"
	servicerooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateRoomHandler struct {
	roomService        servicerooms.RoomService
	participantService serviceparticipants.ParticipantService
}

func NewCreateRoomHandler(roomService servicerooms.RoomService, participantService serviceparticipants.ParticipantService) *CreateRoomHandler {
	return &CreateRoomHandler{roomService: roomService, participantService: participantService}
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type CreateRoomResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

func (h *CreateRoomHandler) Handle(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req CreateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Errorf(c.Context(), "CreateRoom Handle BodyParser error: %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "Invalid request body"},
		)
	}

	room, err := h.roomService.Create(c.Context(), rooms.Room{
		ID:      uuid.New().String(),
		Name:    req.Name,
		OwnerID: userID,
	})

	if err != nil {
		logger.Errorf(c.Context(), "CreateRoom Handle Create error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to create room"},
		)
	}

	_, err = h.participantService.Add(c.Context(), rooms.RoomParticipant{
		ID:     uuid.New().String(),
		RoomID: room.ID,
		UserID: userID,
		Role:   "owner",
	})
	if err != nil {
		logger.Errorf(c.Context(), "CreateRoom Handle AddParticipant error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to add participant to room"},
		)
	}

	response := CreateRoomResponse{
		ID:      room.ID,
		Name:    room.Name,
		OwnerID: room.OwnerID,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
