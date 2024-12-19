package rest

import (
	"context"
	"fmt"
	"net/http"
	"santapan/domain"
	"santapan/pkg/json"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type BundlingService interface {
	GetByID(ctx context.Context, id int64) (domain.Bundling, error)
	Fetch(ctx context.Context, cursor string, num int64) ([]domain.Bundling, string, error)
	FetchBundlingMenuByBundlingID(ctx context.Context, bundlingID int64) ([]domain.BundlingMenu, error)
}

// BundlingHandler manages HTTP endpoints for bundling
type BundlingHandler struct {
	BundlingService BundlingService
	Validator       *validator.Validate
}

// NewBundlingHandler initializes bundling routes
func NewBundlingHandler(e *echo.Echo, bundlingService BundlingService) {
	handler := &BundlingHandler{
		BundlingService: bundlingService,
		Validator:       validator.New(),
	}

	e.GET("/bundling/:id", handler.GetByID)
	e.GET("/bundling", handler.Fetch)
	e.GET("/bundling/:id/menu", handler.FetchBundlingMenu)
	e.GET("/bundling/:id/menu/grouped", handler.FetchBundlingMenuGroupedByWeekAndDay)
}

// GetBundling handles fetching a bundling by ID
func (bh *BundlingHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid bundling ID")
	}

	bundling, err := bh.BundlingService.GetByID(c.Request().Context(), id)
	if err != nil {
		return json.Response(c, http.StatusNotFound, false, "", "Bundling not found")
	}

	return json.Response(c, http.StatusOK, true, "Bundling fetched successfully", bundling)
}

// Fetch bundlings
func (bh *BundlingHandler) Fetch(c echo.Context) error {
	cursor := c.QueryParam("cursor")
	num := c.QueryParam("num")

	parseNum, err := strconv.Atoi(num)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid Num")
	}

	bundlings, nextCursor, err := bh.BundlingService.Fetch(c.Request().Context(), cursor, int64(parseNum))
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", err.Error())
	}

	responseData := map[string]interface{}{
		"bundlings":  bundlings,
		"nextCursor": nextCursor,
	}

	return json.Response(c, http.StatusOK, true, "Successfully Get Articles!", responseData)
}

func (bh *BundlingHandler) FetchBundlingMenu(c echo.Context) error {
	bundlingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid bundling ID")
	}

	bundlingMenus, err := bh.BundlingService.FetchBundlingMenuByBundlingID(c.Request().Context(), bundlingID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", "Failed to fetch bundling menu")
	}

	// Group by DayNumber
	groupedData := make(map[int][]domain.BundlingMenu)
	for _, menu := range bundlingMenus {
		day := menu.DayNumber
		groupedData[day] = append(groupedData[day], menu)
	}

	// Sort day numbers
	var days []int
	for day := range groupedData {
		days = append(days, day)
	}
	sort.Ints(days)

	// Build response in order
	response := []map[string]interface{}{}
	for _, day := range days {
		dayGroup := map[string]interface{}{
			"day":  fmt.Sprintf("Hari %d", day),
			"menu": groupedData[day],
		}
		response = append(response, dayGroup)
	}

	return json.Response(c, http.StatusOK, true, "Bundling menu fetched successfully", response)
}

func (bh *BundlingHandler) FetchBundlingMenuGroupedByWeekAndDay(c echo.Context) error {
	bundlingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return json.Response(c, http.StatusBadRequest, false, "", "Invalid bundling ID")
	}

	// Fetch bundling details to check the type
	bundling, err := bh.BundlingService.GetByID(c.Request().Context(), bundlingID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", "Failed to fetch bundling details")
	}

	// Only allow if bundly_type is "monthly"
	if bundling.BundlingType != "monthly" {
		return json.Response(c, http.StatusBadRequest, false, "", "Bundling type must be 'monthly'")
	}

	// Fetch the menus
	bundlingMenus, err := bh.BundlingService.FetchBundlingMenuByBundlingID(c.Request().Context(), bundlingID)
	if err != nil {
		return json.Response(c, http.StatusInternalServerError, false, "", "Failed to fetch bundling menu")
	}

	// Group by week and day
	groupedByWeek := make(map[int]map[int][]domain.BundlingMenu)
	for _, menu := range bundlingMenus {
		week := (menu.DayNumber-1)/7 + 1
		if groupedByWeek[week] == nil {
			groupedByWeek[week] = make(map[int][]domain.BundlingMenu)
		}
		groupedByWeek[week][menu.DayNumber] = append(groupedByWeek[week][menu.DayNumber], menu)
	}

	// Build response
	response := []map[string]interface{}{}
	for week, days := range groupedByWeek {
		weekData := map[string]interface{}{
			"week": fmt.Sprintf("Minggu %d", week),
			"days": []map[string]interface{}{},
		}

		// Sort days in ascending order
		sortedDays := make([]int, 0, len(days))
		for day := range days {
			sortedDays = append(sortedDays, day)
		}
		sort.Ints(sortedDays)

		// Build day data in sorted order
		for _, day := range sortedDays {
			dayData := map[string]interface{}{
				"day":  fmt.Sprintf("Hari %d", day),
				"menu": days[day],
			}
			weekData["days"] = append(weekData["days"].([]map[string]interface{}), dayData)
		}
		response = append(response, weekData)
	}

	return json.Response(c, http.StatusOK, true, "Bundling menu grouped by week and day fetched successfully", response)
}
