package domain

import (
	"time"
)

// Nutrition represents the model for food nutrition information
type Nutrition struct {
	ID            int64     `json:"id"`
	FoodName      string    `json:"food_name"`
	Calories      int       `json:"calories"`
	Protein       int       `json:"protein"`
	Fat           int       `json:"fat"`
	Carbohydrates int       `json:"carbohydrates"`
	Sugar         int       `json:"sugar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
