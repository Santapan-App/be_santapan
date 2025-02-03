package rest

import (
	"context"
	"net/http"
	"santapan/domain"
	"santapan/pkg/json"

	"github.com/labstack/echo/v4"
)

//go:generate mockery --name ArticleService
type NutritionService interface {
	GetByClassification(ctx context.Context, classification string) (domain.Nutrition, error)
}

// NutritionHandler  represent the httphandler for article

type NutritionHandler struct {
	NutritionService NutritionService
}

func NewNutritionHandler(e *echo.Echo, nutritionService NutritionService) {
	handler := &NutritionHandler{
		NutritionService: nutritionService,
	}

	e.GET("/nutrition", handler.GetByClassification)
}

func (ah *NutritionHandler) GetByClassification(c echo.Context) error {
	classification := c.QueryParam("classification")

	nutrition, err := ah.NutritionService.GetByClassification(c.Request().Context(), classification)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Successfully get nutrition by classification!", nutrition)
}
