package rest

import (
	"context"
	"net/http"
	"santapan_transaction_service/domain"
	"santapan_transaction_service/internal/rest/middleware"
	"santapan_transaction_service/pkg/json"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

//go:generate mockery --name ArticleService
type HistoryService interface {
	GetByUserID(ctx context.Context, userID int64) ([]domain.Transaction, error)
	GetByID(ctx context.Context, historyID int64) (domain.Transaction, error)
	Validate(ctx context.Context, userID int64, historyID int64) error
}

// HistoryHandler represent the httphandler for histories
type HistoryHandler struct {
	HistoryService HistoryService
	Validator      *validator.Validate
}

// NewHistoryHandler will initialize the histories/ resources endpoint
func NewHistoryHandler(e *echo.Echo, historyService HistoryService) {
	validator := validator.New()

	handler := &HistoryHandler{
		HistoryService: historyService,
		Validator:      validator,
	}

	e.GET("/history", handler.Fetch, middleware.AuthMiddleware)
	e.POST("/history/:id", handler.GetDetailHistory, middleware.AuthMiddleware)
}

// Fetch History Based On Token User ID
func (a *HistoryHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Retrieve user ID from the context (assuming it's stored in the token claims)
	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	history, err := a.HistoryService.GetByUserID(ctx, userID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return json.Response(c, http.StatusOK, true, "Success", history)
}

// GetDetailHistory get detail history
func (a *HistoryHandler) GetDetailHistory(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Retrieve user ID from the context (assuming it's stored in the token claims)
	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	historyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid history ID", nil)
	}

	err = a.HistoryService.Validate(ctx, userID, historyID)
	if err != nil {
		return json.Response(c, http.StatusForbidden, false, err.Error(), nil)
	}

	history, err := a.HistoryService.GetByID(ctx, historyID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return json.Response(c, http.StatusOK, true, "Success", history)
}
