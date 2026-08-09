package model

import "time"

type InventoryTransactionType string

const (
	InventoryTransactionRestock    InventoryTransactionType = "RESTOCK"
	InventoryTransactionSale       InventoryTransactionType = "SALE"
	InventoryTransactionReturn     InventoryTransactionType = "RETURN"
	InventoryTransactionDamage     InventoryTransactionType = "DAMAGE"
	InventoryTransactionAdjustment InventoryTransactionType = "ADJUSTMENT"
	InventoryTransactionRelease    InventoryTransactionType = "RELEASE"
)

type InventoryTransaction struct {
	ID                  string                   `json:"id"`
	ProductVariantID    string                   `json:"product_variant_id"`
	TransactionType     InventoryTransactionType `json:"transaction_type"`
	Quantity            int64                    `json:"quantity"`
	StockBefore         int64                    `json:"stock_before"`
	StockAfter          int64                    `json:"stock_after"`
	ReservedStockBefore int64                    `json:"reserved_stock_before"`
	ReservedStockAfter  int64                    `json:"reserved_stock_after"`
	CreatedBy           string                   `json:"created_by"`
	CreatedAt           time.Time                `json:"created_at"`
}

type InventoryTransactionInput struct {
	ProductVariantID string                   `json:"product_variant_id" validate:"required,uuid"`
	TransactionType  InventoryTransactionType `json:"transaction_type" validate:"required"`
	Quantity         int64                    `json:"quantity" validate:"required,gt=0"`
}

type InventoryTransactionUpdateInput struct {
	ProductVariantID string                   `json:"product_variant_id" validate:"required,uuid"`
	TransactionType  InventoryTransactionType `json:"transaction_type" validate:"required"`
	Quantity         int64                    `json:"quantity" validate:"required,gt=0"`
}

type InventoryOperationInput struct {
	ProductVariantID string `json:"product_variant_id" validate:"required,uuid"`
	Quantity         int64  `json:"quantity" validate:"required,gt=0"`
}

type InventoryAdjustmentInput struct {
	ProductVariantID string `json:"product_variant_id" validate:"required,uuid"`
	Quantity         int64  `json:"quantity" validate:"required"`
}
