package usecase

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/repository"
	"gorm.io/gorm"
)

type ICategoryUsecase interface {
	GetCategories(ctx context.Context, filter *model.CategoryFilter) ([]*model.Category, error)
	GetCategoryByID(ctx context.Context, id int64) (*model.Category, error)
	CreateCategory(ctx context.Context, in *model.CategoryInput) (*model.Category, error)
	UpdateCategory(ctx context.Context, id int64, in *model.CategoryUpdateInput) (*model.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryUsecase struct {
	repo repository.ICategoryRepository
}

func NewCategoryUsecase(repo repository.ICategoryRepository) ICategoryUsecase {
	return &categoryUsecase{repo: repo}
}

func (c *categoryUsecase) GetCategories(ctx context.Context, filter *model.CategoryFilter) ([]*model.Category, error) {
	log := logrus.WithFields(logrus.Fields{"filter": filter})

	categories, err := c.repo.GetCategories(ctx, filter)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return categories, nil
}

func (c *categoryUsecase) GetCategoryByID(ctx context.Context, id int64) (*model.Category, error) {
	log := logrus.WithFields(logrus.Fields{"id": id})

	category, err := c.repo.GetCategoryByID(ctx, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return category, nil
}

func (c *categoryUsecase) CreateCategory(ctx context.Context, in *model.CategoryInput) (*model.Category, error) {
	log := logrus.WithFields(logrus.Fields{
		"slug": in.Slug,
		"name": in.Name,
	})

	_, err := c.repo.GetCategoryBySlug(ctx, in.Slug)
	if err == nil {
		log.Warn("Category slug already exists")

		return nil, errors.New("category slug already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Error("Failed to check category slug")

		return nil, err
	}

	category := &model.Category{
		Name:        in.Name,
		Slug:        in.Slug,
		Description: in.Description,
		ImageURL:    in.ImageURL,
		IsActive:    in.IsActive,
	}

	result, err := c.repo.CreateCategory(ctx, category)
	if err != nil {
		log.WithError(err).Error("Failed to create category")

		return nil, err
	}

	log.WithField("category_id", result.ID).
		Info("Category created successfully")

	return result, nil
}

func (c *categoryUsecase) UpdateCategory(ctx context.Context, id int64, in *model.CategoryUpdateInput) (*model.Category, error) {
	log := logrus.WithFields(logrus.Fields{
		"category_id": id,
		"slug":        in.Slug,
		"name":        in.Name,
	})

	category, err := c.repo.GetCategoryByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("Failed to get category")

		return nil, err
	}

	if category.Slug != in.Slug {

		existing, err := c.repo.GetCategoryBySlug(ctx, in.Slug)

		if err == nil && existing.ID != category.ID {
			log.Warn("Category slug already exists")

			return nil, errors.New("category slug already exists")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("Failed to check category slug")

			return nil, err
		}
	}

	category.Name = in.Name
	category.Slug = in.Slug
	category.Description = in.Description
	category.ImageURL = in.ImageURL
	category.IsActive = in.IsActive

	result, err := c.repo.UpdateCategory(ctx, category)
	if err != nil {
		log.WithError(err).Error("Failed to update category")

		return nil, err
	}

	log.WithField("category_id", result.ID).
		Info("Category updated successfully")

	return result, nil
}

func (c *categoryUsecase) DeleteCategory(ctx context.Context, id int64) error {
	log := logrus.WithField("category_id", id)

	category, err := c.repo.GetCategoryByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("Category not found")

		return err
	}

	err = c.repo.DeleteCategory(ctx, category.ID)
	if err != nil {
		log.WithError(err).Error("Failed to delete category")

		return err
	}

	log.WithFields(logrus.Fields{"category_id": category.ID}).
		Info("Category deleted successfully")

	return nil
}
