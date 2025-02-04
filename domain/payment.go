package domain

import "time"

type PaymentBody struct {
	Channel string    `json:"channel" validate:"required"`
	Method  string    `json:"method" validate:"required"`
	Amount  float64   `json:"amount" validate:"required"`
	Image   []string  `json:"image" validate:"required"`
	Name    []string  `json:"name" validate:"required,dive,required"` // Ensures each name is non-empty
	Qty     []int64   `json:"qty" validate:"required,dive,gt=0"`      // Ensures each qty is > 0
	Price   []float64 `json:"price" validate:"required,dive,gt=0"`    // Ensures each price is > 0
}

type PaymentResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    PaymentData `json:"data"`
}

type PaymentData struct {
	ID          int64     `json:"id"`
	ReferenceID string    `json:"reference_id"`
	SessionID   string    `json:"session_id"`
	UserID      int64     `json:"user_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Url         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
