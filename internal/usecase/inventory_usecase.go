package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInventoryNotFound   = errors.New("inventory not found")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrInsufficientReserve = errors.New("insufficient reserved stock")
	ErrInvalidQuantity     = errors.New("quantity must be greater than zero")
)

type IInventoryUsecase interface {
	GetInventory(ctx context.Context, productVariantID string) (*model.Inventory, error)
	GetTransactions(ctx context.Context, productVariantID string) ([]*model.InventoryTransaction, error)
	Restock(ctx context.Context, productVariantID string, quantity int64, createdBy string) error
	Sale(ctx context.Context, productVariantID string, quantity int64, createdBy string) error
	Return(ctx context.Context, productVariantID string, quantity int64, createdBy string) error
	Damage(ctx context.Context, productVariantID string, quantity int64, createdBy string) error
	Adjustment(ctx context.Context, productVariantID string, quantity int64, createdBy string) error
	Release(ctx context.Context, productVariantID string, quantity int64, createdBy string) error
}

type inventoryUsecase struct {
	inventoryRepository            repository.IInventoryRepository
	inventoryTransactionRepository repository.IInventoryTransactionRepository
	db                             *gorm.DB
}

func NewInventoryUsecase(inventoryRepository repository.IInventoryRepository, inventoryTransactionRepository repository.IInventoryTransactionRepository, db *gorm.DB) IInventoryUsecase {
	return &inventoryUsecase{
		inventoryRepository:            inventoryRepository,
		inventoryTransactionRepository: inventoryTransactionRepository,
		db:                             db,
	}
}

func (u *inventoryUsecase) GetInventory(ctx context.Context, productVariantID string) (*model.Inventory, error) {
	log := logrus.WithFields(logrus.Fields{"product_variant_id": productVariantID})

	inventory, err := u.inventoryRepository.GetByProductVariantID(ctx, productVariantID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return inventory, nil
}

func (u *inventoryUsecase) GetTransactions(ctx context.Context, productVariantID string) ([]*model.InventoryTransaction, error) {
	log := logrus.WithFields(logrus.Fields{"product_variant_id": productVariantID})

	transactions, err := u.inventoryTransactionRepository.GetByProductVariantID(ctx, productVariantID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return transactions, nil
}

func (u *inventoryUsecase) changeStock(ctx context.Context, productVariantID string, transactionType model.InventoryTransactionType, quantity int64, createdBy string, increase bool) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inventory model.Inventory

		err := tx.
			Clauses(clause.Locking{
				Strength: "UPDATE",
			}).
			Where("product_variant_id = ?", productVariantID).
			First(&inventory).
			Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInventoryNotFound
			}

			return err
		}

		stockBefore := inventory.Stock
		reservedStockBefore := inventory.ReservedStock

		if increase {
			inventory.Stock += quantity
		} else {
			if inventory.Stock < quantity {
				return ErrInsufficientStock
			}

			inventory.Stock -= quantity
		}

		stockAfter := inventory.Stock
		reservedStockAfter := inventory.ReservedStock

		err = tx.
			Model(&model.Inventory{}).
			Where("id = ?", inventory.ID).
			Updates(map[string]interface{}{
				"stock":      stockAfter,
				"updated_at": time.Now(),
			}).
			Error

		if err != nil {
			return err
		}

		transaction := &model.InventoryTransaction{
			ID:                  uuid.New().String(),
			ProductVariantID:    productVariantID,
			TransactionType:     transactionType,
			Quantity:            quantity,
			StockBefore:         stockBefore,
			StockAfter:          stockAfter,
			ReservedStockBefore: reservedStockBefore,
			ReservedStockAfter:  reservedStockAfter,
			CreatedBy:           createdBy,
		}

		return tx.Create(transaction).Error
	})
}

