package games

import (
	"context"
	"database/sql"
	"errors"

	entitiesrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/games/gen"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type GameRepository interface {
	Add(context.Context, AddParams) (entitiesrooms.Game, error)
	GetAllRoomGames(context.Context, uuid.UUID) ([]entitiesrooms.Game, error)
	Delete(context.Context, uuid.UUID) error
	Get(context.Context, uuid.UUID) (entitiesrooms.Game, error)
}

type Repository struct {
	db *gen.Queries
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: gen.New(db)}
}

type AddParams struct {
	ID     uuid.UUID
	RoomID uuid.UUID
	Title  string
}

func (r *Repository) Add(ctx context.Context, params AddParams) (entitiesrooms.Game, error) {
	created, err := r.db.Add(ctx, gen.AddParams{
		ID:     params.ID,
		RoomID: params.RoomID,
		Title:  params.Title,
	})
	if err != nil {
		logger.Errorf(ctx, "AddGame error: %v; data: %v", err, params)

		return entitiesrooms.Game{}, err
	}

	return entitiesrooms.Game{
		ID:        created.ID.String(),
		RoomID:    created.RoomID.String(),
		Title:     created.Title,
		CreatedAt: created.CreatedAt.Time,
	}, nil
}

func (r *Repository) GetAllRoomGames(ctx context.Context, roomID uuid.UUID) ([]entitiesrooms.Game, error) {
	items, err := r.db.GetAllRoomGames(ctx, roomID)
	if err != nil {
		logger.Errorf(ctx, "GetAllRoomGames error: %v; roomID: %v", err, roomID)

		return nil, err
	}

	res := make([]entitiesrooms.Game, 0, len(items))
	for _, it := range items {
		res = append(res, entitiesrooms.Game{
			ID:        it.ID.String(),
			RoomID:    it.RoomID.String(),
			Title:     it.Title,
			CreatedAt: it.CreatedAt.Time,
		})
	}

	return res, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.Delete(ctx, id)
	if err != nil {
		logger.Errorf(ctx, "DeleteGame error: %v; id: %v", err, id)
		return err
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (entitiesrooms.Game, error) {
	item, err := r.db.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return entitiesrooms.Game{}, nil
	}

	if err != nil {
		logger.Errorf(ctx, "GetGame error: %v; id: %v", err, id)
		return entitiesrooms.Game{}, err
	}

	return entitiesrooms.Game{
		ID:        item.ID.String(),
		RoomID:    item.RoomID.String(),
		Title:     item.Title,
		CreatedAt: item.CreatedAt.Time,
	}, nil
}
