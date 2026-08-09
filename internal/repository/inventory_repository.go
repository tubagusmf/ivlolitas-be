package repository

import (
	"context"

	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IInventoryRepository interface {
	Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)
	GetByID(ctx context.Context, id string) (*model.Inventory, error)
	GetByProductVariantID(ctx context.Context, productVariantID string) (*model.Inventory, error)
	GetByProductVariantIDForUpdate(ctx context.Context, productVariantID string) (*model.Inventory, error)
	Update(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)
	Delete(ctx context.Context, id string) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) IInventoryRepository {
	return &inventoryRepository{
		db: db,
	}
}

func (r *inventoryRepository) Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	if err := r.db.WithContext(ctx).Create(inventory).Error; err != nil {
		return nil, err
	}

	return inventory, nil
}

func (r *inventoryRepository) GetByID(ctx context.Context, id string) (*model.Inventory, error) {
	var inventory model.Inventory

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&inventory).Error; err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (r *inventoryRepository) GetByProductVariantID(ctx context.Context, productVariantID string) (*model.Inventory, error) {
	var inventory model.Inventory

	if err := r.db.WithContext(ctx).Where("product_variant_id = ?", productVariantID).First(&inventory).Error; err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (r *inventoryRepository) GetByProductVariantIDForUpdate(ctx context.Context, productVariantID string) (*model.Inventory, error) {
	var inventory model.Inventory

	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("product_variant_id = ?", productVariantID).
		First(&inventory).
		Error

	if err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (r *inventoryRepository) Update(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	if err := r.db.WithContext(ctx).Model(inventory).Updates(inventory).Error; err != nil {
		return nil, err
	}

	return inventory, nil
}

func (r *inventoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Inventory{}).Error
}
