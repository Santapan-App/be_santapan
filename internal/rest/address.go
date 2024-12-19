package rest

import (
	"context"
	"net/http"
	"santapan/domain"
	"santapan/internal/rest/middleware"
	"santapan/pkg/json"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

// PostgresRepositoryQueries defines the methods for querying the category repository.
type AddressService interface {
	GetByUserID(ctx context.Context, userID int64) ([]domain.Address, error)
	Create(ctx context.Context, a domain.Address) (domain.Address, error)
	Update(ctx context.Context, a domain.Address) (domain.Address, error)
	Delete(ctx context.Context, id int64) error
}

type AddressHandler struct {
	AddressService AddressService
	Validator      *validator.Validate
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewAddressHandler(e *echo.Echo, addressService AddressService) {
	validator := validator.New()
	// Register the custom validation function
	validator.RegisterValidation("date", validateDate)

	handler := &AddressHandler{
		AddressService: addressService,
		Validator:      validator,
	}

	e.GET("/address", handler.GetByUserID, middleware.AuthMiddleware)
	e.POST("/address", handler.Create, middleware.AuthMiddleware)
	e.PUT("/address", handler.Update, middleware.AuthMiddleware)
}

// GetByUserID retrieves an address by user ID.
func (ah *AddressHandler) GetByUserID(c echo.Context) error {
	userID := c.Get("userID").(int64)
	res, err := ah.AddressService.GetByUserID(c.Request().Context(), userID)
	logrus.Info(err)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}
	return json.Response(c, http.StatusOK, true, "Successfully Get Address!", res)
}

// Create creates a new address.
func (ah *AddressHandler) Create(c echo.Context) error {
	var a domain.Address
	if err := c.Bind(&a); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", err.Error())
	}

	if err := ah.Validator.Struct(a); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", err.Error())
	}

	userID := c.Get("userID").(int64)
	a.UserID = userID

	res, err := ah.AddressService.Create(c.Request().Context(), a)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Successfully Create Address!", res)
}

// Update updates an address.
func (ah *AddressHandler) Update(c echo.Context) error {
	var a domain.Address
	if err := c.Bind(&a); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", err.Error())
	}

	if err := ah.Validator.Struct(a); err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", err.Error())
	}

	userID := c.Get("userID").(int64)
	a.UserID = userID

	res, err := ah.AddressService.Update(c.Request().Context(), a)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Successfully Update Address!", res)
}

// Delete deletes an address.
func (ah *AddressHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid ID")
	}

	err = ah.AddressService.Delete(c.Request().Context(), int64(id))
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Successfully Delete Address!", nil)
}
