package data

import (
	"context"
	"database/sql"
	"errors"
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

func (m StudentModel) Insert(student *Student) error {
	query := `
        INSERT INTO students (
		student_name,
		stage,
		student_id,
		state
		) 
        VALUES ($1, $2, $3, $4)
        RETURNING created_at, seq_in_college`
	args := []interface{}{student.StudentName,
		student.Stage,
		student.StudentId,
		student.State,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&student.CreatedAt, &student.SeqInCollege)
}

func (m StudentModel) GetAll() ([]*Student, error) {
	query := `SELECT * FROM students`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
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

func (m StudentModel) Get(id int64) (*Student, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	query := `
	SELECT seq_in_college, student_name, stage, student_id, state, created_at
	FROM students
	WHERE student_id = $1`
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

func (m StudentModel) Update(student *Student) error {
	query := `
	UPDATE students
	SET student_name = $1, stage = $2, state = $4
	WHERE student_id = $3`
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

func (m StudentModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM students
	WHERE student_id = $1`
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
