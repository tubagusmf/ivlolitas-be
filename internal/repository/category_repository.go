package repository

import (
	"context"

	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"gorm.io/gorm"
)

type ICategoryRepository interface {
	GetCategories(ctx context.Context, filter *model.CategoryFilter) ([]*model.Category, error)
	GetCategoryByID(ctx context.Context, id int64) (*model.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetCategories(ctx context.Context, filter *model.CategoryFilter) ([]*model.Category, error) {
	var categories []*model.Category

	db := r.db.WithContext(ctx)

	if filter.Search != "" {
		db = db.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	if filter.IsActive != nil {
		db = db.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Sort != "" {
		db = db.Order(buildCategorySort(filter.Sort))
	} else {
		db = db.Order("created_at DESC")
	}

	err := db.
		Offset((filter.Page - 1) * filter.Limit).
		Limit(filter.Limit).
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int64) (*model.Category, error) {
	var category model.Category

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetCategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	var category model.Category

	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&category).Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	err := r.db.WithContext(ctx).
		Model(category).
		Updates(category).Error

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	category := &model.Category{ID: id}
	return r.db.WithContext(ctx).Delete(category).Error
}

func buildCategorySort(sort string) string {
	switch sort {
	case "name":
		return "name ASC"
	case "-name":
		return "name DESC"
	case "created_at":
		return "created_at ASC"
	case "-created_at":
		return "created_at DESC"
	default:
		return "created_at DESC"
	}
}
