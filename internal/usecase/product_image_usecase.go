package usecase

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/repository"
	"github.com/tubagusmf/ivlolitas-be/internal/storage"
)

type IProductImageUsecase interface {
	UploadImages(ctx context.Context, productID string, files []*multipart.FileHeader) ([]*model.ProductImage, error)
	GetImagesByProductID(ctx context.Context, productID string) ([]*model.ProductImage, error)
	SetPrimaryImage(ctx context.Context, imageID string) error
	DeleteImage(ctx context.Context, imageID string) error
}

type productImageUsecase struct {
	productRepo      repository.IProductRepository
	productImageRepo repository.IProductImageRepository
	storage          storage.Storage
}

func NewProductImageUsecase(
	productRepo repository.IProductRepository,
	productImageRepo repository.IProductImageRepository,
	storage storage.Storage,
) IProductImageUsecase {
	return &productImageUsecase{
		productRepo:      productRepo,
		productImageRepo: productImageRepo,
		storage:          storage,
	}
}

func (u *productImageUsecase) GetImagesByProductID(ctx context.Context, productID string) ([]*model.ProductImage, error) {
	log := logrus.WithField("product_id", productID)

	_, err := u.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		log.WithError(err).Error("Product not found")
		return nil, err
	}

	images, err := u.productImageRepo.GetByProductID(ctx, productID)
	if err != nil {
		log.WithError(err).Error("Failed to get product images")
		return nil, err
	}

	return images, nil
}

func (u *productImageUsecase) UploadImages(ctx context.Context, productID string, files []*multipart.FileHeader) ([]*model.ProductImage, error) {
	log := logrus.WithField("product_id", productID)

	product, err := u.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		log.WithError(err).Error("Product not found")
		return nil, err
	}

	existingImages, err := u.productImageRepo.GetByProductID(ctx, product.ID)
	if err != nil {
		log.WithError(err).Error("Failed to get product images")
		return nil, err
	}

	var images []*model.ProductImage

	for i, file := range files {

		src, err := file.Open()
		if err != nil {
			log.WithError(err).Error("Failed to open file")
			return nil, err
		}

		result, err := u.storage.Upload(
			ctx,
			src,
			"products",
		)

		src.Close()

		if err != nil {
			log.WithError(err).Error("Failed to upload image")
			return nil, err
		}

		image := &model.ProductImage{
			ID:            uuid.New().String(),
			ProductID:     product.ID,
			ImageURL:      result.URL,
			ImagePublicID: result.PublicID,
			IsPrimary:     len(existingImages) == 0 && i == 0,
			SortOrder:     int64(len(existingImages) + i + 1),
		}

		images = append(images, image)
	}

	if err := u.productImageRepo.CreateImages(ctx, images); err != nil {

		for _, img := range images {
			_ = u.storage.Delete(ctx, img.ImagePublicID)
		}

		log.WithError(err).Error("Failed to save product images")

		return nil, err
	}

	return images, nil
}

func (u *productImageUsecase) SetPrimaryImage(ctx context.Context, imageID string) error {
	image, err := u.productImageRepo.GetByID(ctx, imageID)
	if err != nil {
		return err
	}

	images, err := u.productImageRepo.GetByProductID(ctx, image.ProductID)
	if err != nil {
		return err
	}

	for _, img := range images {

		if img.ID == image.ID {
			continue
		}

		if img.IsPrimary {

			img.IsPrimary = false

			if err := u.productImageRepo.Update(ctx, img); err != nil {
				return err
			}
		}
	}

	image.IsPrimary = true

	return u.productImageRepo.Update(ctx, image)
}

func (u *productImageUsecase) DeleteImage(ctx context.Context, imageID string) error {
	log := logrus.WithField("image_id", imageID)

	_, err := u.productImageRepo.GetByID(ctx, imageID)
	if err != nil {
		log.WithError(err).Error("Image not found")
		return err
	}

	image, err := u.productImageRepo.GetByID(ctx, imageID)
	if err != nil {
		return err
	}

	if err := u.storage.Delete(ctx, image.ImagePublicID); err != nil {
		return err
	}

	return u.productImageRepo.DeleteByID(ctx, imageID)
}
