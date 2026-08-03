package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/usecase"
)

type categoryHandler struct {
	categoryUsecase usecase.ICategoryUsecase
}

func NewCategoryHandler(e *echo.Echo, uc usecase.ICategoryUsecase, auth *AuthMiddleware) {
	handler := &categoryHandler{
		categoryUsecase: uc,
	}

	group := e.Group("/v1/categories")

	group.POST("", handler.CreateCategory, auth.JWT)
	group.GET("", handler.GetCategories)
	group.GET("/:id", handler.GetCategoryByID)
	group.PUT("/:id", handler.UpdateCategory, auth.JWT)
	group.DELETE("/:id", handler.DeleteCategory, auth.JWT)
}

// CreateCategory godoc
//
//	@Summary		Create Category
//	@Description	Create new category
//	@Tags			Categories
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.CategoryInput	true	"Category"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Router			/categories [post]
func (h *categoryHandler) CreateCategory(c echo.Context) error {
	var body model.CategoryInput

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	// TODO: validation
	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	category, err := h.categoryUsecase.CreateCategory(
		c.Request().Context(),
		&body,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "data created successfully",
		"data":    category,
	})
}

// GetCategories godoc
//
//	@Summary		Get Categories
//	@Description	Get all categories
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			search	query		string	false	"Search Category"
//	@Param			sort	query		string	false	"Sort field (name,-name,created_at,-created_at)"
//	@Param			page	query		int		false	"Page number"
//	@Param			limit	query		int		false	"Page size"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Router			/categories [get]
func (h *categoryHandler) GetCategories(c echo.Context) error {
	search := c.QueryParam("search")
	sort := c.QueryParam("sort")

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	filter := &model.CategoryFilter{
		Search: search,
		Sort:   sort,
		Page:   page,
		Limit:  limit,
	}

	categories, err := h.categoryUsecase.GetCategories(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "data fetched successfully",
		"data":    categories,
	})
}

// GetCategoryByID godoc
//
//	@Summary		Get Category By ID
//	@Description	Get category detail
//	@Tags			Categories
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/categories/{id} [get]
func (h *categoryHandler) GetCategoryByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid category id")
	}

	category, err := h.categoryUsecase.GetCategoryByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "data fetched successfully",
		"data":    category,
	})
}

// UpdateCategory godoc
//
//	@Summary		Update Category
//	@Description	Update category by ID
//	@Tags			Categories
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Category ID"
//	@Param			request	body		model.CategoryUpdateInput	true	"Category"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Router			/categories/{id} [put]
func (h *categoryHandler) UpdateCategory(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid category id")
	}

	var body model.CategoryUpdateInput

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	// TODO: validation
	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	category, err := h.categoryUsecase.UpdateCategory(
		c.Request().Context(),
		id,
		&body,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "data updated successfully",
		"data":    category,
	})
}

// DeleteCategory godoc
//
//	@Summary		Delete Category
//	@Description	Delete category by ID
//	@Tags			Categories
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/categories/{id} [delete]
func (h *categoryHandler) DeleteCategory(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid category id")
	}

	if err := h.categoryUsecase.DeleteCategory(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "data deleted successfully",
	})
}
