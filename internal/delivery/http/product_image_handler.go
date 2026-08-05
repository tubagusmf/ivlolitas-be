package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ivlolitas-be/internal/usecase"
)

type productImageHandler struct {
	productImageUsecase usecase.IProductImageUsecase
}

func NewProductImageHandler(e *echo.Echo, productImageUsecase usecase.IProductImageUsecase, auth *AuthMiddleware) {
	handler := &productImageHandler{
		productImageUsecase: productImageUsecase,
	}

	group := e.Group("/v1/products", auth.JWT)

	group.POST("/:id/images", handler.UploadImages)
	group.GET("/:id/images", handler.GetImages)
	group.DELETE("/images/:imageId", handler.DeleteImage)
	group.PATCH("/images/:imageId/primary", handler.SetPrimaryImage)
}

// UploadImages godoc
//
//	@Summary		Upload Product Images
//	@Description	Upload one or multiple images for a product
//	@Tags			Product Images
//	@Security		BearerAuth
//	@Accept			mpfd
//	@Produce		json
//	@Param			id		path		string	true	"Product UUID"
//	@Param			images	formData	file	true	"Product Images (multiple files)"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/products/{id}/images [post]
func (h *productImageHandler) UploadImages(c echo.Context) error {
	productID := c.Param("id")

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid multipart form",
		})
	}

	files := form.File["images"]

	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "images is required",
		})
	}

	result, err := h.productImageUsecase.UploadImages(
		c.Request().Context(),
		productID,
		files,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Images uploaded successfully",
		"data":    result,
	})
}

// GetImages godoc
//
//	@Summary		Get Product Images
//	@Description	Get all images of a product
//	@Tags			Product Images
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Product UUID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/products/{id}/images [get]
func (h *productImageHandler) GetImages(c echo.Context) error {
	productID := c.Param("id")

	images, err := h.productImageUsecase.GetImagesByProductID(
		c.Request().Context(),
		productID,
	)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    images,
	})
}

// DeleteImage godoc
//
//	@Summary		Delete Product Image
//	@Description	Delete product image by image id
//	@Tags			Product Images
//	@Security		BearerAuth
//	@Produce		json
//	@Param			imageId	path		string	true	"Image UUID"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/products/images/{imageId} [delete]
func (h *productImageHandler) DeleteImage(c echo.Context) error {
	imageID := c.Param("imageId")

	err := h.productImageUsecase.DeleteImage(
		c.Request().Context(),
		imageID,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Image deleted successfully",
	})
}

// SetPrimaryImage godoc
//
//	@Summary		Set Primary Product Image
//	@Description	Set selected image as primary image
//	@Tags			Product Images
//	@Security		BearerAuth
//	@Produce		json
//	@Param			imageId	path		string	true	"Image UUID"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/products/images/{imageId}/primary [patch]
func (h *productImageHandler) SetPrimaryImage(c echo.Context) error {
	imageID := c.Param("imageId")

	err := h.productImageUsecase.SetPrimaryImage(
		c.Request().Context(),
		imageID,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Primary image updated successfully",
	})
}
