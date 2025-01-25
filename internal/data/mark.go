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

type Mark struct {
	Id              int64     `json:"id"`
	StudentId       int64     `json:"student_id"`
	SubjectId       int64     `json:"subject_id"`
	StudentName     string    `json:"student_name"`
	SubjectName     string    `json:"subject_name"`
	SemesterMark    int       `json:"semester_mark"`
	MaxSemesterMark int       `json:"max_semester_mark"`
	FinalMark       int       `json:"final_mark"`
	MaxFinalExam    int       `json:"max_final_exam"`
	CreatedAt       time.Time `json:"-"`
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

func (m MarkModel) Insert(year string, mark *Mark) error {
	if strings.TrimSpace(year) == "" {
		return errors.New("invalid year")
	}
	tableName := fmt.Sprintf("marks_%s", year)
	query := fmt.Sprintf(`
        INSERT INTO %s (
		student_id,
		subject_id,
		semester_mark,
		final_mark
		) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at`, tableName)
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

func (m MarkModel) GetAll(year, stage string) ([]*Mark, error) {
	if strings.TrimSpace(year) == "" {
		return nil, errors.New("invalid year")
	}
	marksTable := fmt.Sprintf("marks_%s", year)
	studentsTable := fmt.Sprintf("students_%s", year)
	subjectsTable := fmt.Sprintf("subjects_%s", year)
	query := fmt.Sprintf(`
	SELECT
	c.id,
	s.student_name AS student_name,
	sub.subject_name AS subject_name,
	c.semester_mark,
	sub.max_semester_mark AS max_semester_mark,
	c.final_mark,
	sub.max_final_exam AS max_final_exam
	FROM %s c
	JOIN %s s ON c.student_id = s.student_id
	JOIN %s sub ON c.subject_id = sub.subject_id
	`, marksTable, studentsTable, subjectsTable)
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

	var marks []*Mark
	for rows.Next() {
		var mark Mark
		err := rows.Scan(
			&mark.Id,
			&mark.StudentName,
			&mark.SubjectName,
			&mark.SemesterMark,
			&mark.MaxSemesterMark,
			&mark.FinalMark,
			&mark.MaxFinalExam,
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

func (m MarkModel) Get(year string, id int64) (*Mark, error) {
	if strings.TrimSpace(year) == "" {
		return nil, errors.New("invalid year")
	}
	marksTable := fmt.Sprintf("marks_%s", year)
	studentsTable := fmt.Sprintf("students_%s", year)
	subjectsTable := fmt.Sprintf("subjects_%s", year)
	query := fmt.Sprintf(`
	SELECT
	c.id,
	s.student_name AS student_name,
	sub.subject_name AS subject_name,
	c.semester_mark,
	sub.max_semester_mark AS max_semester_mark,
	c.final_mark,
	sub.max_final_exam AS max_final_exam
	FROM %s c
	JOIN %s s ON c.student_id = s.student_id
	JOIN %s sub ON c.subject_id = sub.subject_id
	WHERE c.id = $1;`, marksTable, studentsTable, subjectsTable)
	var mark Mark
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&mark.Id,
		&mark.StudentName,
		&mark.SubjectName,
		&mark.SemesterMark,
		&mark.MaxSemesterMark,
		&mark.FinalMark,
		&mark.MaxFinalExam,
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

func (m MarkModel) GetRaw(year string, id int64) (*Mark, error) {
	if strings.TrimSpace(year) == "" {
		return nil, errors.New("invalid year")
	}
	marksTable := fmt.Sprintf("marks_%s", year)
	query := fmt.Sprintf(`SELECT id, student_id, subject_id, semester_mark, final_mark from %s WHERE id = $1;`, marksTable)
	var mark Mark
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&mark.Id,
		&mark.StudentId,
		&mark.SubjectId,
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

func (m MarkModel) Update(year string, mark *Mark) error {
	if strings.TrimSpace(year) == "" {
		return errors.New("invalid year")
	}
	marksTable := fmt.Sprintf("marks_%s", year)
	query := fmt.Sprintf(`
	UPDATE %s
	SET student_id = $2, subject_id = $3, semester_mark = $4, final_mark = $5
	WHERE id = $1`, marksTable)
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

func (m MarkModel) Delete(year string, id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}
	if strings.TrimSpace(year) == "" {
		return errors.New("invalid year")
	}
	tableName := fmt.Sprintf("marks_%s", year)
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
