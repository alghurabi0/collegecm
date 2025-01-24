package data

import (
	"context"
	"database/sql"
	"time"
)

type Year struct {
	Year string `json:"year"`
}

type YearModel struct {
	DB *sql.DB
}

func (y YearModel) GetAll() ([]*Year, error) {
	q := `SELECT * FROM years;`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := y.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var years []*Year
	for rows.Next() {
		var year Year
		err := rows.Scan(
			&year.Year,
		)
		if err != nil {
			return nil, err
		}
		years = append(years, &year)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return years, nil
}
