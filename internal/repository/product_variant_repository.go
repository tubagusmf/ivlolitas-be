package repository

import (
	"context"

	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"gorm.io/gorm"
)

type IProductVariantRepository interface {
	GetByID(ctx context.Context, id string) (*model.ProductVariant, error)
	GetByProductID(ctx context.Context, productID string) ([]*model.ProductVariant, error)
	GetProductBySKU(ctx context.Context, sku string) (*model.ProductVariant, error)
	Create(ctx context.Context, variant *model.ProductVariant) (*model.ProductVariant, error)
	Update(ctx context.Context, variant *model.ProductVariant) (*model.ProductVariant, error)
	Delete(ctx context.Context, id string) error
}

type ProductVariantRepository struct {
	db *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) IProductVariantRepository {
	return &ProductVariantRepository{db: db}
}

func (r *ProductVariantRepository) GetByID(ctx context.Context, id string) (*model.ProductVariant, error) {
	var variant model.ProductVariant

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&variant).Error; err != nil {
		return nil, err
	}

	return &variant, nil
}

func (r *ProductVariantRepository) GetByProductID(ctx context.Context, productID string) ([]*model.ProductVariant, error) {
	var variants []*model.ProductVariant

	if err := r.db.
		WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at ASC").
		Find(&variants).Error; err != nil {
		return nil, err
	}

	return variants, nil
}

func (r *ProductVariantRepository) GetProductBySKU(ctx context.Context, sku string) (*model.ProductVariant, error) {
	var variant model.ProductVariant

	if err := r.db.
		WithContext(ctx).
		Where("sku = ?", sku).
		First(&variant).Error; err != nil {
		return nil, err
	}

	return &variant, nil
}

func (r *ProductVariantRepository) Create(ctx context.Context, variant *model.ProductVariant) (*model.ProductVariant, error) {
	if err := r.db.WithContext(ctx).Create(variant).Error; err != nil {
		return nil, err
	}

	return variant, nil
}

func (r *ProductVariantRepository) Update(ctx context.Context, variant *model.ProductVariant) (*model.ProductVariant, error) {
	if err := r.db.WithContext(ctx).Model(variant).Updates(variant).Error; err != nil {
		return nil, err
	}

	return variant, nil
}

func (r *ProductVariantRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ProductVariant{}).Error
}
