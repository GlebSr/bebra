package participants

import (
	"context"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/profile"
	entitiesrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/hub"
	repositoryparticipants "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/participants"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type ParticipantService interface {
	Add(context.Context, entitiesrooms.RoomParticipant) (entitiesrooms.RoomParticipant, error)
	GetAllParticipants(context.Context, string) ([]profile.User, []string, error)
	Delete(context.Context, string, string) error
	Get(context.Context, string, string) (entitiesrooms.RoomParticipant, error)
}

type Service struct {
	repo repositoryparticipants.ParticipantRepository
	hub  hub.Hub
}

func NewService(repo repositoryparticipants.ParticipantRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetHub(h hub.Hub) {
	s.hub = h
}

func (s *Service) Add(ctx context.Context, participant entitiesrooms.RoomParticipant) (entitiesrooms.RoomParticipant, error) {
	id, err := uuid.Parse(participant.ID)
	if err != nil {
		logger.Errorf(ctx, "AddParticipant invalid ID: %v", err)

		return entitiesrooms.RoomParticipant{}, err
	}

	roomID, err := uuid.Parse(participant.RoomID)
	if err != nil {
		logger.Errorf(ctx, "AddParticipant invalid RoomID: %v", err)

		return entitiesrooms.RoomParticipant{}, err
	}

	userID, err := uuid.Parse(participant.UserID)
	if err != nil {
		logger.Errorf(ctx, "AddParticipant invalid UserID: %v", err)

		return entitiesrooms.RoomParticipant{}, err
	}

	result, err := s.repo.Add(ctx, repositoryparticipants.AddParams{
		ID:     id,
		RoomID: roomID,
		UserID: userID,
		Role:   participant.Role,
	})
	if err == nil && s.hub != nil {
		s.hub.Broadcast(participant.RoomID, hub.RoomEvent{
			Type:   hub.EventParticipantAdded,
			RoomID: participant.RoomID,
			Payload: map[string]any{
				"id":      result.ID,
				"user_id": result.UserID,
				"role":    result.Role,
			},
		})
	}
	return result, err
}

func (s *Service) GetAllParticipants(ctx context.Context, roomID string) ([]profile.User, []string, error) {
	uuidRoomID, err := uuid.Parse(roomID)
	if err != nil {
		logger.Errorf(ctx, "GetAllParticipants invalid RoomID: %v", err)

		return nil, nil, err
	}
	list, err := s.repo.GetAllParticipants(ctx, uuidRoomID)
	if err != nil {
		return nil, nil, err
	}

	res := make([]profile.User, 0, len(list))
	roles := make([]string, 0, len(list))
	for _, it := range list {
		res = append(res, it.User)
		roles = append(roles, it.Role)
	}

	return res, roles, nil
}

func (s *Service) Delete(ctx context.Context, roomID, userID string) error {
	uuidRoomID, err := uuid.Parse(roomID)
	if err != nil {
		logger.Errorf(ctx, "DeleteParticipant invalid RoomID: %v", err)

		return err
	}
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf(ctx, "DeleteParticipant invalid UserID: %v", err)

		return err
	}
	err = s.repo.Delete(ctx, uuidRoomID, uuidUserID)
	if err == nil && s.hub != nil {
		s.hub.Broadcast(roomID, hub.RoomEvent{
			Type:   hub.EventParticipantLeft,
			RoomID: roomID,
			Payload: map[string]any{
				"user_id": userID,
			},
		})
	}
	return err
}

func (s *Service) Get(ctx context.Context, roomID, userID string) (entitiesrooms.RoomParticipant, error) {
	uuidRoomID, err := uuid.Parse(roomID)
	if err != nil {
		logger.Errorf(ctx, "GetParticipant invalid RoomID: %v", err)

		return entitiesrooms.RoomParticipant{}, err
	}
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf(ctx, "GetParticipant invalid UserID: %v", err)

		return entitiesrooms.RoomParticipant{}, err
	}

	return s.repo.Get(ctx, uuidRoomID, uuidUserID)
}
