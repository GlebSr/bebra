package users

import (
	"context"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/profile"
	repositoryusers "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/users"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type UserService interface {
	Add(context.Context, profile.User) (profile.User, error)
	GetByID(context.Context, string) (profile.User, error)
	GetByName(context.Context, string) (profile.User, error)
	Update(context.Context, profile.User) (profile.User, error)
	Delete(context.Context, string) error
}

type Service struct {
	repo repositoryusers.UserRepository
}

func NewService(repo repositoryusers.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(ctx context.Context, user profile.User) (profile.User, error) {
	id, err := uuid.Parse(user.ID)
	if err != nil {
		logger.Errorf(ctx, "Invalid UUID: %v", err)

		return profile.User{}, err
	}

	params := repositoryusers.AddParams{
		ID:           id,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
	}
	return s.repo.Add(ctx, params)
}

func (s *Service) GetByID(ctx context.Context, id string) (profile.User, error) {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorf(ctx, "GetByID invalid ID: %v", err)

		return profile.User{}, err
	}

	return s.repo.GetByID(ctx, uuidId)
}

func (s *Service) Update(ctx context.Context, user profile.User) (profile.User, error) {
	id, err := uuid.Parse(user.ID)
	if err != nil {
		logger.Errorf(ctx, "UpdateUser invalid ID: %v", err)

		return profile.User{}, err
	}

	params := repositoryusers.UpdateParams{
		ID:           id,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
	}
	return s.repo.Update(ctx, params)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorf(ctx, "DeleteUser invalid ID: %v", err)

		return err
	}

	return s.repo.Delete(ctx, uuidId)
}

func (s *Service) GetByName(ctx context.Context, name string) (profile.User, error) {
	return s.repo.GetByName(ctx, name)
}
