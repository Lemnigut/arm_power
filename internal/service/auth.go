package service

import (
	"context"
	"errors"
	"time"

	"arm_back/internal/model"
	"arm_back/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: []byte(jwtSecret)}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.TokenResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Username:     req.Username,
		PasswordHash: string(hash),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, model.ErrConflict
	}

	return s.generateTokens(user.ID)
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.TokenResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, model.ErrUnauthorized
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, model.ErrUnauthorized
	}

	return s.generateTokens(user.ID)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.TokenResponse, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, model.ErrUnauthorized
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, model.ErrUnauthorized
	}

	// Verify user still exists
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return nil, model.ErrUnauthorized
	}

	return s.generateTokens(userID)
}

func (s *AuthService) ValidateToken(tokenStr string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, model.ErrUnauthorized
	}

	return uuid.Parse(claims.Subject)
}

func (s *AuthService) generateTokens(userID uuid.UUID) (*model.TokenResponse, error) {
	now := time.Now()

	accessClaims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(now.Add(90 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
