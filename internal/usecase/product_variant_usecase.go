package usecase

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/repository"
	"gorm.io/gorm"
)

type IProductVariantUsecase interface {
	GetByID(ctx context.Context, id string) (*model.ProductVariant, error)
	GetByProductID(ctx context.Context, productID string) ([]*model.ProductVariant, error)
	Create(ctx context.Context, in *model.ProductVariantInput) (*model.ProductVariant, error)
	Update(ctx context.Context, id string, in *model.ProductVariantUpdateInput) (*model.ProductVariant, error)
	Delete(ctx context.Context, id string) error
}

type productVariantUsecase struct {
	repo        repository.IProductVariantRepository
	productRepo repository.IProductRepository
}

func NewProductVariantUsecase(repo repository.IProductVariantRepository, productRepo repository.IProductRepository) IProductVariantUsecase {
	return &productVariantUsecase{repo: repo, productRepo: productRepo}
}

func (p *productVariantUsecase) GetByID(ctx context.Context, id string) (*model.ProductVariant, error) {
	log := logrus.WithFields(logrus.Fields{"id": id})

	productVariant, err := p.repo.GetByID(ctx, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return productVariant, nil
}

func (p *productVariantUsecase) GetByProductID(ctx context.Context, productID string) ([]*model.ProductVariant, error) {
	log := logrus.WithFields(logrus.Fields{"product_id": productID})

	productVariant, err := p.repo.GetByProductID(ctx, productID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return productVariant, nil
}

func (p *productVariantUsecase) Create(ctx context.Context, in *model.ProductVariantInput) (*model.ProductVariant, error) {
	log := logrus.WithFields(logrus.Fields{
		"input": in,
	})

	_, err := p.productRepo.GetProductByID(ctx, in.ProductID)
	if err != nil {
		log.WithError(err).Error("Product not found")
		return nil, errors.New("product not found")
	}

	_, err = p.repo.GetProductBySKU(ctx, in.SKU)

	if err == nil {
		return nil, errors.New("product variant SKU already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Error("Failed to check product variant SKU")
		return nil, err
	}

	productVariant := &model.ProductVariant{
		ProductID: in.ProductID,
		SKU:       in.SKU,
		Color:     in.Color,
		Size:      in.Size,
		Price:     in.Price,
		Weight:    in.Weight,
		Barcode:   in.Barcode,
		IsActive:  in.IsActive,
	}

	result, err := p.repo.Create(ctx, productVariant)
	if err != nil {
		log.WithError(err).Error("Failed to create product variant")
		return nil, err
	}

	return result, nil
}

func (p *productVariantUsecase) Update(ctx context.Context, id string, in *model.ProductVariantUpdateInput) (*model.ProductVariant, error) {
	log := logrus.WithFields(logrus.Fields{
		"id":    id,
		"input": in,
	})

	productVariant, err := p.repo.GetByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("Product variant not found")
		return nil, err
	}

	if productVariant.SKU != in.SKU {
		existing, err := p.repo.GetProductBySKU(ctx, in.SKU)

		if err == nil && existing.ID != id {
			return nil, errors.New("product variant SKU already exists")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("Failed to check product variant SKU")
			return nil, err
		}
	}

	if productVariant.ProductID != in.ProductID {
		_, err := p.productRepo.GetProductByID(ctx, in.ProductID)
		if err != nil {
			log.WithError(err).Error("Product not found")
			return nil, errors.New("product not found")
		}
	}

	productVariant.ProductID = in.ProductID
	productVariant.SKU = in.SKU
	productVariant.Color = in.Color
	productVariant.Size = in.Size
	productVariant.Price = in.Price
	productVariant.Weight = in.Weight
	productVariant.Barcode = in.Barcode
	productVariant.IsActive = in.IsActive

	result, err := p.repo.Update(ctx, productVariant)
	if err != nil {
		log.WithError(err).Error("Failed to update product variant")
		return nil, err
	}

	return result, nil
}

func (p *productVariantUsecase) Delete(ctx context.Context, id string) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	_, err := p.repo.GetByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("Product variant not found")
		return err
	}

	if err := p.repo.Delete(ctx, id); err != nil {
		log.WithError(err).Error("Failed to delete product variant")
		return err
	}

	return nil
}
