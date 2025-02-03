package domain

import (
	"time"
)

type PersonalisasiDiseaseBody struct {
	Diabetes  bool `json:"diabetes"`
	GERD      bool `json:"gerd"`
	AsamUrat  bool `json:"asam_urat"`
	Kolestrol bool `json:"kolestrol"`
}

// Personalisasi represents the personalisasi model
type Personalisasi struct {
	ID                int64     `json:"id"`
	UserID            int64     `json:"user_id"`
	Diabetes          bool      `json:"diabetes"`
	GERD              bool      `json:"gerd"`
	AsamUrat          bool      `json:"asam_urat"`
	Kolestrol         bool      `json:"kolestrol"`
	RendahKarbohidrat bool      `json:"rendah_karbohidrat"`
	TinggiProtein     bool      `json:"tinggi_protein"`
	Vegetarian        bool      `json:"vegetarian"`
	RendahGula        bool      `json:"rendah_gula"`
	RendahKalori      bool      `json:"rendah_kalori"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
