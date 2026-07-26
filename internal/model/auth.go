package model

import "github.com/golang-jwt/jwt/v4"

type ContextAuthKey string

const BearerAuthKey ContextAuthKey = "BearerAuth"

type CustomClaims struct {
	UserID   string `json:"user_id"`
	RoleID   int64  `json:"role_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

type Auth struct {
	UserID   string `json:"user_id"`
	RoleID   int64  `json:"role_id"`
	Role     *Role  `json:"role"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *Auth  `json:"user"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
