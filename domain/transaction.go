package domain

import (
	"time"
)

type Transaction struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CartID    int64     `json:"cart_id"`
	PaymentID int64     `json:"payment_id"`
	CourierID int64     `json:"courier_id"`
	AddressID int64     `json:"address_id"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type TransactionBody struct {
	CourierID  int64     `json:"courier_id" validate:"required"`
	AddressID  int64     `json:"address_id" validate:"required"`
	Amount     float64   `json:"amount" validate:"required,gt=0"`
	ItemNames  []string  `json:"item_names" validate:"required,dive,required"` // Array of item names
	ItemQtys   []int64   `json:"item_qtys" validate:"required,dive,gt=0"`      // Array of item quantities
	ItemPrices []float64 `json:"item_prices" validate:"required,dive,gt=0"`    // Array of item prices
}
