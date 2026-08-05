package repository

import (
	"context"

	"github.com/tubagusmf/ivlolitas-be/internal/helper"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"gorm.io/gorm"
)

type IProductRepository interface {
	GetProducts(ctx context.Context, filter *model.ProductFilter) ([]*model.Product, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*model.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*model.Product, error)
	CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProducts(ctx context.Context, filter *model.ProductFilter) ([]*model.Product, error) {
	var products []*model.Product

	db := r.db.WithContext(ctx)

	if filter.Search != "" {
		db = db.Where("name ILIKE ? OR sku ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	if filter.IsActive != nil {
		db = db.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Sort != "" {
		db = db.Order(helper.BuildProductSort(filter.Sort))
	} else {
		db = db.Order("created_at DESC")
	}

	helper.NormalizePagination(filter)

	err := db.
		Offset((filter.Page - 1) * filter.Limit).
		Limit(filter.Limit).
		Preload("ProductImages").
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product

	db := r.db.WithContext(ctx)

	if err := db.Where("id = ?", id).First(&product).Preload("ProductImages").Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetProductBySlug(ctx context.Context, slug string) (*model.Product, error) {
	var product model.Product

	db := r.db.WithContext(ctx)

	if err := db.Where("slug = ?", slug).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetProductBySKU(ctx context.Context, sku string) (*model.Product, error) {
	var product model.Product

	db := r.db.WithContext(ctx)

	if err := db.Where("sku = ?", sku).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	if err := r.db.WithContext(ctx).Model(product).Updates(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	product := &model.Product{ID: id}
	return r.db.WithContext(ctx).Delete(product).Error
}
