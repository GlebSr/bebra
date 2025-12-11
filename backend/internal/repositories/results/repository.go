package results

import (
	"context"
	"database/sql"
	"errors"

	entitiesrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/results/gen"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type ResultRepository interface {
	PickResult(context.Context, uuid.UUID) (string, error)
	GetLastResult(context.Context, uuid.UUID) (entitiesrooms.Result, error)
	GetAllResults(context.Context, uuid.UUID) ([]entitiesrooms.Result, error)
	Delete(context.Context, uuid.UUID) error
	Add(context.Context, AddParams) (entitiesrooms.Result, error)
}

type Repository struct {
	db *gen.Queries
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: gen.New(db)}
}

func (r *Repository) PickResult(ctx context.Context, roomID uuid.UUID) (string, error) {
	gameID, err := r.db.PickResult(ctx, roomID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(ctx, "PickResult error: %v; roomID: %v", err, roomID)

		return "", err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	return gameID.String(), nil
}

func (r *Repository) GetLastResult(ctx context.Context, roomID uuid.UUID) (entitiesrooms.Result, error) {
	res, err := r.db.GetLastResult(ctx, roomID)
	if errors.Is(err, sql.ErrNoRows) {
		return entitiesrooms.Result{}, nil
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(ctx, "GetLastResult error: %v; roomID: %v", err, roomID)

		return entitiesrooms.Result{}, err
	}

	return entitiesrooms.Result{
		ID:        res.ID.String(),
		RoomID:    res.RoomID.String(),
		GameID:    res.GameID.String(),
		ChosenBy:  res.ChosenBy.String(),
		CreatedAt: res.CreatedAt.Time,
	}, nil
}

func (r *Repository) GetAllResults(ctx context.Context, roomID uuid.UUID) ([]entitiesrooms.Result, error) {
	items, err := r.db.GetAllResults(ctx, roomID)
	if err != nil {
		logger.Errorf(ctx, "GetAllResults error: %v; roomID: %v", err, roomID)

		return nil, err
	}

	res := make([]entitiesrooms.Result, 0, len(items))
	for _, it := range items {
		res = append(res, entitiesrooms.Result{
			ID:        it.ID.String(),
			RoomID:    it.RoomID.String(),
			GameID:    it.GameID.String(),
			ChosenBy:  it.ChosenBy.String(),
			CreatedAt: it.CreatedAt.Time,
		})
	}

	return res, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.Delete(ctx, id)
	if err != nil {
		logger.Errorf(ctx, "DeleteResult error: %v; id: %v", err, id)

		return err
	}

	return nil
}

type AddParams struct {
	ID       uuid.UUID
	RoomID   uuid.UUID
	GameID   uuid.UUID
	ChosenBy uuid.UUID
}

func (r *Repository) Add(ctx context.Context, params AddParams) (entitiesrooms.Result, error) {
	result, err := r.db.Add(ctx, gen.AddParams{
		ID:       params.ID,
		RoomID:   params.RoomID,
		GameID:   params.GameID,
		ChosenBy: params.ChosenBy,
	})
	if err != nil {
		logger.Errorf(ctx, "AddResult error: %v; data: %v", err, params)

		return entitiesrooms.Result{}, err
	}

	return entitiesrooms.Result{
		ID:        result.ID.String(),
		RoomID:    result.RoomID.String(),
		GameID:    result.GameID.String(),
		ChosenBy:  result.ChosenBy.String(),
		CreatedAt: result.CreatedAt.Time,
	}, nil
}
