package results

import (
	"context"

	entitiesrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	repositoryresults "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/results"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type ResultService interface {
	PickResult(context.Context, string) (string, error)
	GetLastResult(context.Context, string) (string, error)
	GetAllResults(context.Context, string) ([]string, error)
	Delete(context.Context, string) error
	Add(context.Context, entitiesrooms.Result) (entitiesrooms.Result, error)
}

type Service struct {
	repo repositoryresults.ResultRepository
}

func NewService(repo repositoryresults.ResultRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) PickResult(ctx context.Context, roomID string) (string, error) {
	uuidRoomID, err := uuid.Parse(roomID)
	if err != nil {
		logger.Errorf(ctx, "PickResult invalid RoomID: %v", err)

		return "", err
	}

	return s.repo.PickResult(ctx, uuidRoomID)
}

func (s *Service) GetLastResult(ctx context.Context, roomID string) (string, error) {
	uuidRoomID, err := uuid.Parse(roomID)
	if err != nil {
		logger.Errorf(ctx, "GetLastResult invalid RoomID: %v", err)

		return "", err
	}

	res, err := s.repo.GetLastResult(ctx, uuidRoomID)
	if err != nil {
		return "", err
	}
	return res.GameID, nil
}

func (s *Service) GetAllResults(ctx context.Context, roomID string) ([]string, error) {
	uuidRoomID, err := uuid.Parse(roomID)
	if err != nil {
		logger.Errorf(ctx, "GetAllResults invalid RoomID: %v", err)

		return nil, err
	}

	results, err := s.repo.GetAllResults(ctx, uuidRoomID)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(results))
	for _, result := range results {
		res = append(res, result.GameID)
	}
	return res, nil
}

func (s *Service) Delete(ctx context.Context, roomID string) error {
	uuidRoomID, err := uuid.Parse(roomID)
	if err != nil {
		logger.Errorf(ctx, "DeleteResults invalid RoomID: %v", err)

		return err
	}

	return s.repo.Delete(ctx, uuidRoomID)
}

func (s *Service) Add(ctx context.Context, result entitiesrooms.Result) (entitiesrooms.Result, error) {
	id, err := uuid.Parse(result.ID)
	if err != nil {
		logger.Errorf(ctx, "AddResult invalid ID: %v", err)

		return entitiesrooms.Result{}, err
	}
	gameID, err := uuid.Parse(result.GameID)
	if err != nil {
		logger.Errorf(ctx, "AddResult invalid GameID: %v", err)

		return entitiesrooms.Result{}, err
	}
	roomID, err := uuid.Parse(result.RoomID)
	if err != nil {
		logger.Errorf(ctx, "AddResult invalid RoomID: %v", err)

		return entitiesrooms.Result{}, err
	}
	chosenBy, err := uuid.Parse(result.ChosenBy)
	if err != nil {
		logger.Errorf(ctx, "AddResult invalid ChosenBy: %v", err)

		return entitiesrooms.Result{}, err
	}
	return s.repo.Add(ctx, repositoryresults.AddParams{
		ID:       id,
		GameID:   gameID,
		RoomID:   roomID,
		ChosenBy: chosenBy,
	})
}
