package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"collegecm.hamid.net/internal/validator"
)

type Student struct {
	SeqInCollege int       `json:"seq_in_college" csv:"-"`
	StudentName  string    `json:"student_name" csv:"student_name"`
	Stage        string    `json:"stage" csv:"stage"`
	StudentId    int       `json:"student_id" csv:"student_id"`
	State        string    `json:"state" csv:"state"`
	CreatedAt    time.Time `json:"-" csv:"-"`
	Year         string    `json:"-"`
}

func ValidateStudent(v *validator.Validator, student *Student) {
	// TODO - handle strings length with varchar
	v.Check(student.StudentName != "", "اسم الطالب", "يجب تزويد المعلومات")
	v.Check(student.Stage != "", "المرحلة", "يجب تزويد المعلومات")
	v.Check(student.StudentId >= 0, "رقم الطالب", "يجب تزويد المعلومات")
	v.Check(student.State != "", "الوضع", "يجب تزويد المعلومات")
}

type StudentModel struct {
	DB *sql.DB
}

func (m StudentModel) Insert(year string, student *Student) error {
	if strings.TrimSpace(year) == "" {
		return errors.New("invalid year")
	}
	tableName := fmt.Sprintf("students_%s", year)
	query := fmt.Sprintf(`
        INSERT INTO %s (
		student_name,
		stage,
		student_id,
		state
		) 
        VALUES ($1, $2, $3, $4)
        RETURNING created_at, seq_in_college`, tableName)
	args := []interface{}{student.StudentName,
		student.Stage,
		student.StudentId,
		student.State,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&student.CreatedAt, &student.SeqInCollege)
}

func (m StudentModel) GetAll(year, stage string) ([]*Student, error) {
	if strings.TrimSpace(year) == "" {
		return nil, errors.New("invalid year")
	}
	tableName := fmt.Sprintf("students_%s", year)
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	var args []interface{}
	if stage != "all" {
		query += " WHERE stage = $1"
		args = append(args, stage)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*Student
	for rows.Next() {
		var student Student
		err := rows.Scan(
			&student.SeqInCollege,
			&student.StudentName,
			&student.Stage,
			&student.StudentId,
			&student.State,
			&student.CreatedAt,
			&student.Year,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, &student)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

func (m StudentModel) Get(year string, id int64) (*Student, error) {
	// Define the SQL query for retrieving the movie data.
	tableName := fmt.Sprintf("students_%s", year)
	query := fmt.Sprintf(`
	SELECT seq_in_college, student_name, stage, student_id, state, created_at
	FROM %s
	WHERE student_id = $1`, tableName)
	var student Student
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&student.SeqInCollege,
		&student.StudentName,
		&student.Stage,
		&student.StudentId,
		&student.State,
		&student.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &student, nil
}

func (m StudentModel) Update(year string, student *Student) error {
	if student.StudentId < 0 {
		return ErrRecordNotFound
	}
	if strings.TrimSpace(year) == "" {
		return errors.New("invalid year")
	}
	tableName := fmt.Sprintf("students_%s", year)
	query := fmt.Sprintf(`
	UPDATE %s
	SET student_name = $1, stage = $2, state = $4
	WHERE student_id = $3`, tableName)
	args := []interface{}{
		&student.StudentName,
		&student.Stage,
		&student.StudentId,
		&student.State,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m StudentModel) Delete(year string, id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}
	if strings.TrimSpace(year) == "" {
		return errors.New("invalid year")
	}
	tableName := fmt.Sprintf("students_%s", year)
	query := fmt.Sprintf(`
	DELETE FROM %s
	WHERE student_id = $1`, tableName)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)
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

func (m StudentModel) GetCustom(tableName string, id int64) (*Student, error) {
	query := fmt.Sprintf(`
	SELECT student_id, student_name, stage
	FROM %s
	WHERE student_id = $1`, tableName)
	var student Student
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&student.StudentId,
		&student.StudentName,
		&student.Stage,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &student, nil
}

func (m StudentModel) GetStage(id int64, year string) (string, error) {
	studentsTable := fmt.Sprintf("students_%s", year)
	query := fmt.Sprintf(`
	SELECT stage
	FROM %s
	WHERE student_id = $1`, studentsTable)
	var stage string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&stage)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err
		}
	}
	return stage, nil
}
