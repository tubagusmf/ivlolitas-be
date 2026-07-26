package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ivlolitas-be/internal/model"

	jwt "github.com/tubagusmf/ivlolitas-be/internal/jwt"
)

type AuthMiddleware struct {
	jwt *jwt.JWT
}

func NewAuthMiddleware(jwt *jwt.JWT) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}

func (m *AuthMiddleware) JWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")

		if auth == "" {
			return c.JSON(http.StatusUnauthorized, "missing token")
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, "invalid authorization header")
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		if tokenString == "" {
			return c.JSON(
				http.StatusUnauthorized,
				"token required",
			)
		}

		claims, err := m.jwt.ParseAccessToken(tokenString)
		if err != nil {
			c.Logger().Error(err)

			return c.JSON(
				http.StatusUnauthorized,
				"invalid token",
			)
		}

		c.Set(string(model.BearerAuthKey), &model.Auth{
			UserID:   claims.UserID,
			RoleID:   claims.RoleID,
			Email:    claims.Email,
			FullName: claims.FullName,
		})

		return next(c)
	}
}
