package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/usecase"
)

type AuthHandler struct {
	authUsecase usecase.IAuthUsecase
}

func NewAuthHandler(e *echo.Echo, uc usecase.IAuthUsecase, middleware *AuthMiddleware) {
	handler := &AuthHandler{
		authUsecase: uc,
	}

	auth := e.Group("/v1/auth")

	auth.POST("/login", handler.Login)
	auth.POST("/refresh", handler.Refresh)
	auth.POST("/logout", handler.Logout, middleware.JWT)
	auth.POST("/logout-all", handler.LogoutAll, middleware.JWT)
}

// Login godoc
//
//	@Summary		User Login
//	@Description	Authenticate user using email and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.LoginInput	true	"Login Request"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	result, err := h.authUsecase.Login(
		c.Request().Context(),
		&req,
	)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "login success",
		"data":    result,
	})
}

// Refresh godoc
//
//	@Summary		Refresh Access Token
//	@Description	Get new access token using refresh token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.RefreshTokenRequest	true	"Refresh Token"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/auth/refresh [post]
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req model.RefreshTokenRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	result, err := h.authUsecase.RefreshToken(
		c.Request().Context(),
		req.RefreshToken,
	)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "refresh token success",
		"data":    result,
	})
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Logout current user by invalidating refresh token
//	@Tags			Authentication
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.LogoutRequest	true	"Logout Request"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	var req model.LogoutRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.authUsecase.Logout(
		c.Request().Context(),
		req.RefreshToken,
	); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "logout success",
	})
}

// LogoutAll godoc
//
//	@Summary		Logout All Devices
//	@Description	Invalidate all refresh tokens of current user
//	@Tags			Authentication
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c echo.Context) error {
	auth := c.Get(string(model.BearerAuthKey)).(*model.Auth)

	if err := h.authUsecase.LogoutAll(
		c.Request().Context(),
		auth.UserID,
	); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "logout all success",
	})
}
