package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	ImageURL    string         `json:"image_url"`
	IsActive    bool           `json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CategoryInput struct {
	Name        string `json:"name" validate:"required,max=150"`
	Slug        string `json:"slug" validate:"required,max=150"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`
}

type CategoryUpdateInput struct {
	Name        string `json:"name" validate:"required,max=150"`
	Slug        string `json:"slug" validate:"required,max=150"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`
}

type CategoryFilter struct {
	Search   string
	Sort     string
	Page     int
	Limit    int
	IsActive *bool
}
