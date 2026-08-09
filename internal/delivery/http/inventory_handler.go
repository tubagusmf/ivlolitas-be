package http

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/ivlolitas-be/internal/model"
	"github.com/tubagusmf/ivlolitas-be/internal/usecase"
)

type inventoryHandler struct {
	inventoryUsecase usecase.IInventoryUsecase
}

func NewInventoryHandler(e *echo.Echo, uc usecase.IInventoryUsecase, auth *AuthMiddleware) {
	handler := &inventoryHandler{
		inventoryUsecase: uc,
	}

	group := e.Group("/v1/inventories")

	group.GET("/:productVariantID", handler.GetInventory, auth.JWT)
	group.GET("/:productVariantID/transactions", handler.GetTransactions, auth.JWT)
	group.POST("/restock", handler.Restock, auth.JWT)
	group.POST("/sale", handler.Sale, auth.JWT)
	group.POST("/return", handler.Return, auth.JWT)
	group.POST("/damage", handler.Damage, auth.JWT)
	group.POST("/release", handler.Release, auth.JWT)
	group.POST("/adjustment", handler.Adjustment, auth.JWT)
}

// GetInventory
//
//	@Summary		Get Inventory
//	@Description	Get current inventory by product variant ID
//	@Tags			Inventory
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			productVariantID	path	string	true	"Product Variant ID"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		400					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Failure		404					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/inventories/{productVariantID} [get]
func (h *inventoryHandler) GetInventory(c echo.Context) error {
	productVariantID := c.Param("productVariantID")

	if _, err := uuid.Parse(productVariantID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid product variant ID",
		})
	}

	inventory, err := h.inventoryUsecase.GetInventory(
		c.Request().Context(),
		productVariantID,
	)

	if err != nil {
		if errors.Is(err, usecase.ErrInventoryNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "inventory not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed to get inventory",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": inventory,
	})
}

// GetTransactions
//
//	@Summary		Get Inventory Transactions
//	@Description	Get inventory transaction history by product variant ID
//	@Tags			Inventory
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			productVariantID	path	string	true	"Product Variant ID"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		400					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Failure		404					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/inventories/{productVariantID}/transactions [get]
func (h *inventoryHandler) GetTransactions(c echo.Context) error {
	productVariantID := c.Param("productVariantID")

	if _, err := uuid.Parse(productVariantID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid product variant ID",
		})
	}

	transactions, err := h.inventoryUsecase.GetTransactions(
		c.Request().Context(),
		productVariantID,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": transactions,
	})
}

// Restock
//
//	@Summary		Restock
//	@Description	Restock
//	@Tags			Inventory
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request body model.InventoryOperationInput true "Restock Request"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		400					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Failure		404					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/inventories/restock [post]
func (h *inventoryHandler) Restock(c echo.Context) error {
	var req model.InventoryOperationInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "unauthorized",
		})
	}

	err := h.inventoryUsecase.Restock(
		c.Request().Context(),
		req.ProductVariantID,
		req.Quantity,
		userID,
	)

	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInventoryNotFound):
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "inventory not found",
			})

		case errors.Is(err, usecase.ErrInsufficientStock):
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"message": err.Error(),
			})

		case errors.Is(err, usecase.ErrInvalidQuantity):
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
			})

		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "failed to restock inventory",
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "stock successfully restocked",
	})
}

// Sale
//
//	@Summary		Sale
//	@Description	Sale
//	@Tags			Inventory
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request body model.InventoryOperationInput true "Sale Request"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		400					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Failure		404					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/inventories/sale [post]
func (h *inventoryHandler) Sale(c echo.Context) error {
	var req model.InventoryOperationInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	userID, ok := c.Get("user_id").(string)

	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "unauthorized",
		})
	}

	err := h.inventoryUsecase.Sale(
		c.Request().Context(),
		req.ProductVariantID,
		req.Quantity,
		userID,
	)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "stock successfully deducted",
	})
}

// Return
//
//	@Summary		Return
//	@Description	Return
//	@Tags			Inventory
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request body model.InventoryOperationInput true "Return Request"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		400					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Failure		404					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/inventories/return [post]
func (h *inventoryHandler) Return(c echo.Context) error {
	var req model.InventoryOperationInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	userID, ok := c.Get("user_id").(string)

	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "unauthorized",
		})
	}

	err := h.inventoryUsecase.Return(
		c.Request().Context(),
		req.ProductVariantID,
		req.Quantity,
		userID,
	)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "stock successfully returned",
	})
}

// Damage
//
//	@Summary		Damage
//	@Description	Damage
//	@Tags			Inventory
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request body model.InventoryOperationInput true "Damage Request"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		400					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Failure		404					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/inventories/damage [post]
func (h *inventoryHandler) Damage(c echo.Context) error {
	var req model.InventoryOperationInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	userID, ok := c.Get("user_id").(string)

	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "unauthorized",
		})
	}

	err := h.inventoryUsecase.Damage(
		c.Request().Context(),
		req.ProductVariantID,
		req.Quantity,
		userID,
	)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "stock successfully damaged",
	})
}

// Release
//
//	@Summary		Release
//	@Description	Release
//	@Tags			Inventory
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request body model.InventoryOperationInput true "Release Request"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		400					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Failure		404					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/inventories/release [post]
func (h *inventoryHandler) Release(c echo.Context) error {
	var req model.InventoryOperationInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	userID, ok := c.Get("user_id").(string)

	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "unauthorized",
		})
	}

	err := h.inventoryUsecase.Release(
		c.Request().Context(),
		req.ProductVariantID,
		req.Quantity,
		userID,
	)

	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInventoryNotFound):
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "inventory not found",
			})

		case errors.Is(err, usecase.ErrInsufficientReserve):
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"message": err.Error(),
			})

		case errors.Is(err, usecase.ErrInvalidQuantity):
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
			})

		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "failed to release reserved stock",
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "reserved stock successfully released",
	})
}

// Adjustment
//
//	@Summary		Adjustment
//	@Description	Adjustment
//	@Tags			Inventory
//	@Security       BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request body model.InventoryAdjustmentInput true "Adjustment Request"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		400					{object}	map[string]interface{}
//	@Failure		401					{object}	map[string]interface{}
//	@Failure		404					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/inventories/adjustment [post]
func (h *inventoryHandler) Adjustment(c echo.Context) error {
	var req model.InventoryAdjustmentInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	userID, ok := c.Get("user_id").(string)

	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "unauthorized",
		})
	}

	err := h.inventoryUsecase.Adjustment(
		c.Request().Context(),
		req.ProductVariantID,
		req.Quantity,
		userID,
	)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "stock successfully adjusted",
	})
}
