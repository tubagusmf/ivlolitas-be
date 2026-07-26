package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	GetUsers(ctx context.Context, email string) ([]*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, in *model.CreateUserInput) (*model.User, error)
	UpdateUser(ctx context.Context, id string, in *model.UpdateUserInput) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userUsecase struct {
	repo repository.IUserRepository
}

func NewUserUsecase(repo repository.IUserRepository) IUserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) GetUsers(ctx context.Context, email string) ([]*model.User, error) {
	return u.repo.GetUsers(ctx, email)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *userUsecase) CreateUser(ctx context.Context, in *model.CreateUserInput) (*model.User, error) {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(in.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:          uuid.New().String(),
		FullName:    in.FullName,
		Email:       in.Email,
		Password:    string(hashed),
		Address:     in.Address,
		PhoneNumber: in.PhoneNumber,
		RoleID:      in.RoleID,
	}

	return u.repo.CreateUser(ctx, user)
}

func (u *userUsecase) UpdateUser(ctx context.Context, id string, in *model.UpdateUserInput) (*model.User, error) {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(in.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = in.FullName
	user.Email = in.Email
	user.Password = string(hashed)
	user.Address = in.Address
	user.PhoneNumber = in.PhoneNumber
	user.RoleID = in.RoleID

	return u.repo.UpdateUser(ctx, user)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	return u.repo.DeleteUser(ctx, id)
}
