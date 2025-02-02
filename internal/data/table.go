package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Table struct {
	ID        int64  `json:"id"`
	TableName string `json:"table_name"`
}

type TableModel struct {
	DB *sql.DB
}

func (t TableModel) GetByName(name, year string) (*Table, error) {
	tableName := name + "_" + year
	query := `SELECT id, table_name FROM tables WHERE table_name = $1`
	var table Table
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := t.DB.QueryRowContext(ctx, query, tableName).Scan(&table.ID, &table.TableName)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &table, nil
}
