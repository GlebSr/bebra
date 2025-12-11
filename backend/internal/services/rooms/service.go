package rooms

import (
	"context"

	entitiesrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/hub"
	repositoryrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type RoomService interface {
	Create(context.Context, entitiesrooms.Room) (entitiesrooms.Room, error)
	GetByID(context.Context, string) (entitiesrooms.Room, error)
	GetAllForUser(context.Context, string) ([]entitiesrooms.Room, error)
	Update(context.Context, entitiesrooms.Room) (entitiesrooms.Room, error)
	Delete(context.Context, string) error
}

type Service struct {
	repo repositoryrooms.RoomRepository
	hub  hub.Hub
}

func NewService(repo repositoryrooms.RoomRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetHub(h hub.Hub) {
	s.hub = h
}

func (s *Service) Create(ctx context.Context, room entitiesrooms.Room) (entitiesrooms.Room, error) {
	id, err := uuid.Parse(room.ID)
	if err != nil {
		logger.Errorf(ctx, "CreateRoom invalid ID: %v", err)

		return entitiesrooms.Room{}, err
	}

	ownerID, err := uuid.Parse(room.OwnerID)
	if err != nil {
		logger.Errorf(ctx, "CreateRoom invalid OwnerID: %v", err)

		return entitiesrooms.Room{}, err
	}

	return s.repo.Create(ctx, repositoryrooms.CreateParams{
		ID:      id,
		Name:    room.Name,
		OwnerID: ownerID,
	})
}

func (s *Service) GetByID(ctx context.Context, id string) (entitiesrooms.Room, error) {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorf(ctx, "GetByID invalid ID: %v", err)

		return entitiesrooms.Room{}, err
	}

	return s.repo.GetByID(ctx, uuidId)
}

func (s *Service) GetAllForUser(ctx context.Context, userID string) ([]entitiesrooms.Room, error) {
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf(ctx, "GetAllForUser invalid UserID: %v", err)

		return nil, err
	}

	return s.repo.GetAllForUser(ctx, uuidUserID)
}

func (s *Service) Update(ctx context.Context, room entitiesrooms.Room) (entitiesrooms.Room, error) {
	id, err := uuid.Parse(room.ID)
	if err != nil {
		logger.Errorf(ctx, "UpdateRoom invalid ID: %v", err)

		return entitiesrooms.Room{}, err
	}

	params := repositoryrooms.UpdateParams{
		ID:   id,
		Name: room.Name,
	}

	result, err := s.repo.Update(ctx, params)
	if err == nil && s.hub != nil {
		s.hub.Broadcast(room.ID, hub.RoomEvent{
			Type:   hub.EventRoomUpdated,
			RoomID: room.ID,
			Payload: map[string]any{
				"id":   result.ID,
				"name": result.Name,
			},
		})
	}
	return result, err
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorf(ctx, "DeleteRoom invalid ID: %v", err)

		return err
	}

	return s.repo.Delete(ctx, uuidId)
}
