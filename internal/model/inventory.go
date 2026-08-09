package model

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID               string         `json:"id"`
	ProductVariantID string         `json:"product_variant_id"`
	Stock            int64          `json:"stock"`
	ReservedStock    int64          `json:"reserved_stock"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type InventoryInput struct {
	ProductVariantID string `json:"product_variant_id" validate:"required,uuid"`
}

type InventoryUpdateInput struct {
	Stock         int64 `json:"stock" validate:"gte=0"`
	ReservedStock int64 `json:"reserved_stock" validate:"gte=0"`
}
