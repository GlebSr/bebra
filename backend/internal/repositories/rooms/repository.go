package rooms

import (
	"context"
	"database/sql"
	"errors"

	entitiesrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/rooms/gen"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type RoomRepository interface {
	Create(context.Context, CreateParams) (entitiesrooms.Room, error)
	GetByID(context.Context, uuid.UUID) (entitiesrooms.Room, error)
	GetAllForUser(context.Context, uuid.UUID) ([]entitiesrooms.Room, error)
	Update(context.Context, UpdateParams) (entitiesrooms.Room, error)
	Delete(context.Context, uuid.UUID) error
}

type Repository struct {
	db *gen.Queries
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: gen.New(db)}
}

type CreateParams struct {
	ID      uuid.UUID
	Name    string
	OwnerID uuid.UUID
}

func (r *Repository) Create(ctx context.Context, params CreateParams) (entitiesrooms.Room, error) {
	created, err := r.db.Create(ctx, gen.CreateParams{
		ID:      params.ID,
		Name:    params.Name,
		OwnerID: params.OwnerID,
	})
	if err != nil {
		logger.Errorf(ctx, "CreateRoom error: %v; data: %v", err, params)
		return entitiesrooms.Room{}, err
	}

	return entitiesrooms.Room{
		ID:        created.ID.String(),
		Name:      created.Name,
		OwnerID:   created.OwnerID.String(),
		CreatedAt: created.CreatedAt.Time,
	}, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (entitiesrooms.Room, error) {
	res, err := r.db.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return entitiesrooms.Room{}, nil
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(ctx, "GetRoomByID error: %v; id: %v", err, id)
		return entitiesrooms.Room{}, err
	}

	return entitiesrooms.Room{
		ID:        res.ID.String(),
		Name:      res.Name,
		OwnerID:   res.OwnerID.String(),
		CreatedAt: res.CreatedAt.Time,
	}, nil
}

func (r *Repository) GetAllForUser(ctx context.Context, userID uuid.UUID) ([]entitiesrooms.Room, error) {
	items, err := r.db.GetAllForUser(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "GetAllRoomsForUser error: %v; userID: %v", err, userID)
		return nil, err
	}

	res := make([]entitiesrooms.Room, 0, len(items))
	for _, it := range items {
		res = append(res, entitiesrooms.Room{
			ID:        it.ID.String(),
			Name:      it.Name,
			OwnerID:   it.OwnerID.String(),
			CreatedAt: it.CreatedAt.Time,
		})
	}
	return res, nil
}

type UpdateParams struct {
	ID   uuid.UUID
	Name string
}

func (r *Repository) Update(ctx context.Context, params UpdateParams) (entitiesrooms.Room, error) {
	updatedRoom, err := r.db.Update(ctx, gen.UpdateParams{
		ID:   params.ID,
		Name: params.Name,
	})
	if err != nil {
		logger.Errorf(ctx, "UpdateRoom error: %v; data: %v", err, params)
		return entitiesrooms.Room{}, err
	}

	return entitiesrooms.Room{
		ID:        updatedRoom.ID.String(),
		Name:      updatedRoom.Name,
		OwnerID:   updatedRoom.OwnerID.String(),
		CreatedAt: updatedRoom.CreatedAt.Time,
	}, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.Delete(ctx, id)
	if err != nil {
		logger.Errorf(ctx, "DeleteRoom error: %v; id: %v", err, id)
		return err
	}
	return nil
}
