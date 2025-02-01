package data

import "database/sql"

type Table struct {
	ID        int64  `json:"id"`
	TableName string `json:"table_name"`
}

type TableModel struct {
	DB *sql.DB
}
