package data

import (
	"context"
	"database/sql"
	"time"

	"collegecm.hamid.net/internal/validator"
)

type Privilege struct {
	UserId    int       `json:"user_id"`
	Year      string    `json:"year"`
	TableId   int       `json:"table_id"`
	Stage     string    `json:"stage"`
	SubjectId int       `json:"subject_id"`
	CanRead   bool      `json:"can_read"`
	CanWrite  bool      `json:"can_write"`
	CreatedAt time.Time `json:"created_at"`
	TableName string    `json:"-"`
}

type CustomPrivilegeAccess struct {
	Students   bool
	Subjects   bool
	Carryovers bool
	Exempted   bool
	Marks      bool
}

func ValidatePrivilege(v *validator.Validator, privilege *Privilege) {
	v.Check(privilege.UserId > 0, "المستخدم", "يجب تزويد المعلومات")
	v.Check(privilege.Year != "", "الجدول", "يجب تزويد المعلومات")
	v.Check(privilege.TableId == -1 || privilege.TableId > 0, "الجدول", "يجب تزويد المعلومات")
	v.Check(privilege.Stage != "" && (privilege.Stage == "الاولى" || privilege.Stage == "الثانية" ||
		privilege.Stage == "الثالثة" || privilege.Stage == "الرابعة" ||
		privilege.Stage == "الخامسة" || privilege.Stage == "السادسة" ||
		privilege.Stage == "all"), "المرحلة", "يجب تزويد المعلومات")
	v.Check(privilege.SubjectId == -1 || privilege.SubjectId > 0, "المادة", "يجب تزويد المعلومات")
	v.Check(privilege.CanRead || !privilege.CanRead, "الصلاحيات", "يجب تزويد المعلومات")
	v.Check(privilege.CanWrite || !privilege.CanWrite, "الصلاحيات", "يجب تزويد المعلومات")

}

type PrivilegeModel struct {
	DB *sql.DB
}

func (p PrivilegeModel) Insert(privilege *Privilege) error {
	query := `
    INSERT INTO privileges (user_id, year, table_id, stage, subject_id, can_read, can_write)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    ON CONFLICT (user_id, year, table_id, stage, subject_id) DO UPDATE
    SET can_read = $6, can_write = $7
	RETURNING created_at
`
	args := []interface{}{
		privilege.UserId,
		privilege.Year,
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
			&privilege.Year,
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

func (p PrivilegeModel) CheckAccess(userId int, tableName, stage string) (*Privilege, error) {
	query := `
	SELECT p.user_id, t.table_name as table_name, p.stage, p.can_read, p.can_write
	FROM privileges p
	JOIN tables t ON p.table_id = t.id
	WHERE p.user_id = $1 AND t.table_name = $2 AND (p.stage = $3 OR p.stage = 'all')
	LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var privilege Privilege
	err := p.DB.QueryRowContext(ctx, query, userId, tableName, stage).Scan(
		&privilege.UserId,
		&privilege.TableName,
		&privilege.Stage,
		&privilege.CanRead,
		&privilege.CanWrite,
	)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &privilege, nil
}

func (p PrivilegeModel) CheckWriteAccess(userId int, tableName string) (bool, error) {
	query := `
	SELECT p.user_id, t.table_name as table_name, p.can_read, p.can_write
	FROM privileges p
	JOIN tables t ON p.table_id = t.id
	WHERE p.user_id = $1 AND t.table_name = $2
	LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var privilege Privilege
	err := p.DB.QueryRowContext(ctx, query, userId, tableName).Scan(
		&privilege.UserId,
		&privilege.TableName,
		&privilege.CanRead,
		&privilege.CanWrite,
	)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return false, nil
		default:
			return false, err
		}
	}
	if privilege.CanWrite {
		return true, nil
	}
	return false, nil
}

func (p PrivilegeModel) CheckCustomAccess(userId int, year, stage string) (*CustomPrivilegeAccess, error) {
	query := `
	SELECT p.can_read
	FROM privileges p
	JOIN tables t ON p.table_id = t.id
	WHERE p.user_id = $1 AND t.table_name = $2 AND (p.stage = $3 OR p.stage = 'all')
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var access CustomPrivilegeAccess
	table := "students_" + year
	err := p.DB.QueryRowContext(ctx, query, userId, table, stage).Scan(&access.Students)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			access.Students = false
		default:
			return nil, err
		}
	}
	table = "subjects_" + year
	err = p.DB.QueryRowContext(ctx, query, userId, table, stage).Scan(&access.Subjects)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			access.Subjects = false
		default:
			return nil, err
		}
	}
	table = "carryovers_" + year
	err = p.DB.QueryRowContext(ctx, query, userId, table, stage).Scan(&access.Carryovers)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			access.Carryovers = false
		default:
			return nil, err
		}
	}
	table = "exempted_" + year
	err = p.DB.QueryRowContext(ctx, query, userId, table, stage).Scan(&access.Exempted)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			access.Exempted = false
		default:
			return nil, err
		}
	}
	table = "marks_" + year
	err = p.DB.QueryRowContext(ctx, query, userId, table, stage).Scan(&access.Marks)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			access.Marks = false
		default:
			return nil, err
		}
	}
	return &access, nil
}
