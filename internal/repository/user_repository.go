package repository

import (
	"context"

	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUsers(ctx context.Context, email string) ([]*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) GetUsers(ctx context.Context, email string) ([]*model.User, error) {
	var users []*model.User

	db := u.db.WithContext(ctx).Preload("Role")

	if email != "" {
		db = db.Where("email = ?", email)
	}

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	err := u.db.
		WithContext(ctx).
		Preload("Role").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := u.db.
		WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := u.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) DeleteUser(ctx context.Context, id string) error {
	user := &model.User{
		ID: id,
	}

	return u.db.WithContext(ctx).Delete(user).Error
}
