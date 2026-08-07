package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/repository"
	"gorm.io/gorm"
)

type IProductUsecase interface {
	GetProducts(ctx context.Context, filter *model.ProductFilter) ([]*model.Product, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
	CreateProduct(ctx context.Context, in *model.ProductInput) (*model.Product, error)
	UpdateProduct(ctx context.Context, id string, in *model.ProductUpdateInput) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type productUsecase struct {
	repo         repository.IProductRepository
	categoryRepo repository.ICategoryRepository
}

func NewProductUsecase(repo repository.IProductRepository, categoryRepo repository.ICategoryRepository) IProductUsecase {
	return &productUsecase{repo: repo, categoryRepo: categoryRepo}
}

func (p *productUsecase) GetProducts(ctx context.Context, filter *model.ProductFilter) ([]*model.Product, error) {
	log := logrus.WithFields(logrus.Fields{"filter": filter})

	products, err := p.repo.GetProducts(ctx, filter)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return products, nil
}

func (p *productUsecase) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	log := logrus.WithFields(logrus.Fields{"id": id})

	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return product, nil
}

func (p *productUsecase) CreateProduct(ctx context.Context, in *model.ProductInput) (*model.Product, error) {
	log := logrus.WithFields(logrus.Fields{"input": in})

	_, err := p.repo.GetProductBySlug(ctx, in.Slug)
	if err == nil {
		log.Warn("Product slug already exists")

		return nil, errors.New("product slug already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Error("Failed to check product slug")

		return nil, err
	}

	_, err = p.categoryRepo.GetCategoryByID(ctx, in.CategoryID)
	if err != nil {
		log.WithError(err).Error("Category not found")

		return nil, err
	}

	claims, ok := ctx.Value(model.BearerAuthKey).(*model.CustomClaims)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	product := &model.Product{
		ID:               uuid.New().String(),
		CategoryID:       in.CategoryID,
		Name:             in.Name,
		Slug:             in.Slug,
		ShortDescription: &in.ShortDescription,
		Description:      &in.Description,
		IsActive:         in.IsActive,
		CreatedBy:        claims.UserID,
		UpdatedBy:        claims.UserID,
	}

	result, err := p.repo.CreateProduct(ctx, product)
	if err != nil {
		log.WithError(err).Error("Failed to create product")

		return nil, err
	}

	return result, nil
}

func (p *productUsecase) UpdateProduct(ctx context.Context, id string, in *model.ProductUpdateInput) (*model.Product, error) {
	log := logrus.WithFields(logrus.Fields{"product_id": id, "input": in})

	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("Product not found")

		return nil, err
	}

	if product.Slug != in.Slug {
		existing, err := p.repo.GetProductBySlug(ctx, in.Slug)

		if err == nil && existing.ID != product.ID {
			log.Warn("Product slug already exists")

			return nil, errors.New("product slug already exists")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("Failed to check product slug")

			return nil, err
		}
	}

	if product.CategoryID != in.CategoryID {
		_, err := p.categoryRepo.GetCategoryByID(ctx, in.CategoryID)
		if err != nil {
			log.WithError(err).Error("Category not found")

			return nil, err
		}
	}

	claims, ok := ctx.Value(model.BearerAuthKey).(*model.CustomClaims)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	product.CategoryID = in.CategoryID
	product.Name = in.Name
	product.Slug = in.Slug
	product.ShortDescription = &in.ShortDescription
	product.Description = &in.Description
	product.IsActive = in.IsActive
	product.UpdatedBy = claims.UserID

	result, err := p.repo.UpdateProduct(ctx, product)
	if err != nil {
		log.WithError(err).Error("Failed to update product")

		return nil, err
	}

	return result, nil
}

func (p *productUsecase) DeleteProduct(ctx context.Context, id string) error {
	log := logrus.WithFields(logrus.Fields{"product_id": id})

	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("Product not found")

		return err
	}

	err = p.repo.DeleteProduct(ctx, product.ID)
	if err != nil {
		log.WithError(err).Error("Failed to delete product")

		return err
	}

	log.WithFields(logrus.Fields{"product_id": product.ID}).Info("Product deleted successfully")

	return nil
}
