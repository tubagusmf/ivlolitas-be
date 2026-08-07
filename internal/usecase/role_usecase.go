package usecase

import (
	"context"

	"github.com/sirupsen/logrus"

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
	return &roleUsecase{
		repo: repo,
	}
}

func (r *roleUsecase) GetRoles(ctx context.Context) ([]*model.Role, error) {
	log := logrus.WithField("action", "GetRoles")

	roles, err := r.repo.GetRoles(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("success get %d roles", len(roles))
	return roles, nil
}

func (r *roleUsecase) GetRoleByID(ctx context.Context, id int64) (*model.Role, error) {
	log := logrus.WithFields(logrus.Fields{
		"role_id": id,
	})

	role, err := r.repo.GetRoleByID(ctx, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Info("success get role")
	return role, nil
}

func (r *roleUsecase) CreateRole(ctx context.Context, in *model.CreateRoleInput) (*model.Role, error) {
	log := logrus.WithFields(logrus.Fields{
		"name": in.Name,
		"desc": in.Description,
	})

	role, err := r.repo.CreateRole(ctx, &model.Role{
		Name:        in.Name,
		Description: in.Description,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.WithField("role_id", role.ID).Info("role created successfully")
	return role, nil
}

func (r *roleUsecase) UpdateRole(ctx context.Context, id int64, in *model.UpdateRoleInput) (*model.Role, error) {
	log := logrus.WithFields(logrus.Fields{
		"role_id": id,
		"name":    in.Name,
		"desc":    in.Description,
	})

	role, err := r.repo.UpdateRole(ctx, &model.Role{
		ID:          id,
		Name:        in.Name,
		Description: in.Description,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Info("role updated successfully")
	return role, nil
}

func (r *roleUsecase) DeleteRole(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"role_id": id,
	})

	if err := r.repo.DeleteRole(ctx, id); err != nil {
		log.Error(err)
		return err
	}

	log.Info("role deleted successfully")
	return nil
}
