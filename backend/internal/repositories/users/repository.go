package users

import (
	"context"
	"database/sql"
	"errors"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/profile"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/users/gen"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type UserRepository interface {
	Add(context.Context, AddParams) (profile.User, error)
	GetByID(context.Context, uuid.UUID) (profile.User, error)
	GetByName(context.Context, string) (profile.User, error)
	Update(context.Context, UpdateParams) (profile.User, error)
	Delete(context.Context, uuid.UUID) error
}

type Repository struct {
	db *gen.Queries
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: gen.New(db)}
}

type AddParams struct {
	ID           uuid.UUID
	PasswordHash string
	Name         string
}

func (r *Repository) Add(ctx context.Context, params AddParams) (profile.User, error) {
	createdUser, err := r.db.Add(ctx, gen.AddParams{
		ID:           params.ID,
		PasswordHash: params.PasswordHash,
		Name:         params.Name,
	})
	if err != nil {
		logger.Errorf(ctx, "CreateUser error: %v; data: %v", err, params)

		return profile.User{}, err
	}

	return profile.User{
		ID:           createdUser.ID.String(),
		PasswordHash: createdUser.PasswordHash,
		Name:         createdUser.Name,
		CreatedAt:    createdUser.CreatedAt.Time,
	}, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (profile.User, error) {
	user, err := r.db.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return profile.User{}, nil
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(ctx, "GetUserByID error: %v; id: %v", err, id)

		return profile.User{}, err
	}

	return profile.User{
		ID:           user.ID.String(),
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt.Time,
	}, nil
}

type UpdateParams struct {
	ID           uuid.UUID
	PasswordHash string
	Name         string
}

func (r *Repository) Update(ctx context.Context, params UpdateParams) (profile.User, error) {
	updatedUser, err := r.db.Update(ctx, gen.UpdateParams{
		ID:           params.ID,
		PasswordHash: params.PasswordHash,
		Name:         params.Name,
	})

	if err != nil {
		logger.Errorf(ctx, "UpdateUser error: %v; data: %v", err, params)

		return profile.User{}, err
	}

	return profile.User{
		ID:           updatedUser.ID.String(),
		PasswordHash: updatedUser.PasswordHash,
		Name:         updatedUser.Name,
		CreatedAt:    updatedUser.CreatedAt.Time,
	}, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.Delete(ctx, id)
	if err != nil {
		logger.Errorf(ctx, "DeleteUser error: %v; id: %v", err, id)

		return err
	}

	return nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (profile.User, error) {
	user, err := r.db.GetByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return profile.User{}, nil
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf(ctx, "GetUserByName error: %v; name: %v", err, name)

		return profile.User{}, err
	}

	return profile.User{
		ID:           user.ID.String(),
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt.Time,
	}, nil
}
