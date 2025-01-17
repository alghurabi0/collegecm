package data

import (
	"context"
	"database/sql"
	"errors"
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

func (m CarryoverModel) Insert(carryover *Carryover) error {
	query := `
        INSERT INTO carryovers (
		student_id,
		subject_id
		) 
        VALUES ($1, $2,)
        RETURNING id, created_at`
	args := []interface{}{
		carryover.StudentId,
		carryover.SubjectId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&carryover.Id, &carryover.CreatedAt)
}

// dd
func (m CarryoverModel) GetAll() ([]*Carryover, error) {
	query := `
	SELECT c.id, s.student_name AS student_name, sub.subject_name AS subject_name
	FROM carryovers c
	JOIN students s ON c.student_id = s.student_id
	JOIN subjects sub ON c.subject_id = sub.subject_id;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
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

func (m CarryoverModel) Get(id int64) (*Carryover, error) {
	query := `
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM carryovers c
	JOIN students s ON c.student_id = s.id
	JOIN subjects sub ON c.subject.id = sub.id
	WHERE c.id = $1;
	`
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

func (m CarryoverModel) GetSubjects(student_id int64) ([]*Carryover, error) {
	query := `
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM carryovers c
	JOIN students s ON c.student_id = s.id
	JOIN subjects sub ON c.subject.id = sub.id
	c.student_id = $1;
	`
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

func (m CarryoverModel) GetStudents(subject_id int64) ([]*Carryover, error) {
	query := `
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM carryovers c
	JOIN students s ON c.student_id = s.id
	JOIN subjects sub ON c.subject.id = sub.id
	c.subject_id = $1;
	`
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

func (m CarryoverModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM carryovers
	WHERE id = $1`
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
