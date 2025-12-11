package participants

import (
	"context"
	"database/sql"
	"errors"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/profile"
	entitiesrooms "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/participants/gen"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type ParticipantRepository interface {
	Add(context.Context, AddParams) (entitiesrooms.RoomParticipant, error)
	GetAllParticipants(context.Context, uuid.UUID) ([]ParticipantWithUser, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	Get(context.Context, uuid.UUID, uuid.UUID) (entitiesrooms.RoomParticipant, error)
}

type ParticipantWithUser struct {
	User profile.User `json:"user"`
	Role string       `json:"role"`
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
	UserID uuid.UUID
	Role   string
}

func (r *Repository) Add(ctx context.Context, params AddParams) (entitiesrooms.RoomParticipant, error) {
	created, err := r.db.Add(ctx, gen.AddParams{
		ID:     params.ID,
		RoomID: params.RoomID,
		UserID: params.UserID,
		Role:   params.Role,
	})
	if err != nil {
		logger.Errorf(ctx, "AddParticipant error: %v; data: %v", err, params)

		return entitiesrooms.RoomParticipant{}, err
	}

	return entitiesrooms.RoomParticipant{
		ID:        created.ID.String(),
		RoomID:    created.RoomID.String(),
		UserID:    created.UserID.String(),
		Role:      created.Role,
		CreatedAt: created.CreatedAt.Time,
	}, nil
}

func (r *Repository) GetAllParticipants(ctx context.Context, roomID uuid.UUID) ([]ParticipantWithUser, error) {
	items, err := r.db.GetAllParticipants(ctx, roomID)
	if err != nil {
		logger.Errorf(ctx, "GetAllParticipants error: %v; roomID: %v", err, roomID)

		return nil, err
	}

	res := make([]ParticipantWithUser, 0, len(items))
	for _, it := range items {
		res = append(res, ParticipantWithUser{
			User: profile.User{
				ID:   it.ID.String(),
				Name: it.Name,
			},
			Role: it.Role,
		})
	}

	return res, nil
}

func (r *Repository) Delete(ctx context.Context, roomID, userID uuid.UUID) error {
	err := r.db.Delete(ctx, gen.DeleteParams{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(ctx, "DeleteParticipant error: %v; roomID: %v, userID: %v", err, roomID, userID)

		return err
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, roomID, userID uuid.UUID) (entitiesrooms.RoomParticipant, error) {
	participant, err := r.db.Get(ctx, gen.GetParams{
		RoomID: roomID,
		UserID: userID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return entitiesrooms.RoomParticipant{}, nil
	}

	if err != nil {
		logger.Errorf(ctx, "GetParticipant error: %v; roomID: %v, userID: %v", err, roomID, userID)

		return entitiesrooms.RoomParticipant{}, err
	}

	return entitiesrooms.RoomParticipant{
		ID:        participant.ID.String(),
		RoomID:    participant.RoomID.String(),
		UserID:    participant.UserID.String(),
		Role:      participant.Role,
		CreatedAt: participant.CreatedAt.Time,
	}, nil
}
