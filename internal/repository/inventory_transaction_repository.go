package repository

import (
	"context"

	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"gorm.io/gorm"
)

type IInventoryTransactionRepository interface {
	Create(ctx context.Context, transaction *model.InventoryTransaction) (*model.InventoryTransaction, error)
	GetByProductVariantID(ctx context.Context, productVariantID string) ([]*model.InventoryTransaction, error)
}

type inventoryTransactionRepository struct {
	db *gorm.DB
}

func NewInventoryTransactionRepository(db *gorm.DB) IInventoryTransactionRepository {
	return &inventoryTransactionRepository{db: db}
}

func (r *inventoryTransactionRepository) Create(ctx context.Context, transaction *model.InventoryTransaction) (*model.InventoryTransaction, error) {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *inventoryTransactionRepository) GetByProductVariantID(ctx context.Context, productVariantID string) ([]*model.InventoryTransaction, error) {
	var transactions []*model.InventoryTransaction
	if err := r.db.WithContext(ctx).Where("product_variant_id = ?", productVariantID).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
