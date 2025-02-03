package rest

import (
	"context"
	"net/http"
	"santapan/domain"
	"santapan/pkg/json"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type CourierService interface {
	Fetch(ctx context.Context) ([]domain.Courier, error)
	FetchByID(ctx context.Context, id int) (*domain.Courier, error)
}

// CourierHandler  represent the httphandler for courier
type CourierHandler struct {
	CourierService CourierService
	Validator      *validator.Validate
}

// NewCourierHandler will initialize the courier/ resources endpoint
func NewCourierHandler(e *echo.Echo, courierService CourierService) {
	validator := validator.New()
	// Register the custom validation function
	validator.RegisterValidation("date", validateDate)

	handler := &CourierHandler{
		CourierService: courierService,
		Validator:      validator,
	}

	e.GET("/courier", handler.Fetch)
	e.GET("/courier/:id", handler.GetByID)
}

func (ah *CourierHandler) Fetch(c echo.Context) error {

	couriers, err := ah.CourierService.Fetch(c.Request().Context())

	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	responseData := map[string]interface{}{
		"couriers": couriers,
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Articles!", responseData)
}

func (ah *CourierHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	// convert id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid ID")
	}
	courier, err := ah.CourierService.FetchByID(c.Request().Context(), idInt)

	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	responseData := map[string]interface{}{
		"courier": courier,
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Courier!", responseData)
}
