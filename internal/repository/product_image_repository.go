package repository

import (
	"context"

	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"gorm.io/gorm"
)

type IProductImageRepository interface {
	CreateImages(ctx context.Context, images []*model.ProductImage) error
	GetByProductID(ctx context.Context, productID string) ([]*model.ProductImage, error)
	GetByID(ctx context.Context, id string) (*model.ProductImage, error)
	Update(ctx context.Context, image *model.ProductImage) error
	DeleteByID(ctx context.Context, id string) error
	DeleteByProductID(ctx context.Context, productID string) error
}

type ProductImageRepository struct {
	db *gorm.DB
}

func NewProductImageRepository(db *gorm.DB) IProductImageRepository {
	return &ProductImageRepository{db: db}
}

func (r *ProductImageRepository) CreateImages(ctx context.Context, images []*model.ProductImage) error {
	return r.db.WithContext(ctx).Create(images).Error
}

func (r *ProductImageRepository) GetByProductID(ctx context.Context, productID string) ([]*model.ProductImage, error) {
	var images []*model.ProductImage

	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("sort_order ASC").Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ProductImageRepository) GetByID(ctx context.Context, id string) (*model.ProductImage, error) {
	var image model.ProductImage

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&image).Error; err != nil {
		return nil, err
	}

	return &image, nil
}

func (r *ProductImageRepository) Update(ctx context.Context, image *model.ProductImage) error {
	return r.db.WithContext(ctx).Model(&model.ProductImage{}).Where("id = ?", image.ID).Updates(image).Error
}

func (r *ProductImageRepository) DeleteByID(ctx context.Context, id string) error {
	image := &model.ProductImage{ID: id}
	return r.db.WithContext(ctx).Delete(image).Error
}

func (r *ProductImageRepository) DeleteByProductID(ctx context.Context, productID string) error {
	return r.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&model.ProductImage{}).Error
}
