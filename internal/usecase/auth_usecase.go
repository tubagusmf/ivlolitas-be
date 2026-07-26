package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tubagusmf/ivlolitas-be/internal/jwt"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type IAuthUsecase interface {
	Login(ctx context.Context, req *model.LoginInput) (*model.LoginResponse, error)
	RefreshToken(ctx context.Context, token string) (*model.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	LogoutAll(ctx context.Context, userID string) error
}

type authUsecase struct {
	userRepo         repository.IUserRepository
	refreshTokenRepo repository.IRefreshTokenRepository
	jwt              *jwt.JWT
}

func NewAuthUsecase(userRepo repository.IUserRepository, refreshTokenRepo repository.IRefreshTokenRepository, jwt *jwt.JWT) IAuthUsecase {
	return &authUsecase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwt:              jwt,
	}
}

func (a *authUsecase) Login(ctx context.Context, req *model.LoginInput) (*model.LoginResponse, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("email or password invalid")
	}

	if !user.IsActive {
		return nil, errors.New("user inactive")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("email or password invalid")
	}

	accessToken, err := a.jwt.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refresh := &model.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	if err := a.refreshTokenRepo.Create(ctx, refresh); err != nil {
		return nil, err
	}

	now := time.Now()
	user.LastLogin = &now

	_, _ = a.userRepo.UpdateUser(ctx, user)

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &model.Auth{
			UserID:   user.ID,
			RoleID:   user.RoleID,
			Role:     user.Role,
			Email:    user.Email,
			FullName: user.FullName,
		},
	}, nil
}

func (a *authUsecase) RefreshToken(ctx context.Context, token string) (*model.LoginResponse, error) {
	refresh, err := a.refreshTokenRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if refresh.Revoked {
		return nil, errors.New("refresh token revoked")
	}

	if time.Now().After(refresh.ExpiredAt) {
		return nil, errors.New("refresh token expired")
	}

	user, err := a.userRepo.GetUserByID(ctx, refresh.UserID)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.jwt.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := a.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := a.refreshTokenRepo.Revoke(ctx, token); err != nil {
		return nil, err
	}

	err = a.refreshTokenRepo.Create(ctx, &model.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour),
	})

	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: &model.Auth{
			UserID:   user.ID,
			RoleID:   user.RoleID,
			Role:     user.Role,
			Email:    user.Email,
			FullName: user.FullName,
		},
	}, nil
}

func (a *authUsecase) Logout(ctx context.Context, token string) error {
	return a.refreshTokenRepo.Revoke(ctx, token)
}

func (a *authUsecase) LogoutAll(ctx context.Context, userID string) error {
	return a.refreshTokenRepo.RevokeAllByUserID(ctx, userID)
}
