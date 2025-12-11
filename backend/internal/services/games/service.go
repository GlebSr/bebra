package games

import (
	"context"

	entitiesrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/hub"
	repositorygames "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/games"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type GameService interface {
	Add(context.Context, entitiesrooms.Game) (entitiesrooms.Game, error)
	GetAllRoomGames(context.Context, string) ([]entitiesrooms.Game, error)
	Delete(context.Context, string, string) error
	Get(context.Context, string) (entitiesrooms.Game, error)
}

type Service struct {
	repo repositorygames.GameRepository
	hub  hub.Hub
}

func NewService(repo repositorygames.GameRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetHub(h hub.Hub) {
	s.hub = h
}

func (s *Service) Add(ctx context.Context, game entitiesrooms.Game) (entitiesrooms.Game, error) {
	id, err := uuid.Parse(game.ID)
	if err != nil {
		logger.Errorf(ctx, "AddGame invalid ID: %v", err)

		return entitiesrooms.Game{}, err
	}

	roomID, err := uuid.Parse(game.RoomID)
	if err != nil {
		logger.Errorf(ctx, "AddGame invalid RoomID: %v", err)

		return entitiesrooms.Game{}, err
	}

	gameRes, err := s.repo.Add(ctx, repositorygames.AddParams{
		ID:     id,
		RoomID: roomID,
		Title:  game.Title,
	})
	if err == nil && s.hub != nil {
		s.hub.Broadcast(game.RoomID, hub.RoomEvent{
			Type:   hub.EventGameAdded,
			RoomID: game.RoomID,
			Payload: map[string]any{
				"id":    gameRes.ID,
				"title": gameRes.Title,
			},
		})
	}
	return gameRes, err
}

func (s *Service) GetAllRoomGames(ctx context.Context, roomID string) ([]entitiesrooms.Game, error) {
	uuidRoomID, err := uuid.Parse(roomID)
	if err != nil {
		logger.Errorf(ctx, "GetAllRoomGames invalid RoomID: %v", err)

		return nil, err
	}

	return s.repo.GetAllRoomGames(ctx, uuidRoomID)
}

func (s *Service) Delete(ctx context.Context, id string, roomID string) error {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorf(ctx, "DeleteGame invalid ID: %v", err)

		return err
	}

	err = s.repo.Delete(ctx, uuidId)
	if err == nil && s.hub != nil {
		s.hub.Broadcast(roomID, hub.RoomEvent{
			Type:   hub.EventGameDeleted,
			RoomID: roomID,
			Payload: map[string]any{
				"id": id,
			},
		})
	}
	return err
}

func (s *Service) Get(ctx context.Context, id string) (entitiesrooms.Game, error) {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorf(ctx, "GetGame invalid ID: %v", err)

		return entitiesrooms.Game{}, err
	}

	return s.repo.Get(ctx, uuidId)
}