func (u *inventoryUsecase) releaseReservedStock(ctx context.Context, productVariantID string, quantity int64, createdBy string) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inventory model.Inventory

		err := tx.
			Clauses(clause.Locking{
				Strength: "UPDATE",
			}).
			Where("product_variant_id = ?", productVariantID).
			First(&inventory).
			Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInventoryNotFound
			}

			return err
		}

		if inventory.ReservedStock < quantity {
			return ErrInsufficientReserve
		}

		stockBefore := inventory.Stock
		reservedStockBefore := inventory.ReservedStock

		inventory.ReservedStock -= quantity

		stockAfter := inventory.Stock
		reservedStockAfter := inventory.ReservedStock

		err = tx.
			Model(&model.Inventory{}).
			Where("id = ?", inventory.ID).
			Updates(map[string]interface{}{
				"reserved_stock": reservedStockAfter,
				"updated_at":     time.Now(),
			}).
			Error

		if err != nil {
			return err
		}

		transaction := &model.InventoryTransaction{
			ID:                  uuid.New().String(),
			ProductVariantID:    productVariantID,
			TransactionType:     model.InventoryTransactionRelease,
			Quantity:            quantity,
			StockBefore:         stockBefore,
			StockAfter:          stockAfter,
			ReservedStockBefore: reservedStockBefore,
			ReservedStockAfter:  reservedStockAfter,
			CreatedBy:           createdBy,
		}

		return tx.Create(transaction).Error
	})
}

func (u *inventoryUsecase) adjustStock(ctx context.Context, productVariantID string, quantity int64, createdBy string) error {
	if quantity == 0 {
		return errors.New("adjustment quantity cannot be zero")
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inventory model.Inventory

		err := tx.
			Clauses(clause.Locking{
				Strength: "UPDATE",
			}).
			Where("product_variant_id = ?", productVariantID).
			First(&inventory).
			Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInventoryNotFound
			}

			return err
		}

		stockBefore := inventory.Stock
		reservedStockBefore := inventory.ReservedStock

		stockAfter := inventory.Stock + quantity

		if stockAfter < 0 {
			return ErrInsufficientStock
		}

		inventory.Stock = stockAfter

		reservedStockAfter := inventory.ReservedStock

		err = tx.
			Model(&model.Inventory{}).
			Where("id = ?", inventory.ID).
			Updates(map[string]interface{}{
				"stock":      stockAfter,
				"updated_at": time.Now(),
			}).
			Error

		if err != nil {
			return err
		}

		transaction := &model.InventoryTransaction{
			ID:                  uuid.New().String(),
			ProductVariantID:    productVariantID,
			TransactionType:     model.InventoryTransactionAdjustment,
			Quantity:            quantity,
			StockBefore:         stockBefore,
			StockAfter:          stockAfter,
			ReservedStockBefore: reservedStockBefore,
			ReservedStockAfter:  reservedStockAfter,
			CreatedBy:           createdBy,
		}

		return tx.Create(transaction).Error
	})
}

func (u *inventoryUsecase) Restock(ctx context.Context, productVariantID string, quantity int64, createdBy string) error {
	return u.changeStock(ctx, productVariantID, model.InventoryTransactionRestock, quantity, createdBy, true)
}

func (u *inventoryUsecase) Sale(ctx context.Context, productVariantID string, quantity int64, createdBy string) error {
	return u.changeStock(ctx, productVariantID, model.InventoryTransactionSale, quantity, createdBy, false)
}

func (u *inventoryUsecase) Return(ctx context.Context, productVariantID string, quantity int64, createdBy string) error {
	return u.changeStock(ctx, productVariantID, model.InventoryTransactionReturn, quantity, createdBy, true)
}

func (u *inventoryUsecase) Damage(ctx context.Context, productVariantID string, quantity int64, createdBy string) error {
	return u.changeStock(ctx, productVariantID, model.InventoryTransactionDamage, quantity, createdBy, false)
}

func (u *inventoryUsecase) Adjustment(ctx context.Context, productVariantID string, quantity int64, createdBy string) error {
	return u.adjustStock(ctx, productVariantID, quantity, createdBy)
}

func (u *inventoryUsecase) Release(ctx context.Context, productVariantID string, quantity int64, createdBy string) error {
	return u.releaseReservedStock(ctx, productVariantID, quantity, createdBy)
}
