package rest

import (
	"context"
	"net/http"
	"santapan_transaction_service/domain"
	"santapan_transaction_service/internal/rest/middleware"
	"santapan_transaction_service/pkg/json"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

//go:generate mockery --name ArticleService
type CartService interface {
	GetByUserID(ctx context.Context, userID int64, status string) (domain.Cart, error)
	GetByID(ctx context.Context, cartID int64) (domain.Cart, error)
	Validate(ctx context.Context, userID int64, cartID int64) error
	Store(ctx context.Context, cart *domain.Cart) error
	Update(ctx context.Context, cart *domain.Cart) error          // Explicit update method
	Delete(ctx context.Context, userID int64, cartID int64) error // New method for deleting a cart
}

type CartItemService interface {
	GetByCartID(ctx context.Context, cartID int64) ([]domain.CartItem, error)
	GetByID(ctx context.Context, id int64) (domain.CartItem, error)
	Delete(ctx context.Context, id int64) error
	Store(ctx context.Context, item *domain.CartItem) error
	Update(ctx context.Context, item *domain.CartItem) error
}

// ArticleHandler  represent the httphandler for article
type CartHandler struct {
	CartService     CartService
	CartItemService CartItemService
	Validator       *validator.Validate
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewCartHandler(e *echo.Echo, cartService CartService, cartItemService CartItemService) {
	validator := validator.New()

	handler := &CartHandler{
		CartService:     cartService,
		CartItemService: cartItemService,
		Validator:       validator,
	}

	e.GET("/cart", handler.Fetch, middleware.AuthMiddleware)
	e.GET("/cart/:id", handler.GetByID, middleware.AuthMiddleware)
	e.GET("/cart/home", handler.GetCartHome, middleware.AuthMiddleware)
	e.POST("/cart", handler.Store, middleware.AuthMiddleware)
	e.POST("/cart/:id", handler.UpdateCartItem, middleware.AuthMiddleware)
	e.DELETE("/cart/item/:id", handler.DeleteItem, middleware.AuthMiddleware)
}

// Fetch Cart Based On Token User ID
func (a *CartHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Retrieve user ID from the context (assuming it's stored in the token claims)
	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	cart, err := a.CartService.GetByUserID(ctx, userID, "active")
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", err.Error())
	}

	// Fetch the cart items
	items, err := a.CartItemService.GetByCartID(ctx, cart.ID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to fetch cart items", nil)
	}

	responseData := map[string]interface{}{
		"cart":  cart,
		"items": items,
	}

	return json.Response(c, http.StatusOK, true, "Cart Fetched", responseData)
}

// GetCartHome fetches the user's cart items and the total amount.
func (a *CartHandler) GetCartHome(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Retrieve user ID from the context (assuming it's stored in the token claims)
	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	// Fetch the user's cart
	cart, err := a.CartService.GetByUserID(ctx, userID, "active")
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	// Fetch the cart items
	items, err := a.CartItemService.GetByCartID(ctx, cart.ID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to fetch cart items", nil)
	}

	// Calculate the total amount
	var totalAmount float64
	var totalQuantity int

	for _, item := range items {
		totalAmount += item.Price
		totalQuantity += item.Quantity
	}

	responseData := map[string]interface{}{
		"total_amount":   totalAmount,
		"total_quantity": totalQuantity,
	}

	return json.Response(c, http.StatusOK, true, "Cart items fetched successfully", responseData)
}

// Store accepts cart items, creates a cart if one doesn't exist, and inserts the items into it.
func (a *CartHandler) Store(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Retrieve user ID from the context (assuming it's stored in the token claims)
	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	// Bind request body to a list of cart items
	var items domain.CartItem
	if err := c.Bind(&items); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid request body", nil)
	}

	// Validate each cart item
	if err := a.Validator.Struct(items); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid cart item fields", nil)
	}

	// Check if the user already has a cart
	cart, err := a.CartService.GetByUserID(ctx, userID, "active")

	if err != nil {
		// Handle error (e.g., cart not found)
		if err != domain.ErrNotFound {
			return json.Response(c, http.StatusInternalServerError, false, "Failed to fetch the cart", nil)
		}

		// Create a new cart if one doesn't exist
		cart = domain.Cart{
			UserID:     userID,
			TotalPrice: items.Price,
			Status:     "active",
		}

		err = a.CartService.Store(ctx, &cart)

		if err != nil {
			return json.Response(c, http.StatusInternalServerError, false, "Failed to create a new cart", nil)
		}
	}

	items.CartID = cart.ID // Associate each item with the user's cart
	if err := a.CartItemService.Store(ctx, &items); err != nil {
		logrus.Info("Failed to store cart items:", err)
		return json.Response(c, http.StatusInternalServerError, false, "Failed to store cart items", nil)
	}

	return json.Response(c, http.StatusCreated, true, "Cart items stored successfully", nil)
}

// Delete deletes a cart item based on its ID and ensures the user owns it.
func (a *CartHandler) DeleteItem(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Retrieve user ID from the context (assuming it's stored in the token claims)
	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	// Get the cart item ID from the URL parameter
	itemIDStr := c.Param("id") // Assuming cart item ID is passed as a URL parameter
	if itemIDStr == "" {
		return json.Response(c, http.StatusBadRequest, false, "Cart item ID is required", nil)
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid cart item ID", nil)
	}

	// Ensure the user owns the cart item
	_, err = a.CartService.GetByUserID(ctx, userID, "active") // Fetch the cart the item belongs to
	if err != nil {
		return json.Response(c, http.StatusForbidden, false, "You do not own this cart item", nil)
	}

	// Delete the cart item
	err = a.CartItemService.Delete(ctx, itemID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to delete the cart item", nil)
	}

	return json.Response(c, http.StatusOK, true, "Cart item deleted successfully", nil)
}

// Delete deletes a cart based on its ID and ensures the user owns it.
func (a *CartHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Retrieve user ID from the context (assuming it's stored in the token claims)
	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	// Get the cart ID from the URL parameter
	cartIDStr := c.Param("id") // Assuming cart ID is passed as a URL parameter
	if cartIDStr == "" {
		return json.Response(c, http.StatusBadRequest, false, "Cart ID is required", nil)
	}

	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid cart ID", nil)
	}

	// Ensure the user owns the cart
	err = a.CartService.Validate(ctx, userID, cartID)
	if err != nil {
		return json.Response(c, http.StatusForbidden, false, "You do not own this cart", nil)
	}

	// Delete the cart
	err = a.CartService.Delete(ctx, userID, cartID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to delete the cart", nil)
	}

	return json.Response(c, http.StatusOK, true, "Cart deleted successfully", nil)
}

// UpdateCartItem updates the quantity of a specific cart item based on its ID.
func (a *CartHandler) UpdateCartItem(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Retrieve user ID from the context (assuming it's stored in the token claims)
	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	// Get the cart item ID from the URL parameter
	itemIDStr := c.Param("id")
	if itemIDStr == "" {
		return json.Response(c, http.StatusBadRequest, false, "Cart item ID is required", nil)
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid cart item ID", nil)
	}

	// Fetch the existing cart item to ensure it exists and belongs to the user
	existingItem, err := a.CartItemService.GetByID(ctx, itemID)
	if err != nil {
		return json.Response(c, http.StatusNotFound, false, "Cart item not found", nil)
	}

	// Ensure the user owns the cart item (implement this in your service layer)
	if err := a.CartService.Validate(ctx, userID, existingItem.CartID); err != nil {
		return json.Response(c, http.StatusForbidden, false, "You do not own this cart item", nil)
	}

	// Bind only the quantity field from the request body
	var updateRequest struct {
		Quantity int `json:"quantity"` // Only quantity is expected
	}

	if err := c.Bind(&updateRequest); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid request body", nil)
	}

	// Validate that the quantity is greater than 0
	if updateRequest.Quantity <= 0 {
		return json.Response(c, http.StatusBadRequest, false, "Quantity must be greater than 0", nil)
	}

	// Set the new quantity for the cart item
	existingItem.Quantity = updateRequest.Quantity

	// Update the cart item in the service
	err = a.CartItemService.Update(ctx, &existingItem)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to update the cart item", nil)
	}

	return json.Response(c, http.StatusOK, true, "Cart item updated", nil)
}

// GetByID fetches a cart based on its ID and ensures the user owns it.
func (a *CartHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Get the cart ID from the URL parameter
	cartIDStr := c.Param("id") // Assuming cart ID is passed as a URL parameter

	if cartIDStr == "" {
		return json.Response(c, http.StatusBadRequest, false, "Cart ID is required", nil)
	}

	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid cart ID", nil)
	}

	// Fetch the cart
	cart, err := a.CartService.GetByID(ctx, cartID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to fetch the cart", nil)
	}

	// Fetch the cart items
	items, err := a.CartItemService.GetByCartID(ctx, cart.ID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to fetch cart items", nil)
	}

	responseData := map[string]interface{}{
		"cart":  cart,
		"items": items,
	}

	return json.Response(c, http.StatusOK, true, "Cart fetched successfully", responseData)
}
