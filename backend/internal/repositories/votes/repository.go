package votes

import (
	"context"
	"database/sql"
	"errors"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/rooms"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/votes/gen"
	"github.com/google/uuid"
)

type VoteRepository interface {
	Add(ctx context.Context, params AddParams) (rooms.Vote, error)
	Get(ctx context.Context, id uuid.UUID) (rooms.Vote, error)
	GetForRoom(ctx context.Context, roomID uuid.UUID) ([]rooms.Vote, error)
	Delete(ctx context.Context, id uuid.UUID) error
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
	GameID uuid.UUID
	UserID uuid.UUID
}

func (r *Repository) Add(ctx context.Context, params AddParams) (rooms.Vote, error) {
	createdVote, err := r.db.Add(ctx, gen.AddParams{
		ID:     params.ID,
		RoomID: params.RoomID,
		GameID: params.GameID,
		UserID: params.UserID,
	})
	if err != nil {
		return rooms.Vote{}, err
	}

	return rooms.Vote{
		ID:        createdVote.ID.String(),
		RoomID:    createdVote.RoomID.String(),
		GameID:    createdVote.GameID.String(),
		UserID:    createdVote.UserID.String(),
		CreatedAt: createdVote.CreatedAt.Time,
	}, nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (rooms.Vote, error) {
	vote, err := r.db.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return rooms.Vote{}, nil
	}

	if err != nil {
		return rooms.Vote{}, err
	}

	return rooms.Vote{
		ID:        vote.ID.String(),
		RoomID:    vote.RoomID.String(),
		GameID:    vote.GameID.String(),
		UserID:    vote.UserID.String(),
		CreatedAt: vote.CreatedAt.Time,
	}, nil
}

func (r *Repository) GetForRoom(ctx context.Context, roomID uuid.UUID) ([]rooms.Vote, error) {
	votes, err := r.db.GetForRoom(ctx, roomID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	result := make([]rooms.Vote, len(votes))
	for i, vote := range votes {
		result[i] = rooms.Vote{
			ID:        vote.ID.String(),
			RoomID:    vote.RoomID.String(),
			GameID:    vote.GameID.String(),
			UserID:    vote.UserID.String(),
			CreatedAt: vote.CreatedAt.Time,
		}
	}

	return result, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Delete(ctx, id)
}
