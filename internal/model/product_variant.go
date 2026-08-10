package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	ID        string         `json:"id"`
	ProductID string         `json:"product_id"`
	Product   *Product       `json:"product"`
	SKU       string         `json:"sku"`
	Color     *string        `json:"color"`
	Size      *string        `json:"size"`
	Price     float64        `json:"price"`
	Weight    float64        `json:"weight"`
	Barcode   *string        `json:"barcode"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductVariantInput struct {
	SKU      string  `json:"sku" validate:"required,max=50"`
	Color    *string `json:"color" validate:"max=100"`
	Size     *string `json:"size" validate:"max=50"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	Weight   float64 `json:"weight" validate:"gte=0"`
	Barcode  *string `json:"barcode" validate:"max=100"`
	IsActive bool    `json:"is_active"`
}

type ProductVariantUpdateInput struct {
	SKU      string  `json:"sku" validate:"required,max=50"`
	Color    *string `json:"color" validate:"max=100"`
	Size     *string `json:"size" validate:"max=50"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	Weight   float64 `json:"weight" validate:"gte=0"`
	Barcode  *string `json:"barcode" validate:"max=100"`
	IsActive bool    `json:"is_active"`
}
