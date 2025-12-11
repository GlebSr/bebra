package tokens

import (
	"context"
	"errors"
	"time"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/config"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/entities/profile"
	repositoryrefreshtokens "code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/repositories/refresh_tokens"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateJWTToken(ctx context.Context, userID string) (string, error)
	ValidateJWTToken(ctx context.Context, token string) (string, error)
	GenerateSSEToken(ctx context.Context, userID, jobID string) (string, error)
	ValidateSSEToken(ctx context.Context, token string) (string, string, error) //TODO: from 3 return to struct
	CreateRefreshToken(ctx context.Context, userID string) (profile.RefreshToken, error)
	ValidateRefreshToken(ctx context.Context, token string) (string, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	GetRefreshTokensByUserID(ctx context.Context, userID string) ([]profile.RefreshToken, error)
	ClearRefreshTokens(ctx context.Context) error
}

type Service struct {
	conf              config.AppConfig
	refreshTokensRepo repositoryrefreshtokens.RefreshTokenRepository
}

func NewService(
	conf config.AppConfig,
	refreshTokensRepo repositoryrefreshtokens.RefreshTokenRepository,
) *Service {
	return &Service{
		conf:              conf,
		refreshTokensRepo: refreshTokensRepo,
	}
}

func (s *Service) GenerateJWTToken(ctx context.Context, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(s.conf.GetJWTConfig().TokenTTL) * time.Second).Unix(),
	})

	return token.SignedString([]byte(s.conf.GetJWTConfig().SecretKey))
}

func (s *Service) GenerateSSEToken(ctx context.Context, userID, jobID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"job_id":  jobID,
		"exp":     time.Now().Add(2 * time.Minute).Unix(),
	})

	return token.SignedString([]byte(s.conf.GetJWTConfig().SecretKey))
}

func (s *Service) ValidateJWTToken(ctx context.Context, token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.conf.GetJWTConfig().SecretKey), nil
	})

	if err != nil {
		logger.Errorf(ctx, "Parse jwt token error: %v", err)

		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims["user_id"].(string), nil
	}

	return "", errors.New("invalid token")
}

func (s *Service) ValidateSSEToken(ctx context.Context, token string) (string, string, error) { //TODO: from 3 return to struct
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.conf.GetJWTConfig().SecretKey), nil
	})

	if err != nil {
		logger.Errorf(ctx, "Parse sse token error: %v", err)

		return "", "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims["user_id"].(string), claims["job_id"].(string), nil
	}

	return "", "", errors.New("invalid token")
}

func (s *Service) CreateRefreshToken(ctx context.Context, userID string) (profile.RefreshToken, error) {
	token := profile.RefreshToken{
		UserID:    userID,
		Token:     uuid.New().String(), //TODO: change to more secure token generation
		ExpiresAt: time.Now().Add(time.Duration(s.conf.GetJWTConfig().TokenTTL) * time.Second),
	}

	token, err := s.refreshTokensRepo.Create(ctx, repositoryrefreshtokens.CreateParams{
		UserID:    uuid.MustParse(token.UserID),
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
	})
	if err != nil {
		logger.Errorf(ctx, "Create refresh token error: %v", err)

		return profile.RefreshToken{}, err
	}

	return token, nil
}

func (s *Service) ClearRefreshTokens(ctx context.Context) error {
	err := s.refreshTokensRepo.Clear(ctx)
	if err != nil {
		logger.Errorf(ctx, "Clear refresh tokens error: %v", err)

		return err
	}

	return nil
}

func (s *Service) DeleteRefreshToken(ctx context.Context, token string) error {
	err := s.refreshTokensRepo.Delete(ctx, token)
	if err != nil {
		logger.Errorf(ctx, "Delete refresh token error: %v", err)

		return err
	}

	return nil
}

func (s *Service) ValidateRefreshToken(ctx context.Context, token string) (string, error) {
	storedToken, err := s.refreshTokensRepo.Get(ctx, token)
	if err != nil {
		logger.Errorf(ctx, "Get refresh token error: %v", err)

		return "", err
	}

	if storedToken.ExpiresAt.Before(time.Now()) {
		err := s.ClearRefreshTokens(ctx)
		if err != nil {
			logger.Errorf(ctx, "Clear refresh tokens error: %v", err)
		}

		return "", errors.New("token expired") //TODO: custom error
	}

	err = s.DeleteRefreshToken(ctx, token)
	if err != nil {
		logger.Errorf(ctx, "Delete refresh token error: %v", err)
	}

	return storedToken.UserID, nil
}

func (s *Service) GetRefreshTokensByUserID(ctx context.Context, userID string) ([]profile.RefreshToken, error) {
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf(ctx, "Get refresh tokens by user id invalid ID: %v", err)

		return nil, err
	}

	tokens, err := s.refreshTokensRepo.GetByUserID(ctx, uuidUserID)
	if err != nil {
		logger.Errorf(ctx, "Get refresh tokens by user id error: %v", err)

		return nil, err
	}

	return tokens, nil
}
