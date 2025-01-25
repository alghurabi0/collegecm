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

type Carryover struct {
	Id          int64     `json:"id"`
	StudentId   int64     `json:"student_id"`
	SubjectId   int64     `json:"subject_id"`
	StudentName string    `json:"student_name"`
	SubjectName string    `json:"subject_name"`
	CreatedAt   time.Time `json:"-"`
}

func ValidateCarryover(v *validator.Validator, carryover *Carryover) {
	// TODO - handle strings length with varchar
	v.Check(carryover.StudentId >= 0, "رقم الطالب", "يجب ان يكون 0 او اكبر")
	v.Check(carryover.SubjectId >= 0, "رقم المادة", "يجب ان يكون 0 او اكبر")
}

type CarryoverModel struct {
	DB *sql.DB
}

func (m CarryoverModel) Insert(year string, carryover *Carryover) error {
	if strings.TrimSpace(year) == "" {
		return errors.New("invalid year")
	}
	tableName := fmt.Sprintf("carryovers_%s", year)
	query := fmt.Sprintf(`
        INSERT INTO %s (
		student_id,
		subject_id
		) 
        VALUES ($1, $2)
        RETURNING id, created_at`, tableName)
	args := []interface{}{
		carryover.StudentId,
		carryover.SubjectId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&carryover.Id, &carryover.CreatedAt)
}

// ddd
func (m CarryoverModel) GetAll(year, stage string) ([]*Carryover, error) {
	if strings.TrimSpace(year) == "" {
		return nil, errors.New("invalid year")
	}
	carryoversTable := fmt.Sprintf("carryovers_%s", year)
	studentsTable := fmt.Sprintf("students_%s", year)
	subjectsTable := fmt.Sprintf("subjects_%s", year)
	query := fmt.Sprintf(`
		SELECT c.id, s.student_name AS student_name, sub.subject_name AS subject_name
		FROM %s c
		JOIN %s s ON c.student_id = s.student_id
		JOIN %s sub ON c.subject_id = sub.subject_id
	`, carryoversTable, studentsTable, subjectsTable)
	var args []interface{}
	if stage != "all" {
		query += " WHERE s.stage = $1"
		args = append(args, stage)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carryovers []*Carryover
	for rows.Next() {
		var carryover Carryover
		err := rows.Scan(
			&carryover.Id,
			&carryover.StudentName,
			&carryover.SubjectName,
		)
		if err != nil {
			return nil, err
		}
		carryovers = append(carryovers, &carryover)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return carryovers, nil
}

func (m CarryoverModel) Get(year string, id int64) (*Carryover, error) {
	if id < 0 {
		return nil, ErrRecordNotFound
	}
	if strings.TrimSpace(year) == "" {
		return nil, errors.New("invalid year")
	}
	carryoversTable := fmt.Sprintf("carryovers_%s", year)
	studentsTable := fmt.Sprintf("students_%s", year)
	subjectsTable := fmt.Sprintf("subjects_%s", year)
	query := fmt.Sprintf(`
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM %s c
	JOIN %s s ON c.student_id = s.id
	JOIN %s sub ON c.subject.id = sub.id
	WHERE c.id = $1;`, carryoversTable, studentsTable, subjectsTable)
	var carryover Carryover
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&carryover.Id,
		&carryover.StudentName,
		&carryover.SubjectName,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &carryover, nil
}

func (m CarryoverModel) Find(student_id, subject_id int64) (*Carryover, error) {
	query := `
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM carryovers c
	JOIN students s ON c.student_id = s.id
	JOIN subjects sub ON c.subject.id = sub.id
	c.student_id = $1 AND c.subject_id = $2;
	`
	var carryover Carryover
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, student_id, subject_id).Scan(
		&carryover.Id,
		&carryover.StudentName,
		&carryover.SubjectName,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &carryover, nil
}

func (m CarryoverModel) GetSubjects(year string, student_id int64) ([]*Carryover, error) {
	if student_id < 0 {
		return nil, ErrRecordNotFound
	}
	if strings.TrimSpace(year) == "" {
		return nil, errors.New("invalid year")
	}
	carryoversTable := fmt.Sprintf("carryovers_%s", year)
	studentsTable := fmt.Sprintf("students_%s", year)
	subjectsTable := fmt.Sprintf("subjects_%s", year)
	query := fmt.Sprintf(`
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM %s c
	JOIN %s s ON c.student_id = s.id
	JOIN %s sub ON c.subject.id = sub.id
	c.student_id = $1;`, carryoversTable, studentsTable, subjectsTable)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, student_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carryovers []*Carryover
	for rows.Next() {
		var carryover Carryover
		err := rows.Scan(
			&carryover.Id,
			&carryover.StudentName,
			&carryover.SubjectName,
		)
		if err != nil {
			return nil, err
		}
		carryovers = append(carryovers, &carryover)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return carryovers, nil
}

func (m CarryoverModel) GetStudents(year string, subject_id int64) ([]*Carryover, error) {
	if subject_id < 0 {
		return nil, ErrRecordNotFound
	}
	if strings.TrimSpace(year) == "" {
		return nil, errors.New("invalid year")
	}
	carryoversTable := fmt.Sprintf("carryovers_%s", year)
	studentsTable := fmt.Sprintf("students_%s", year)
	subjectsTable := fmt.Sprintf("subjects_%s", year)
	query := fmt.Sprintf(`
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM %s c
	JOIN %s s ON c.student_id = s.id
	JOIN %s sub ON c.subject.id = sub.id
	c.subject_id = $1;`, carryoversTable, studentsTable, subjectsTable)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, subject_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carryovers []*Carryover
	for rows.Next() {
		var carryover Carryover
		err := rows.Scan(
			&carryover.Id,
			&carryover.StudentName,
			&carryover.SubjectName,
		)
		if err != nil {
			return nil, err
		}
		carryovers = append(carryovers, &carryover)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return carryovers, nil
}

func (m CarryoverModel) Delete(year string, id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}
	if strings.TrimSpace(year) == "" {
		return errors.New("invalid year")
	}
	tableName := fmt.Sprintf("carryovers_%s", year)
	query := fmt.Sprintf(`
	DELETE FROM %s
	WHERE id = $1`, tableName)
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
