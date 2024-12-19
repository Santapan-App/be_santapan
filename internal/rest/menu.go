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

//go:generate mockery --name ArticleService
type MenuService interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]domain.Menu, string, error)
	GetByID(ctx context.Context, id int64) (domain.Menu, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]domain.Menu, error)
}

// ArticleHandler  represent the httphandler for article
type MenuHandler struct {
	MenuService MenuService
	Validator   *validator.Validate
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewMenuHandler(e *echo.Echo, menuService MenuService) {
	validator := validator.New()
	// Register the custom validation function
	validator.RegisterValidation("date", validateDate)

	handler := &MenuHandler{
		MenuService: menuService,
		Validator:   validator,
	}

	e.GET("/menu", handler.Fetch)
	e.GET("/menu/:id", handler.GetByID)
	e.GET("/menu/category/:category_id", handler.GetByCategoryID)
}

func (ah *MenuHandler) Fetch(c echo.Context) error {
	cursor := c.QueryParam("cursor")
	num := c.QueryParam("num")

	parseNum, err := strconv.Atoi(num)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid Num")
	}

	menus, nextCursor, err := ah.MenuService.Fetch(c.Request().Context(), cursor, int64(parseNum))

	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	responseData := map[string]interface{}{
		"menus":      menus,
		"nextCursor": nextCursor,
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Articles!", responseData)
}

func (ah *MenuHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid ID")
	}

	menu, err := ah.MenuService.GetByID(c.Request().Context(), int64(id))
	if err != nil {
		return json.Response(c, http.StatusNotFound, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Article!", menu)
}

func (ah *MenuHandler) GetByCategoryID(c echo.Context) error {
	categoryID, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid Category ID")
	}

	menus, err := ah.MenuService.GetByCategoryID(c.Request().Context(), int64(categoryID))
	if err != nil {
		return json.Response(c, http.StatusNotFound, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Articles!", menus)
}
