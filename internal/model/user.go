package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          string         `json:"id"`
	RoleID      int64          `json:"role_id"`
	Role        *Role          `json:"role"`
	FullName    string         `json:"full_name"`
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	Address     string         `json:"address"`
	PhoneNumber string         `json:"phone_number"`
	IsActive    bool           `json:"is_active"`
	LastLogin   *time.Time     `json:"last_login"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateUserInput struct {
	FullName    string `json:"full_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"omitempty,min=8,max=100"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	RoleID      int64  `json:"role_id" validate:"required"`
}

type UpdateUserInput struct {
	FullName    string `json:"full_name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"omitempty,min=8,max=100"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	RoleID      int64  `json:"role_id" validate:"required"`
}
