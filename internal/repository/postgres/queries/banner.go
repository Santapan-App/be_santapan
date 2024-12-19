package queries

import (
	"context"
	"database/sql"
	"santapan/domain"
	"santapan/internal/repository"

	"github.com/sirupsen/logrus"
)

type BannerRepository struct {
	Conn *sql.DB
}

// NewBannerRepository creates an instance of BannerRepository.
func NewBannerRepository(conn *sql.DB) *BannerRepository {
	return &BannerRepository{conn}
}

func (m *BannerRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Banner, err error) {
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

	result = make([]domain.Banner, 0)
	for rows.Next() {
		t := domain.Banner{}
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.ImageURL,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

// Fetch retrieves banners with ID-based pagination.
func (m *BannerRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Banner, nextCursor string, err error) {
	query := `SELECT id, title, image_url, created_at, updated_at
			  FROM banner WHERE id > $1 ORDER BY id LIMIT $2`

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
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

// GetByID retrieves a banner by its ID.
func (m *BannerRepository) GetByID(ctx context.Context, id int64) (res domain.Banner, err error) {
	query := `SELECT id, title, image_url, created_at, updated_at
			  FROM banner WHERE id = $1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Banner{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}
