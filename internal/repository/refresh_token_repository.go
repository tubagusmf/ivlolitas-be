package repository

import (
	"context"
	"time"

	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"gorm.io/gorm"
)

type IRefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	RevokeAllByUserID(ctx context.Context, userID string) error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) IRefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *refreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var refresh model.RefreshToken

	err := r.db.
		WithContext(ctx).
		Preload("User").
		Preload("User.Role").
		Where("token = ?", token).
		First(&refresh).Error

	if err != nil {
		return nil, err
	}

	return &refresh, nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).Where("token = ?", token).Update("revoked", true).Error
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expired_at < ?", time.Now()).Delete(&model.RefreshToken{}).Error
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}
