package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateRoleInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateRoleInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}
