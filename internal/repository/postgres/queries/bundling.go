package queries

import (
	"context"
	"database/sql"
	"santapan/domain"
	"santapan/internal/repository"

	"github.com/sirupsen/logrus"
)

type BundlingRepository struct {
	Conn *sql.DB
}

// NewBundlingRepository creates an instance of BundlingRepository.
func NewBundlingRepository(conn *sql.DB) *BundlingRepository {
	return &BundlingRepository{conn}
}

func (m *BundlingRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Bundling, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			return
		}
	}()

	result = make([]domain.Bundling, 0)
	for rows.Next() {
		t := domain.Bundling{}
		err = rows.Scan(
			&t.ID,
			&t.BundlingName,
			&t.BundlingType,
			&t.Price,
			&t.ImageURL,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return
		}
		result = append(result, t)
	}

	return result, nil

}

// Fetch retrieves bundlings with ID-based pagination.
func (m *BundlingRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Bundling, nextCursor string, err error) {
	query := `SELECT id, bundling_name, bundling_type, price, image_url, created_at, updated_at
			  FROM bundling WHERE id > $1 ORDER BY id LIMIT $2`

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil {
		return
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return
	}

	if len(res) == 0 {
		return
	}

	nextCursor, err = repository.EncodeCursor(res[len(res)-1].ID)
	if err != nil {
		return
	}
	return
}

// fetchBundlingMenu retrieves bundling_menu with ID-based pagination.
func (m *BundlingRepository) fetchBundlingMenu(ctx context.Context, query string, args ...interface{}) (result []domain.BundlingMenu, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			return
		}
	}()

	result = make([]domain.BundlingMenu, 0)
	for rows.Next() {
		t := domain.BundlingMenu{}
		bundlingID := int64(0)
		menuID := int64(0)
		err = rows.Scan(
			&t.ID,
			&bundlingID,
			&menuID,
			&t.DayNumber,
			&t.MealDescription,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return
		}

		t.Bundling = domain.Bundling{ID: bundlingID}
		t.Menu = domain.Menu{ID: menuID}
		result = append(result, t)
	}

	return result, nil
}

// GetByID retrieves bundling by the given ID.
func (m *BundlingRepository) GetByID(ctx context.Context, id int64) (res domain.Bundling, err error) {
	query := `SELECT id, bundling_name, bundling_type, price, image_url, created_at, updated_at
		  	  FROM bundling WHERE id=$1`

	list, err := m.fetch(ctx, query, id)

	logrus.Info(err)
	if err != nil {
		return
	}

	if len(list) == 0 {
		return domain.Bundling{}, domain.ErrNotFound
	}

	res = list[0]
	return
}

// Get Bundling Menu By Bundling ID
func (m *BundlingRepository) FetchBundlingMenuByBundlingID(ctx context.Context, bundlingID int64) (res []domain.BundlingMenu, err error) {
	query := `SELECT id, bundling_id, menu_id, day_number, meal_description, created_at, updated_at
			  FROM bundling_menu WHERE bundling_id=$1 ORDER BY day_number`

	res, err = m.fetchBundlingMenu(ctx, query, bundlingID)
	if err != nil {
		return
	}

	return
}
