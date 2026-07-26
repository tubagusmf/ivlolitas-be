package usecase

import (
	"context"

	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/repository"
)

type IRoleUsecase interface {
	GetRoles(ctx context.Context) ([]*model.Role, error)
	GetRoleByID(ctx context.Context, id int64) (*model.Role, error)
	CreateRole(ctx context.Context, in *model.CreateRoleInput) (*model.Role, error)
	UpdateRole(ctx context.Context, id int64, in *model.UpdateRoleInput) (*model.Role, error)
	DeleteRole(ctx context.Context, id int64) error
}

type roleUsecase struct {
	repo repository.IRoleRepository
}

func NewRoleUsecase(repo repository.IRoleRepository) IRoleUsecase {
	return &roleUsecase{repo: repo}
}

func (r *roleUsecase) GetRoles(ctx context.Context) ([]*model.Role, error) {
	return r.repo.GetRoles(ctx)
}

func (r *roleUsecase) GetRoleByID(ctx context.Context, id int64) (*model.Role, error) {
	return r.repo.GetRoleByID(ctx, id)
}

func (r *roleUsecase) CreateRole(ctx context.Context, in *model.CreateRoleInput) (*model.Role, error) {
	return r.repo.CreateRole(ctx, &model.Role{
		Name:        in.Name,
		Description: in.Description,
	})
}

func (r *roleUsecase) UpdateRole(ctx context.Context, id int64, in *model.UpdateRoleInput) (*model.Role, error) {
	return r.repo.UpdateRole(ctx, &model.Role{
		ID:          id,
		Name:        in.Name,
		Description: in.Description,
	})
}

func (r *roleUsecase) DeleteRole(ctx context.Context, id int64) error {
	return r.repo.DeleteRole(ctx, id)
}
