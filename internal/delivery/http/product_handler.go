package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/usecase"
)

type productHandler struct {
	productUsecase usecase.IProductUsecase
}

func NewProductHandler(e *echo.Echo, uc usecase.IProductUsecase, auth *AuthMiddleware) {
	handler := &productHandler{
		productUsecase: uc,
	}

	group := e.Group("/v1/products")

	group.POST("", handler.CreateProduct, auth.JWT)
	group.GET("", handler.GetProducts, auth.JWT)
	group.GET("/:id", handler.GetProductByID, auth.JWT)
	group.PUT("/:id", handler.UpdateProduct, auth.JWT)
	group.DELETE("/:id", handler.DeleteProduct, auth.JWT)
}

// CreateProduct godoc
//
// @Summary Create a new product
// @Description Create a new product
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.ProductInput true "Product input"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products [post]
func (h *productHandler) CreateProduct(c echo.Context) error {
	var body model.ProductInput

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	// TODO: validation
	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	product, err := h.productUsecase.CreateProduct(c.Request().Context(), &body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "product created successfully",
		"data":    product,
	})
}

// GetProducts godoc
//
// @Summary Get all products
// @Description Get all products
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search by name or sku"
// @Param is_active query bool false "Filter by is_active"
// @Param sort query string false "Sort by created_at or updated_at"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products [get]
func (h *productHandler) GetProducts(c echo.Context) error {
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

	var isActive *bool

	if value := c.QueryParam("is_active"); value != "" {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid is_active")
		}

		isActive = &b
	}

	filter := &model.ProductFilter{
		Search:   search,
		Sort:     sort,
		Page:     page,
		Limit:    limit,
		IsActive: isActive,
	}

	products, err := h.productUsecase.GetProducts(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "products fetched successfully",
		"data":    products,
	})
}

// GetProductByID godoc
//
// @Summary Get a product by ID
// @Description Get a product by ID
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product UUID"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [get]
func (h *productHandler) GetProductByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "invalid product id")
	}

	product, err := h.productUsecase.GetProductByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "product fetched successfully",
		"data":    product,
	})
}

// UpdateProduct godoc
//
// @Summary Update a product by ID
// @Description Update a product by ID
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product UUID"
// @Param request body model.ProductInput true "Product input"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [put]
func (h *productHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "invalid product id")
	}

	var body model.ProductUpdateInput

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	// TODO: validation
	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	product, err := h.productUsecase.UpdateProduct(c.Request().Context(), id, &body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "product updated successfully",
		"data":    product,
	})
}

// DeleteProduct godoc
//
// @Summary Delete a product by ID
// @Description Delete a product by ID
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product UUID"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [delete]
func (h *productHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "invalid product id")
	}

	if err := h.productUsecase.DeleteProduct(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "product deleted successfully",
	})
}
