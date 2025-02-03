package commands

import (
	"context"
	"database/sql"
	"santapan/domain"
)

type PostgresPersonalisasiCommandRepository struct {
	Conn *sql.DB
}

func NewPostgresPersonalisasiCommandRepository(Conn *sql.DB) *PostgresPersonalisasiCommandRepository {
	return &PostgresPersonalisasiCommandRepository{Conn}
}

// Insert or Update method for personalisasi data
func (m *PostgresPersonalisasiCommandRepository) InsertOrUpdate(ctx context.Context, p *domain.Personalisasi) (domain.Personalisasi, error) {
	query := `
		INSERT INTO user_condition (user_id, diabetes, gerd, asam_urat, kolestrol, rendah_karbohidrat, tinggi_protein, vegetarian, rendah_gula, rendah_kalori, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		ON CONFLICT (user_id) 
		DO UPDATE SET
			diabetes = EXCLUDED.diabetes,
			gerd = EXCLUDED.gerd,
			asam_urat = EXCLUDED.asam_urat,
			kolestrol = EXCLUDED.kolestrol,
			rendah_karbohidrat = EXCLUDED.rendah_karbohidrat,
			tinggi_protein = EXCLUDED.tinggi_protein,
			vegetarian = EXCLUDED.vegetarian,
			rendah_gula = EXCLUDED.rendah_gula,
			rendah_kalori = EXCLUDED.rendah_kalori,
			updated_at = NOW()
		RETURNING id`

	err := m.Conn.QueryRowContext(ctx, query, p.UserID, p.Diabetes, p.GERD, p.AsamUrat, p.Kolestrol, p.RendahKarbohidrat, p.TinggiProtein, p.Vegetarian, p.RendahGula, p.RendahKalori).Scan(&p.ID)
	if err != nil {
		return *p, err
	}
	return *p, nil
}
