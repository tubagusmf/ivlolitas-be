package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/usecase"
)

type roleHandler struct {
	roleUsecase usecase.IRoleUsecase
}

func NewroleHandler(e *echo.Echo, uc usecase.IRoleUsecase, auth *AuthMiddleware) {
	handler := &roleHandler{
		roleUsecase: uc,
	}

	group := e.Group("/v1/roles")

	group.POST("", handler.CreateRole, auth.JWT)
	group.GET("", handler.GetRoles, auth.JWT)
	group.GET("/:id", handler.GetRoleByID, auth.JWT)
	group.PUT("/:id", handler.UpdateRole, auth.JWT)
	group.DELETE("/:id", handler.DeleteRole, auth.JWT)
}

// CreateRole godoc
//
//	@Summary		Create Role
//	@Description	Create new role
//	@Tags			Roles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.CreateRoleInput	true	"Role"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Router			/roles [post]
func (h *roleHandler) CreateRole(c echo.Context) error {
	var body model.CreateRoleInput

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	// TODO: validation
	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	role, err := h.roleUsecase.CreateRole(
		c.Request().Context(),
		&body,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "data created successfully",
		"data":    role,
	})
}

// GetRoles godoc
//
//	@Summary		Get Roles
//	@Description	Get all roles
//	@Tags			Roles
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/roles [get]
func (h *roleHandler) GetRoles(c echo.Context) error {
	roles, err := h.roleUsecase.GetRoles(
		c.Request().Context(),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "data fetched successfully",
		"data":    roles,
	})
}

// GetRoleByID godoc
//
//	@Summary		Get Role By ID
//	@Description	Get role detail
//	@Tags			Roles
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Role ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/roles/{id} [get]
func (h *roleHandler) GetRoleByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid role id")
	}

	role, err := h.roleUsecase.GetRoleByID(
		c.Request().Context(),
		id,
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "data fetched successfully",
		"data":    role,
	})
}

// UpdateRole godoc
//
//	@Summary		Update Role
//	@Description	Update role by ID
//	@Tags			Roles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Role ID"
//	@Param			request	body		model.UpdateRoleInput	true	"Role"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Router			/roles/{id} [put]
func (h *roleHandler) UpdateRole(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid role id")
	}

	var body model.UpdateRoleInput

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	role, err := h.roleUsecase.UpdateRole(
		c.Request().Context(),
		id,
		&body,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "data updated successfully",
		"data":    role,
	})
}

// DeleteRole godoc
//
//	@Summary		Delete Role
//	@Description	Delete role by ID
//	@Tags			Roles
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Role ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/roles/{id} [delete]
func (h *roleHandler) DeleteRole(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid role id")
	}

	if err := h.roleUsecase.DeleteRole(
		c.Request().Context(),
		id,
	); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "data deleted successfully",
	})
}
