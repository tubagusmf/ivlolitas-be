package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v4"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
)

type JWT struct {
	SecretKey string
}

func New(secret string) *JWT {
	return &JWT{
		SecretKey: secret,
	}
}

// Generate Access Token
func (j *JWT) GenerateAccessToken(user *model.User) (string, error) {
	claims := model.CustomClaims{
		UserID:   user.ID,
		RoleID:   user.RoleID,
		Email:    user.Email,
		FullName: user.FullName,
		RegisteredClaims: jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwtgo.NewNumericDate(time.Now()),
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.SecretKey))
}

// Generate Refresh Token
func (j *JWT) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// Parse Access Token
func (j *JWT) ParseAccessToken(tokenString string) (*model.CustomClaims, error) {
	token, err := jwtgo.ParseWithClaims(
		tokenString,
		&model.CustomClaims{},
		func(token *jwtgo.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.CustomClaims)
	if !ok || !token.Valid {
		return nil, jwtgo.ErrSignatureInvalid
	}

	return claims, nil
}
