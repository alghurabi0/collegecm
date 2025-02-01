package data

import (
	"context"
	"database/sql"
	"time"

	"collegecm.hamid.net/internal/validator"
)

type Privilege struct {
	UserId    int            `json:"user_id"`
	Year      string         `json:"year"`
	TableId   int            `json:"table_id"`
	Stage     sql.NullString `json:"stage"`
	SubjectId sql.NullInt64  `json:"subject_id"`
	CanRead   bool           `json:"can_read"`
	CanWrite  bool           `json:"can_write"`
	CreatedAt time.Time      `json:"created_at"`
}

func ValidatePrivilege(v *validator.Validator, privilege *Privilege) {
	// TODO - handle strings length with varchar
	v.Check(privilege.UserId != 0, "المستخدم", "يجب تزويد المعلومات")
	v.Check(privilege.TableId != 0, "الجدول", "يجب تزويد المعلومات")
}

type PrivilegeModel struct {
	DB *sql.DB
}

func (p PrivilegeModel) Insert(privilege *Privilege) error {
	query := `
        INSERT INTO privileges (user_id, table_id, stage, subject_id, can_read, can_write) 
        VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, table_id, stage, subject_id) DO UPDATE
        SET can_read = $5, can_write = $6
		RETURNING created_at`
	args := []interface{}{
		privilege.UserId,
		privilege.TableId,
		privilege.Stage,
		privilege.SubjectId,
		privilege.CanRead,
		privilege.CanWrite,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return p.DB.QueryRowContext(ctx, query, args...).Scan(&privilege.CreatedAt)
}

func (p PrivilegeModel) GetAll(userId int) ([]*Privilege, error) {
	query := `SELECT * FROM privileges WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var privileges []*Privilege
	for rows.Next() {
		var privilege Privilege
		err := rows.Scan(
			&privilege.UserId,
			&privilege.TableId,
			&privilege.Stage,
			&privilege.SubjectId,
			&privilege.CanRead,
			&privilege.CanWrite,
			&privilege.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		privileges = append(privileges, &privilege)
	}
	return privileges, nil
}

func (p PrivilegeModel) Delete(userId int, tableId int, stage sql.NullString, subjectId sql.NullInt64) error {
	query := `DELETE FROM privileges WHERE user_id = $1 AND table_id = $2 AND stage = $3 AND subject_id = $4`
	args := []interface{}{userId, tableId, stage, subjectId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := p.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
