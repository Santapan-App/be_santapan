package domain

import (
	"time"
)

// Bundling represents the bundling table
type Bundling struct {
	ID           int64     `json:"id"`
	ImageURL     string    `json:"image_url"`
	BundlingName string    `json:"bundling_name"`
	BundlingType string    `json:"bundling_type"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BundlingMenu represents the bundling_menu table
type BundlingMenu struct {
	ID              int64     `json:"id"`
	DayNumber       int       `json:"day_number"`
	MealDescription string    `json:"meal_description"` // Allows NULL values
	Bundling        Bundling  `json:"bundling"`
	Menu            Menu      `json:"menu"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
