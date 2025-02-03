package rest

import (
	"bytes"
	"context"
	encodingJson "encoding/json"
	"net/http"
	"os"
	"santapan_transaction_service/domain"
	"santapan_transaction_service/internal/rest/middleware"
	"santapan_transaction_service/pkg/json"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

//go:generate mockery --name TransactionService
type TransactionService interface {
	GetByUserID(ctx context.Context, userID int64) ([]domain.Transaction, error)
	GetByID(ctx context.Context, transactionID int64) (domain.Transaction, error)
	GetOngoing(ctx context.Context, userID int64) ([]domain.Transaction, error)
	Validate(ctx context.Context, userID int64, transactionID int64) error
	Store(ctx context.Context, transaction *domain.Transaction) error
	Update(ctx context.Context, transaction *domain.Transaction) error // Explicit update method
}

// TransactionHandler  represent the httphandler for transaction
type TransactionHandler struct {
	TransactionService TransactionService
	CartService        CartService
	Validator          *validator.Validate
}

// NewTransactionHandler will initialize the transactions/ resources endpoint
func NewTransactionHandler(e *echo.Echo, transactionService TransactionService, cartService CartService) {
	validator := validator.New()

	handler := &TransactionHandler{
		TransactionService: transactionService,
		CartService:        cartService,
		Validator:          validator,
	}

	e.GET("/transaction", handler.Fetch, middleware.AuthMiddleware)
	e.GET("/transaction/:id", handler.GetByID, middleware.AuthMiddleware)
	e.GET("/transaction/ongoing", handler.Ongoing, middleware.AuthMiddleware)
	e.POST("/transaction", handler.Store, middleware.AuthMiddleware)
}

// Fetch Transaction Based On Token User ID
func (a *TransactionHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Retrieve user ID from the context (assuming it's stored in the token claims)
	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	transaction, err := a.TransactionService.GetByUserID(ctx, userID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return json.Response(c, http.StatusOK, true, "Success", transaction)
}

// Store a new transaction
func (a *TransactionHandler) Store(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Bind and validate the transaction request
	var transactionBody domain.TransactionBody
	logrus.Info(transactionBody)
	if err := c.Bind(&transactionBody); err != nil {
		return json.Response(c, http.StatusUnprocessableEntity, false, "Invalid request", nil)
	}

	if err := a.Validator.Struct(transactionBody); err != nil {
		return json.Response(c, http.StatusUnprocessableEntity, false, "Invalid Validator Request!", nil)
	}

	// Retrieve user ID from the context
	userID, ok := c.Get("userID").(int64)
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	// Prepare the payment request
	paymentBody := domain.PaymentBody{
		Amount: transactionBody.Amount,
		Name:   transactionBody.ItemNames,  // Assuming this is a slice of strings
		Qty:    transactionBody.ItemQtys,   // Assuming this is a slice of int64
		Price:  transactionBody.ItemPrices, // Assuming this is a slice of float64
	}

	paymentRequestBody, err := encodingJson.Marshal(paymentBody)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to prepare payment request", nil)
	}

	// Make the POST request to the payment service
	url := os.Getenv("PAYMENT_URL") // Replace with the actual URL
	req, err := http.NewRequest("POST", url, bytes.NewReader(paymentRequestBody))
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to create payment request", nil)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Request().Header.Get("Authorization"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to process payment request", nil)
	}
	defer resp.Body.Close()

	// Parse the response from the payment service
	if resp.StatusCode != http.StatusCreated {
		return json.Response(c, http.StatusInternalServerError, false, "Payment service error", nil)
	}

	var paymentResponse domain.PaymentResponse
	if err := encodingJson.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Failed to parse payment response", nil)
	}

	if paymentResponse.Success == false {
		return json.Response(c, http.StatusInternalServerError, false, paymentResponse.Message, nil)
	}

	// Update the cart status
	cart, err := a.CartService.GetByUserID(ctx, userID, "active")
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	transaction := &domain.Transaction{
		UserID:    userID,
		CartID:    cart.ID, // Assuming this is the cart ID
		PaymentID: paymentResponse.Data.ID,
		CourierID: transactionBody.CourierID,
		AddressID: transactionBody.AddressID,
		Status:    "pending",
		Amount:    transactionBody.Amount,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.TransactionService.Store(ctx, transaction); err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	cart.Status = "used"
	if err := a.CartService.Update(ctx, &cart); err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return json.Response(c, http.StatusCreated, true, "Transaction created successfully", map[string]interface{}{
		"transaction": transaction,
		"payment_url": paymentResponse.Data.Url,
	})
}

// Ongoing Transaction Based On Token User ID
func (a *TransactionHandler) Ongoing(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	// Retrieve user ID from the context (assuming it's stored in the token claims)
	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
	}

	transaction, err := a.TransactionService.GetOngoing(ctx, userID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return json.Response(c, http.StatusOK, true, "Success", transaction)
}

func (a *TransactionHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	transactionID := c.Param("id")
	transactionIDInt, err := strconv.ParseInt(transactionID, 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid transaction ID", nil)
	}

	transaction, err := a.TransactionService.GetByID(ctx, transactionIDInt)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return json.Response(c, http.StatusOK, true, "Success", transaction)
}

// GetDetailTransaction get detail transaction
// func (a *TransactionHandler) GetDetailTransaction(c echo.Context) error {
// 	// Retrieve user ID from the context (assuming it's stored in the token claims)

// 	ctx := c.Request().Context()
// 	if ctx == nil {
// 		ctx = context.Background()
// 	}

// 	// Retrieve user ID from the context (assuming it's stored in the token claims)
// 	userID, ok := c.Get("userID").(int64) // Adjust according to how you're storing the user ID in the context
// 	if !ok {
// 		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", nil)
// 	}

// 	transactionID := c.Param("id")
// 	transactionIDInt, err := strconv.ParseInt(transactionID, 10, 64)
// 	if err != nil {
// 		return json.Response(c, http.StatusBadRequest, false, "Invalid transaction ID", nil)
// 	}

// 	err = a.TransactionService.Validate(ctx, userID, transactionIDInt)
// 	if err != nil {
// 		return json.Response(c, http.StatusForbidden, false, err.Error(), nil)
// 	}

// 	transaction, err := a.TransactionService.GetByID(ctx, transactionIDInt)
// 	if err != nil {
// 		return json.Response(c, http.StatusInternalServerError, false, err.Error(), nil)
// 	}

// 	return json.Response(c, http.StatusOK, true, "Success", transaction)
// }
