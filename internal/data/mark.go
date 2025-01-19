package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"collegecm.hamid.net/internal/validator"
)

type Mark struct {
	Id           int64     `json:"id"`
	StudentId    int64     `json:"student_id"`
	SubjectId    int64     `json:"subject_id"`
	StudentName  string    `json:"student_name"`
	SubjectName  string    `json:"subject_name"`
	SemesterMark int       `json:"semester_mark"`
	FinalMark    int       `json:"final_mark"`
	CreatedAt    time.Time `json:"-"`
}

func ValidateMark(v *validator.Validator, mark *Mark, sem, fin int) {
	// TODO - handle strings length with varchar
	v.Check(mark.StudentId >= 0, "رقم الطالب", "يجب ان يكون 0 او اكبر")
	v.Check(mark.SubjectId >= 0, "رقم المادة", "يجب ان يكون 0 او اكبر")
	v.Check(mark.SemesterMark >= 0, "السعي", "يجب ان يكون 0 او اكبر")
	v.Check(mark.FinalMark >= 0, "درجة الامتحان النهائي", "يجب ان يكون 0 او اكبر")
	v.Check(mark.SemesterMark <= sem, "السعي", "يجب ان يساوي او اقل من درجة السعي القصوى")
	v.Check(mark.FinalMark <= fin, "درجة الامتحان النهائي", "يجب ان يساوي او اقل من درجة الامتحان القصوى")
}

type MarkModel struct {
	DB *sql.DB
}

func (m MarkModel) Insert(mark *Mark) error {
	query := `
        INSERT INTO marks (
		student_id,
		subject_id,
		semester_mark,
		final_mark
		) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at`
	args := []interface{}{
		mark.StudentId,
		mark.SubjectId,
		mark.SemesterMark,
		mark.FinalMark,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&mark.Id, &mark.CreatedAt)
}

func (m MarkModel) GetAll() ([]*Mark, error) {
	query := `
	SELECT c.id, s.student_name AS student_name, sub.subject_name AS subject_name, c.semester_mark, c.final_mark
	FROM marks c
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

	var marks []*Mark
	for rows.Next() {
		var mark Mark
		err := rows.Scan(
			&mark.Id,
			&mark.StudentName,
			&mark.SubjectName,
			&mark.SemesterMark,
			&mark.FinalMark,
		)
		if err != nil {
			return nil, err
		}
		marks = append(marks, &mark)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return marks, nil
}

func (m MarkModel) Get(id int64) (*Mark, error) {
	query := `
	SELECT c.id, s.student_name AS student_name, sub.subject_name AS subject_name, c.semester_mark, c.final_mark
	FROM marks c
	JOIN students s ON c.student_id = s.student_id
	JOIN subjects sub ON c.subject_id = sub.subject_id
	WHERE c.id = $1;
	`
	var mark Mark
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&mark.Id,
		&mark.StudentName,
		&mark.SubjectName,
		&mark.SemesterMark,
		&mark.FinalMark,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &mark, nil
}

func (m MarkModel) Update(mark *Mark) error {
	query := `
	UPDATE marks
	SET student_id = $2, subject_id = $3 semester_mark = $4, final_mark = $5
	WHERE id = $1`
	args := []interface{}{
		&mark.Id,
		&mark.StudentId,
		&mark.SubjectId,
		&mark.SemesterMark,
		&mark.FinalMark,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m MarkModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM marks
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
