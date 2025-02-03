package rest

import (
	"context"
	"net/http"
	"santapan/domain"
	"santapan/internal/rest/middleware"
	"santapan/pkg/json"
	"strconv"

	encodingJson "encoding/json"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

//go:generate mockery --name ArticleService
type MenuService interface {
	Fetch(ctx context.Context, cursor string, num int64, search string) ([]domain.Menu, string, error)
	GetByID(ctx context.Context, id int64) (domain.Menu, error)
	GetByCategoryID(ctx context.Context, categoryID int64, search string) ([]domain.Menu, error)
}

// ArticleHandler  represent the httphandler for article
type MenuHandler struct {
	MenuService          MenuService
	PersonalisasiService PersonalisasiService
	Validator            *validator.Validate
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewMenuHandler(e *echo.Echo, menuService MenuService, personalisasiService PersonalisasiService) {
	validator := validator.New()
	// Register the custom validation function
	validator.RegisterValidation("date", validateDate)

	handler := &MenuHandler{
		MenuService:          menuService,
		PersonalisasiService: personalisasiService,
		Validator:            validator,
	}

	e.GET("/menu", handler.Fetch)
	e.GET("/menu/:id", handler.GetByID)
	e.GET("/menu/category/:category_id", handler.GetByCategoryID)
	e.GET("/menu/recommendation", handler.GetRecommendation, middleware.AuthMiddleware)
}

func (ah *MenuHandler) Fetch(c echo.Context) error {
	cursor := c.QueryParam("cursor")
	num := c.QueryParam("num")
	search := c.QueryParam("search")

	parseNum, err := strconv.Atoi(num)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid Num")
	}

	menus, nextCursor, err := ah.MenuService.Fetch(c.Request().Context(), cursor, int64(parseNum), search)

	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	responseData := map[string]interface{}{
		"menus":      menus,
		"nextCursor": nextCursor,
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Menu!", responseData)
}

func (ah *MenuHandler) GetRecommendation(c echo.Context) error {
	// Parse query parameters
	cursor := c.QueryParam("cursor")
	num := c.QueryParam("num")
	search := c.QueryParam("search")

	// Validate and parse the 'num' parameter to determine how many items to fetch
	parseNum, err := strconv.Atoi(num)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "Invalid 'num' parameter", "Please provide a valid number.")
	}

	// Fetch menu items with pagination
	menus, nextCursor, err := ah.MenuService.Fetch(c.Request().Context(), cursor, int64(parseNum), search)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Error fetching menus", err.Error())
	}

	// Get the personalization data for the user
	userID, ok := c.Get("userID").(int64)
	if !ok {
		return json.Response(c, http.StatusUnauthorized, false, "Unauthorized", "User ID missing or invalid in context.")
	}

	personalisasi, err := ah.PersonalisasiService.GetByUserID(c.Request().Context(), userID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "Error fetching personalization data", err.Error())
	}

	// Filter menus based on personalization data
	filteredMenus := make([]domain.Menu, 0)

	for _, menu := range menus {
		var features map[string]interface{}
		if err := encodingJson.Unmarshal(menu.Features, &features); err != nil {
			return json.Response(c, http.StatusInternalServerError, false, "Error unmarshalling features", "Failed to unmarshal menu features.")
		}

		if (personalisasi.Diabetes || personalisasi.RendahGula) && features != nil {
			if isGlutenFree, ok := features["gluten_free"].(bool); ok && isGlutenFree {
				filteredMenus = append(filteredMenus, menu)
			}
		}

		if personalisasi.TinggiProtein && features != nil {
			if protein, ok := features["high_protein"].(float64); ok && protein >= 20 {
				filteredMenus = append(filteredMenus, menu)
			}
		}

		if personalisasi.Vegetarian && features != nil {
			if isVegetarian, ok := features["vegetarian"].(bool); ok && isVegetarian {
				filteredMenus = append(filteredMenus, menu)
			}
		}
	}

	// Prepare the response data
	responseData := map[string]interface{}{
		"menus":      filteredMenus, // Return the filtered menus
		"nextCursor": nextCursor,
	}

	return json.Response(c, http.StatusOK, true, "Successfully fetched menu recommendations", responseData)
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
	search := c.QueryParam("search")

	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid Category ID")
	}

	menus, err := ah.MenuService.GetByCategoryID(c.Request().Context(), int64(categoryID), search)
	if err != nil {
		return json.Response(c, http.StatusNotFound, false, "", err.Error())
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Articles!", menus)
}
