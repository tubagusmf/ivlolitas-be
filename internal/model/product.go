package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID               string         `json:"id"`
	CategoryID       int64          `json:"category_id"`
	SKU              string         `json:"sku"`
	Name             string         `json:"name"`
	Slug             string         `json:"slug"`
	ShortDescription *string        `json:"short_desc"`
	Description      *string        `json:"description"`
	Price            float64        `json:"price"`
	Weight           float64        `json:"weight"`
	Stock            int64          `json:"stock"`
	IsActive         bool           `json:"is_active"`
	ProductImages    []ProductImage `json:"images,omitempty"`
	CreatedBy        string         `json:"created_by"`
	UpdatedBy        string         `json:"updated_by"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductImage struct {
	ID            string         `json:"id"`
	ProductID     string         `json:"product_id"`
	ImageURL      string         `json:"image_url"`
	ImagePublicID string         `json:"image_public_id"`
	IsPrimary     bool           `json:"is_primary"`
	SortOrder     int64          `json:"sort_order"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductFilter struct {
	Search   string
	Sort     string
	Page     int
	Limit    int
	IsActive *bool
}

type ProductInput struct {
	CategoryID       int64   `json:"category_id" validate:"required,gt=0"`
	SKU              string  `json:"sku" validate:"required"`
	Name             string  `json:"name" validate:"required,max=150"`
	Slug             string  `json:"slug" validate:"required,max=150"`
	ShortDescription string  `json:"short_desc"`
	Description      string  `json:"description"`
	Price            float64 `json:"price" validate:"required,gt=0"`
	Weight           float64 `json:"weight" validate:"gte=0"`
	Stock            int64   `json:"stock" validate:"required,gte=0"`
	IsActive         bool    `json:"is_active"`
}

type ProductUpdateInput struct {
	CategoryID       int64   `json:"category_id"`
	SKU              string  `json:"sku"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	ShortDescription string  `json:"short_desc"`
	Description      string  `json:"description"`
	Price            float64 `json:"price"`
	Weight           float64 `json:"weight"`
	Stock            int64   `json:"stock"`
	IsActive         bool    `json:"is_active"`
	UpdatedBy        string  `json:"updated_by"`
}
