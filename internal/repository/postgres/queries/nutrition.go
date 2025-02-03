package queries

import (
	"context"
	"database/sql"
	"santapan/domain"

	"github.com/sirupsen/logrus"
)

type NutritionRepository struct {
	Conn *sql.DB
}

// NewNutritionRepository creates an instance of NutritionRepository.
func NewNutritionRepository(conn *sql.DB) *NutritionRepository {
	return &NutritionRepository{conn}
}

// fetch retrieves nutrition records based on a query and arguments.
func (n *NutritionRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Nutrition, err error) {
	rows, err := n.Conn.QueryContext(ctx, query, args...)
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

	result = make([]domain.Nutrition, 0)
	for rows.Next() {
		var t domain.Nutrition

		err = rows.Scan(
			&t.ID,
			&t.FoodName,
			&t.Calories,
			&t.Protein,
			&t.Fat,
			&t.Carbohydrates,
			&t.Sugar,
		)
		if err != nil {
			logrus.Error("Failed to scan row into struct: ", err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

// GetByClassification retrieves a single nutrition record filtered by food_name.
func (n *NutritionRepository) GetByClassification(ctx context.Context, foodName string) (domain.Nutrition, error) {
	query := `
    SELECT id, food_name, calories, protein, fat, carbohydrates, sugar
    FROM food_nutrition
    WHERE food_name = $1
    ORDER BY id
    LIMIT 1
	`

	result, err := n.fetch(ctx, query, foodName)
	if err != nil {
		return domain.Nutrition{}, err
	}

	if len(result) == 0 {
		return domain.Nutrition{}, domain.ErrNotFound
	}

	return result[0], nil
}
