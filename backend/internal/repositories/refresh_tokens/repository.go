package refreshtokens

import (
	"context"
	"database/sql"
	"time"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/profile"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/refresh_tokens/gen"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/google/uuid"
)

type CreateParams struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}

type RefreshTokenRepository interface {
	Create(context.Context, CreateParams) (profile.RefreshToken, error)
	Get(context.Context, string) (profile.RefreshToken, error)
	Delete(context.Context, string) error
	GetByUserID(context.Context, uuid.UUID) ([]profile.RefreshToken, error)
	Clear(context.Context) error
}

type Repository struct {
	db *gen.Queries
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: gen.New(db)}
}

func (r *Repository) Create(ctx context.Context, params CreateParams) (profile.RefreshToken, error) {
	createdToken, err := r.db.Create(ctx, gen.CreateParams{
		Token:     params.Token,
		UserID:    params.UserID,
		ExpiresAt: params.ExpiresAt,
	})
	if err != nil {
		logger.Errorf(ctx, "CreateRefreshToken error: %v; data: %v", err, params)

		return profile.RefreshToken{}, err
	}

	return profile.RefreshToken{
		Token:     createdToken.Token,
		UserID:    createdToken.UserID.String(),
		ExpiresAt: createdToken.ExpiresAt,
		CreatedAt: createdToken.CreatedAt.Time,
	}, nil
}

func (r *Repository) Get(ctx context.Context, token string) (profile.RefreshToken, error) {
	refreshToken, err := r.db.Get(ctx, token)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorf(ctx, "GetRefreshToken error: %v; token: %v", err, token)

		return profile.RefreshToken{}, err
	}

	return profile.RefreshToken{
		Token:     refreshToken.Token,
		UserID:    refreshToken.UserID.String(),
		ExpiresAt: refreshToken.ExpiresAt,
		CreatedAt: refreshToken.CreatedAt.Time,
	}, nil
}

func (r *Repository) Delete(ctx context.Context, token string) error {
	err := r.db.Delete(ctx, token)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorf(ctx, "DeleteRefreshToken error: %v; token: %v", err, token)

		return err
	}

	return nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]profile.RefreshToken, error) {
	tokens, err := r.db.GetByUserID(ctx, userID)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorf(ctx, "GetRefreshTokensByUserID error: %v; userID: %v", err, userID)

		return nil, err
	}

	result := make([]profile.RefreshToken, len(tokens))
	for i, token := range tokens {
		result[i] = profile.RefreshToken{
			Token:     token.Token,
			UserID:    token.UserID.String(),
			ExpiresAt: token.ExpiresAt,
			CreatedAt: token.CreatedAt.Time,
		}
	}
	return result, nil
}

func (r *Repository) Clear(ctx context.Context) error {
	err := r.db.Clear(ctx)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorf(ctx, "ClearRefreshTokens error: %v", err)

		return err
	}

	return nil
}
