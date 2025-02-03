package domain

import "time"

type Courier struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Logo      string    `json:"logo"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
