package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/usecase"
)

type userHandler struct {
	userUsecase usecase.IUserUsecase
}

func NewUserHandler(e *echo.Echo, uc usecase.IUserUsecase, auth *AuthMiddleware) {
	handler := &userHandler{
		userUsecase: uc,
	}

	group := e.Group("/v1/users")

	group.POST("", handler.CreateUser)
	group.GET("", handler.GetUsers, auth.JWT)
	group.GET("/:id", handler.GetUserByID, auth.JWT)
	group.PUT("/:id", handler.UpdateUser, auth.JWT)
	group.DELETE("/:id", handler.DeleteUser, auth.JWT)
}

// CreateUser godoc
//
//	@Summary		Create User
//	@Description	Register new user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.CreateUserInput	true	"User"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Router			/users [post]
func (h *userHandler) CreateUser(c echo.Context) error {
	var body model.CreateUserInput

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	// TODO: validation
	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.userUsecase.CreateUser(
		c.Request().Context(),
		&body,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "user created successfully",
		"data":    user,
	})
}

// GetUsers godoc
//
//	@Summary		Get Users
//	@Description	Get list of users
//	@Tags			Users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			email	query		string	false	"Filter by email"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/users [get]
func (h *userHandler) GetUsers(c echo.Context) error {
	email := c.QueryParam("email")

	users, err := h.userUsecase.GetUsers(
		c.Request().Context(),
		email,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "users fetched successfully",
		"data":    users,
	})
}

// GetUserByID godoc
//
//	@Summary		Get User By ID
//	@Description	Get user detail
//	@Tags			Users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"User UUID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/users/{id} [get]
func (h *userHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid user id")
	}

	user, err := h.userUsecase.GetUserByID(
		c.Request().Context(),
		id,
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "user fetched successfully",
		"data":    user,
	})
}

// UpdateUser godoc
//
//	@Summary		Update User
//	@Description	Update user by ID
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"User UUID"
//	@Param			request	body		model.UpdateUserInput	true	"User"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Router			/users/{id} [put]
func (h *userHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid user id")
	}

	var body model.UpdateUserInput

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	user, err := h.userUsecase.UpdateUser(
		c.Request().Context(),
		id,
		&body,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "user updated successfully",
		"data":    user,
	})
}

// DeleteUser godoc
//
//	@Summary		Delete User
//	@Description	Delete user by ID
//	@Tags			Users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"User UUID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/users/{id} [delete]
func (h *userHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid user id")
	}

	if err := h.userUsecase.DeleteUser(
		c.Request().Context(),
		id,
	); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "user deleted successfully",
	})
}
