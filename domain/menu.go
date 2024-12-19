package domain

import (
	"encoding/json"
	"time"
)

type Menu struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"` // Allows NULL values
	Price       float64         `json:"price"`       // Change from int64 to float64
	ImageURL    string          `json:"image_url"`   // Allows NULL values
	Nutrition   json.RawMessage `json:"nutrition"`   // Decoded JSON
	Features    json.RawMessage `json:"features"`    // Decoded JSON
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
