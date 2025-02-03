package queries

import (
	"context"
	"database/sql"
	"santapan/domain"

	"github.com/sirupsen/logrus"
)

type PostgresPersonalisasiQueryRepository struct {
	Conn *sql.DB
}

func NewPostgresPersonalisasiQueryRepository(Conn *sql.DB) *PostgresPersonalisasiQueryRepository {
	return &PostgresPersonalisasiQueryRepository{Conn}
}

// Fungsi fetch yang menerima query dan args sebagai parameter
func (m *PostgresPersonalisasiQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Personalisasi, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	// Hasil query akan disimpan dalam slice
	result := make([]domain.Personalisasi, 0)
	for rows.Next() {
		var p domain.Personalisasi

		// Ambil data dari baris hasil query
		err = rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Diabetes,
			&p.GERD,
			&p.AsamUrat,
			&p.Kolestrol,
			&p.RendahKarbohidrat,
			&p.TinggiProtein,
			&p.Vegetarian,
			&p.RendahGula,
			&p.RendahKalori,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			logrus.Error("Failed to scan row into struct: ", err)
			return nil, err
		}

		// Menambahkan hasil ke dalam slice
		result = append(result, p)
	}

	return result, nil
}

// Fungsi GetByUserID untuk mendapatkan data personalisasi berdasarkan user ID
func (m *PostgresPersonalisasiQueryRepository) GetByUserID(ctx context.Context, userID int64) (domain.Personalisasi, error) {
	// Query to fetch the personalization data based on userID
	query := `SELECT id, user_id, diabetes, gerd, asam_urat, kolestrol, rendah_karbohidrat, tinggi_protein, vegetarian, rendah_gula, rendah_kalori, created_at, updated_at
			  FROM user_condition WHERE user_id = $1`

	// Fetch the result using the fetch function
	result := make([]domain.Personalisasi, 0)
	result, err := m.fetch(ctx, query, userID)
	if err != nil {
		// If error occurs, return it
		return domain.Personalisasi{}, err
	}

	// If the result is found, return the first item
	if len(result) > 0 {
		return result[0], nil
	}

	// If no result is found, return an error
	return domain.Personalisasi{}, nil // No data found for the given userID
}
