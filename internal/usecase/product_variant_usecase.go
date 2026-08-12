package usecase

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/repository"
	"github.com/tubagusmf/ivlolitas-be/internal/storage"
	"gorm.io/gorm"
)

type IProductVariantUsecase interface {
	GetByID(ctx context.Context, id string) (*model.ProductVariant, error)
	GetByProductID(ctx context.Context, productID string) ([]*model.ProductVariant, error)
	Create(ctx context.Context, productID string, in *model.ProductVariantInput, file *multipart.FileHeader) (*model.ProductVariant, error)
	Update(ctx context.Context, id string, in *model.ProductVariantUpdateInput, file *multipart.FileHeader) (*model.ProductVariant, error)
	Delete(ctx context.Context, id string) error
}

type productVariantUsecase struct {
	repo        repository.IProductVariantRepository
	productRepo repository.IProductRepository
	storage     storage.Storage
}

func NewProductVariantUsecase(repo repository.IProductVariantRepository, productRepo repository.IProductRepository, storage storage.Storage) IProductVariantUsecase {
	return &productVariantUsecase{repo: repo, productRepo: productRepo, storage: storage}
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

func (p *productVariantUsecase) Create(ctx context.Context, productID string, in *model.ProductVariantInput, file *multipart.FileHeader) (*model.ProductVariant, error) {
	log := logrus.WithFields(logrus.Fields{
		"input": in,
	})

	_, err := p.productRepo.GetProductByID(ctx, productID)
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

	src, err := file.Open()
	if err != nil {
		log.WithError(err).Error("Failed to open file")
		return nil, err
	}
	defer src.Close()

	uploadResult, err := p.storage.Upload(
		ctx,
		src,
		"products/variants",
	)
	if err != nil {
		log.WithError(err).Error("Failed to upload product variant image")
		return nil, err
	}

	productVariant := &model.ProductVariant{
		ID:            uuid.New().String(),
		ProductID:     productID,
		SKU:           in.SKU,
		Color:         in.Color,
		Size:          in.Size,
		Price:         in.Price,
		Weight:        in.Weight,
		Barcode:       in.Barcode,
		IsActive:      in.IsActive,
		ImageURL:      uploadResult.URL,
		ImagePublicID: uploadResult.PublicID,
	}

	result, err := p.repo.Create(ctx, productVariant)
	if err != nil {
		_ = p.storage.Delete(ctx, uploadResult.PublicID)

		log.WithError(err).Error("Failed to create product variant")
		return nil, err
	}

	return result, nil
}

func (p *productVariantUsecase) Update(ctx context.Context, id string, in *model.ProductVariantUpdateInput, file *multipart.FileHeader) (*model.ProductVariant, error) {
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

	productVariant.SKU = in.SKU
	productVariant.Color = in.Color
	productVariant.Size = in.Size
	productVariant.Price = in.Price
	productVariant.Weight = in.Weight
	productVariant.Barcode = in.Barcode
	productVariant.IsActive = in.IsActive

	var oldPublicID string
	var newPublicID string

	if file != nil {
		oldPublicID = productVariant.ImagePublicID

		src, err := file.Open()
		if err != nil {
			log.WithError(err).Error("Failed to open file")
			return nil, err
		}
		defer src.Close()

		uploadResult, err := p.storage.Upload(
			ctx,
			src,
			"products/variants",
		)
		if err != nil {
			log.WithError(err).Error("Failed to upload product variant image")
			return nil, err
		}

		productVariant.ImageURL = uploadResult.URL
		productVariant.ImagePublicID = uploadResult.PublicID

		newPublicID = uploadResult.PublicID
	}

	result, err := p.repo.Update(ctx, productVariant)
	if err != nil {
		if newPublicID != "" {
			_ = p.storage.Delete(ctx, newPublicID)
		}

		log.WithError(err).Error("Failed to update product variant")
		return nil, err
	}

	if oldPublicID != "" {
		if err := p.storage.Delete(ctx, oldPublicID); err != nil {
			log.WithError(err).Warn("Failed to delete old product variant image")
		}
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
