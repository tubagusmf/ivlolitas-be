package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	ID            string         `json:"id"`
	ProductID     string         `json:"product_id"`
	Product       *Product       `json:"product"`
	SKU           string         `json:"sku"`
	Color         *string        `json:"color"`
	Size          *string        `json:"size"`
	Price         float64        `json:"price"`
	Weight        float64        `json:"weight"`
	Barcode       *string        `json:"barcode"`
	IsActive      bool           `json:"is_active"`
	ImageURL      string         `json:"image_url"`
	ImagePublicID string         `json:"image_public_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductVariantInput struct {
	SKU      string  `json:"sku" form:"sku" validate:"required,max=50"`
	Color    *string `json:"color" form:"color" validate:"omitempty,max=100"`
	Size     *string `json:"size" form:"size" validate:"omitempty,max=50"`
	Price    float64 `json:"price" form:"price" validate:"required,gt=0"`
	Weight   float64 `json:"weight" form:"weight" validate:"gte=0"`
	Barcode  *string `json:"barcode" form:"barcode" validate:"omitempty,max=100"`
	IsActive bool    `json:"is_active" form:"is_active"`
}

type ProductVariantUpdateInput struct {
	SKU      string  `json:"sku" form:"sku" validate:"required,max=50"`
	Color    *string `json:"color" form:"color" validate:"omitempty,max=100"`
	Size     *string `json:"size" form:"size" validate:"omitempty,max=50"`
	Price    float64 `json:"price" form:"price" validate:"required,gt=0"`
	Weight   float64 `json:"weight" form:"weight" validate:"gte=0"`
	Barcode  *string `json:"barcode" form:"barcode" validate:"omitempty,max=100"`
	IsActive bool    `json:"is_active" form:"is_active"`
}
