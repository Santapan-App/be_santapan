package queries

import (
	"context"
	"database/sql"
	"santapan/domain"
	"santapan/internal/repository"

	"github.com/sirupsen/logrus"
)

type MenuRepository struct {
	Conn *sql.DB
}

// NewBannerRepository creates an instance of BannerRepository.
func NewMenuRepository(conn *sql.DB) *MenuRepository {
	return &MenuRepository{conn}
}

func (m *MenuRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Menu, err error) {
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

	result = make([]domain.Menu, 0)
	for rows.Next() {
		var t domain.Menu
		var nutritionJSON, featuresJSON []byte

		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Price,
			&t.ImageURL,
			&t.Nutrition,
			&t.Features,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			logrus.Error("Failed to scan row into struct: ", err)
			return nil, err
		}

		logrus.Print(nutritionJSON)
		logrus.Print(featuresJSON)
		// Decode Nutrition JSON
		// if len(nutritionJSON) > 0 {
		// 	if err := json.Unmarshal(nutritionJSON, &t.Nutrition); err != nil {
		// 		logrus.Error("Failed to decode nutrition JSON: ", err)
		// 		return nil, err
		// 	}
		// }

		// Decode Features JSON
		// if len(featuresJSON) > 0 {
		// 	if err := json.Unmarshal(featuresJSON, &t.Features); err != nil {
		// 		logrus.Error("Failed to decode features JSON: ", err)
		// 		return nil, err
		// 	}
		// }

		result = append(result, t)
	}

	return result, nil
}

// Fetch retrieves banners with ID-based pagination.
func (m *MenuRepository) Fetch(ctx context.Context, cursor string, num int64, search string) (res []domain.Menu, nextCursor string, err error) {
	query := `
    SELECT id, title, description, price, image_url, nutrition, features, created_at, updated_at
    FROM menu
    WHERE id > $1 AND title ILIKE '%' || $3 || '%'
    ORDER BY id
    LIMIT $2
`
	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num, search)
	if err != nil {
		return nil, "", err
	}

	// Set the nextCursor if the result count reaches the limit
	if len(res) == int(num) {
		nextCursor, err = repository.EncodeCursor(res[len(res)-1].ID)
		if err != nil {
			logrus.Error("Failed to encode cursor: ", err)
			return res, "", err
		}
	}

	return res, nextCursor, nil
}

func (m *MenuRepository) GetByID(ctx context.Context, id int64) (domain.Menu, error) {
	query := `SELECT id, title, description, price, image_url, nutrition, features, created_at, updated_at
			  FROM menu WHERE id = $1`

	result := make([]domain.Menu, 0)
	result, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Menu{}, err
	}

	if len(result) > 0 {
		return result[0], nil
	}

	return domain.Menu{}, domain.ErrNotFound
}

func (m *MenuRepository) GetByCategoryID(ctx context.Context, categoryID int64, search string) ([]domain.Menu, error) {
	query := `SELECT m.id, m.title, m.description, m.price, m.image_url, m.nutrition, m.features, m.created_at, m.updated_at
			  FROM menu m
			  JOIN category_menu mc ON m.id = mc.menu_id
			  WHERE mc.category_id = $1 AND m.title LIKE '%' || $2 || '%'`

	result := make([]domain.Menu, 0)
	result, err := m.fetch(ctx, query, categoryID, search)
	if err != nil {
		return nil, err
	}

	return result, nil
}
