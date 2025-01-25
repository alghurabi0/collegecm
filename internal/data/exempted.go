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

type Exempted struct {
	Id          int64     `json:"id"`
	StudentId   int64     `json:"student_id"`
	SubjectId   int64     `json:"subject_id"`
	StudentName string    `json:"student_name"`
	SubjectName string    `json:"subject_name"`
	CreatedAt   time.Time `json:"-"`
}

func ValidateExempted(v *validator.Validator, exempted *Exempted) {
	// TODO - handle strings length with varchar
	v.Check(exempted.StudentId >= 0, "رقم الطالب", "يجب ان يكون 0 او اكبر")
	v.Check(exempted.SubjectId >= 0, "رقم المادة", "يجب ان يكون 0 او اكبر")
}

type ExemptedModel struct {
	DB *sql.DB
}

func (m ExemptedModel) Insert(exempted *Exempted) error {
	query := `
        INSERT INTO exempted (
		student_id,
		subject_id
		) 
        VALUES ($1, $2)
        RETURNING id, created_at`
	args := []interface{}{
		exempted.StudentId,
		exempted.SubjectId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&exempted.Id, &exempted.CreatedAt)
}

// ddd
func (m ExemptedModel) GetAll(year, stage string) ([]*Exempted, error) {
	if strings.TrimSpace(year) == "" {
		return nil, errors.New("invalid year")
	}
	exemptedTable := fmt.Sprintf("exempted_%s", year)
	studentsTable := fmt.Sprintf("students_%s", year)
	subjectsTable := fmt.Sprintf("subjects_%s", year)
	query := fmt.Sprintf(`
		SELECT c.id, s.student_name AS student_name, sub.subject_name AS subject_name
		FROM %s c
		JOIN %s s ON c.student_id = s.student_id
		JOIN %s sub ON c.subject_id = sub.subject_id
	`, exemptedTable, studentsTable, subjectsTable)
	var args []interface{}
	if stage != "all" {
		query += " WHERE s.stage = $1 AND sub.stage = $1"
		args = append(args, stage)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exempteds []*Exempted
	for rows.Next() {
		var exempted Exempted
		err := rows.Scan(
			&exempted.Id,
			&exempted.StudentName,
			&exempted.SubjectName,
		)
		if err != nil {
			return nil, err
		}
		exempteds = append(exempteds, &exempted)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return exempteds, nil
}

func (m ExemptedModel) Get(id int64) (*Exempted, error) {
	query := `
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM exempted c
	JOIN students s ON c.student_id = s.id
	JOIN subjects sub ON c.subject.id = sub.id
	WHERE c.id = $1;
	`
	var exempted Exempted
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&exempted.Id,
		&exempted.StudentName,
		&exempted.SubjectName,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &exempted, nil
}

func (m ExemptedModel) Find(student_id, subject_id int64) (*Exempted, error) {
	query := `
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM exempted c
	JOIN students s ON c.student_id = s.id
	JOIN subjects sub ON c.subject.id = sub.id
	c.student_id = $1 AND c.subject_id = $2;
	`
	var exempted Exempted
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, student_id, subject_id).Scan(
		&exempted.Id,
		&exempted.StudentName,
		&exempted.SubjectName,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &exempted, nil
}

func (m ExemptedModel) GetSubjects(student_id int64) ([]*Exempted, error) {
	query := `
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM exempted c
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

	var exempteds []*Exempted
	for rows.Next() {
		var exempted Exempted
		err := rows.Scan(
			&exempted.Id,
			&exempted.StudentName,
			&exempted.SubjectName,
		)
		if err != nil {
			return nil, err
		}
		exempteds = append(exempteds, &exempted)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return exempteds, nil
}

func (m ExemptedModel) GetStudents(subject_id int64) ([]*Exempted, error) {
	query := `
	SELECT c.id, s.name AS student_name, sub.name AS subject_name
	FROM exempted c
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

	var exempteds []*Exempted
	for rows.Next() {
		var exempted Exempted
		err := rows.Scan(
			&exempted.Id,
			&exempted.StudentName,
			&exempted.SubjectName,
		)
		if err != nil {
			return nil, err
		}
		exempteds = append(exempteds, &exempted)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return exempteds, nil
}

func (m ExemptedModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM exempted
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
