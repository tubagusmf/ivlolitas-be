package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/usecase"
)

type productVariantHandler struct {
	productVariantUsecase usecase.IProductVariantUsecase
}

func NewProductVariantHandler(e *echo.Echo, uc usecase.IProductVariantUsecase, auth *AuthMiddleware) {
	handler := &productVariantHandler{
		productVariantUsecase: uc,
	}

	group := e.Group("/v1/products")

	group.POST("/:productId/variants", handler.CreateProductVariant, auth.JWT)
	group.GET("/:productId/variants", handler.GetProductVariants)
	group.GET("/variants/:variantId", handler.GetProductVariantByID)
	group.PUT("/variants/:variantId", handler.UpdateProductVariant, auth.JWT)
	group.DELETE("/variants/:variantId", handler.DeleteProductVariant, auth.JWT)
}

// CreateProductVariant
//
// @Summary Create a new product variant
// @Description Create a new product variant
// @Tags Product Variants
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param productId path string true "Product UUID"
// @Param request body model.ProductVariantInput true "Product variant input"
// @Success 201 {object} model.ProductVariant
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{productId}/variants [post]
func (h *productVariantHandler) CreateProductVariant(c echo.Context) error {
	productID := c.Param("productId")

	if productID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid product id",
		})
	}

	var body model.ProductVariantInput

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	body.ProductID = productID

	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	productVariant, err := h.productVariantUsecase.Create(
		c.Request().Context(),
		&body,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "product variant created successfully",
		"data":    productVariant,
	})
}

// GetProductVariants
//
// @Summary Get all product variants
// @Description Get all product variants
// @Tags Product Variants
// @Accept json
// @Produce json
// @Param productId path string true "Product UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{productId}/variants [get]
func (h *productVariantHandler) GetProductVariants(c echo.Context) error {
	productID := c.Param("productId")
	if productID == "" {
		return c.JSON(http.StatusBadRequest, "invalid product id")
	}

	productVariants, err := h.productVariantUsecase.GetByProductID(c.Request().Context(), productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "product variants fetched successfully",
		"data":    productVariants,
	})
}

// GetProductVariantByID
//
// @Summary Get a product variant by ID
// @Description Get a product variant by ID
// @Tags Product Variants
// @Accept json
// @Produce json
// @Param id path string true "Product variant UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/variants/{id} [get]
func (h *productVariantHandler) GetProductVariantByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "invalid product variant id")
	}

	productVariant, err := h.productVariantUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "product variant fetched successfully",
		"data":    productVariant,
	})
}

// UpdateProductVariant
//
// @Summary Update a product variant by ID
// @Description Update a product variant by ID
// @Tags Product Variants
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product variant UUID"
// @Param request body model.ProductVariantUpdateInput true "Product variant update input"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/variants/{id} [put]
func (h *productVariantHandler) UpdateProductVariant(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid product variant id",
		})
	}

	var body model.ProductVariantUpdateInput

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	productVariant, err := h.productVariantUsecase.Update(
		c.Request().Context(),
		id,
		&body,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "product variant updated successfully",
		"data":    productVariant,
	})
}

// DeleteProductVariant
//
// @Summary Delete a product variant by ID
// @Description Delete a product variant by ID
// @Tags Product Variants
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product variant UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/variants/{id} [delete]
func (h *productVariantHandler) DeleteProductVariant(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid product variant id",
		})
	}

	if err := h.productVariantUsecase.Delete(
		c.Request().Context(),
		id,
	); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "product variant deleted successfully",
	})
}
