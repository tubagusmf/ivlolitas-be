package repository

import (
	"context"

	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	GetRoles(ctx context.Context) ([]*model.Role, error)
	GetRoleByID(ctx context.Context, id int64) (*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role) (*model.Role, error)
	UpdateRole(ctx context.Context, role *model.Role) (*model.Role, error)
	DeleteRole(ctx context.Context, id int64) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetRoles(ctx context.Context) ([]*model.Role, error) {
	var roles []*model.Role

	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *roleRepository) GetRoleByID(ctx context.Context, id int64) (*model.Role, error) {
	var role model.Role

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) CreateRole(ctx context.Context, role *model.Role) (*model.Role, error) {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) UpdateRole(ctx context.Context, role *model.Role) (*model.Role, error) {
	if err := r.db.WithContext(ctx).Save(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) DeleteRole(ctx context.Context, id int64) error {
	role := &model.Role{ID: id}
	return r.db.WithContext(ctx).Delete(role).Error
}
