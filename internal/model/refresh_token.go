package model

import "time"

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`

	User *User `gorm:"foreignKey:UserID"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
