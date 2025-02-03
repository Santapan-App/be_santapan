package rest

import (
	"context"
	"net/http"
	"santapan/domain"
	"santapan/internal/rest/middleware"
	"santapan/pkg/json"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

// PersonalisasiHandler handles HTTP requests related to personalisasi data.
type PersonalisasiHandler struct {
	PersonalisasiService PersonalisasiService
	Validator            *validator.Validate
}

// PersonalisasiService interface defines methods for managing personalisasi data.
type PersonalisasiService interface {
	GetByUserID(context.Context, int64) (domain.Personalisasi, error)
	InsertOrUpdate(context.Context, domain.Personalisasi) (domain.Personalisasi, error)
}

// NewPersonalisasiHandler initializes the PersonalisasiHandler with the service and validator.
func NewPersonalisasiHandler(e *echo.Echo, personalisasiService PersonalisasiService) {
	validator := validator.New()
	// Register the custom validation function if needed
	validator.RegisterValidation("date", validateDate)

	handler := &PersonalisasiHandler{
		PersonalisasiService: personalisasiService,
		Validator:            validator,
	}

	e.POST("/personalisasi", handler.Store, middleware.AuthMiddleware)
}

// Store handles the POST request to insert or update the personalisasi data.
func (ah *PersonalisasiHandler) Store(c echo.Context) error {
	userID := c.Get("userID").(int64)

	var personalisasi domain.Personalisasi

	// Bind the incoming request body to the domain.Personalisasi model
	err := c.Bind(&personalisasi)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", err.Error())
	}

	// Set the user ID from the JWT token
	personalisasi.UserID = userID

	// Call the InsertOrUpdate method from the service to insert or update the record
	updatedPersonalisasi, err := ah.PersonalisasiService.InsertOrUpdate(c.Request().Context(), personalisasi)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	// Prepare the response data
	responseData := map[string]interface{}{
		"personalisasi": updatedPersonalisasi,
	}

	// Return success response
	return json.Response(c, http.StatusOK, true, "Successfully stored or updated personalisasi!", responseData)
}
