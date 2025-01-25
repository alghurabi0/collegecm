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

type Subject struct {
	ID                 int       `json:"subject_id" csv:"subject_id"`
	SubjectName        string    `json:"subject_name" csv:"subject_name"`
	SubjectNameEnglish string    `json:"subject_name_english" csv:"subject_name_english"`
	Stage              string    `json:"stage" csv:"stage"`
	Semester           string    `json:"semester" csv:"semester"`
	Department         string    `json:"department" csv:"department"`
	MaxTheoryMark      int       `json:"max_theory_mark" csv:"max_theory_mark"`
	MaxLabMark         int       `json:"max_lab_mark" csv:"max_lab_mark"`
	MaxSemesterMark    int       `json:"max_semester_mark" csv:"max_semester_mark"`
	MaxFinalExam       int       `json:"max_final_exam" csv:"max_final_exam"`
	Credits            int       `json:"credits" csv:"credits"`
	Active             string    `json:"active" csv:"active"`
	Ministerial        string    `json:"ministerial" csv:"ministerial"`
	CreatedAt          time.Time `json:"-" csv:"-"`
}

func ValidateSubject(v *validator.Validator, subject *Subject) {
	// TODO - handle strings length with varchar
	v.Check(subject.SubjectName != "", "subject_name", "must be provided")
	v.Check(subject.SubjectNameEnglish != "", "subject_name_english", "must be provided")
	v.Check(subject.Stage != "", "stage", "must be provided")
	v.Check(subject.Semester != "", "semester", "must be provided")
	v.Check(subject.Department != "", "department", "must be provided")
	v.Check(subject.MaxTheoryMark >= 0, "max_theory_mark", "must not be less than zero")
	v.Check(subject.MaxLabMark >= 0, "max_lab_mark", "must not be zero less than zero")
	v.Check(subject.MaxSemesterMark >= 0, "max_semester_mark", "must not be zero less than zero")
	v.Check(subject.MaxFinalExam >= 0, "max_final_exam", "must not be zero less than zero")
	v.Check(subject.Credits >= 0, "credits", "must not be less than zero")
	v.Check(subject.Active == "لا" || subject.Active == "نعم", "active", "must equal to لا or نعم")
	v.Check(subject.Ministerial == "لا" || subject.Ministerial == "نعم", "ministerial", "must equal to لا or نعم")
}

type SubjectModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m SubjectModel) Insert(subject *Subject) error {
	query := `
        INSERT INTO subjects (
		subject_id,
		subject_name,
		subject_name_english,
		stage,
		semester,
		department,
		max_theory_mark,
		max_lab_mark,
		max_semester_mark,
		max_final_exam,
		credits,
		active,
		ministerial
		) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING created_at`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{subject.ID,
		subject.SubjectName,
		subject.SubjectNameEnglish,
		subject.Stage,
		subject.Semester,
		subject.Department,
		subject.MaxTheoryMark,
		subject.MaxLabMark,
		subject.MaxSemesterMark,
		subject.MaxFinalExam,
		subject.Credits,
		subject.Active,
		subject.Ministerial,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&subject.CreatedAt)
}

func (m SubjectModel) GetAll(year, stage string) ([]*Subject, error) {
	if strings.TrimSpace(year) == "" {
		return nil, errors.New("invalid year")
	}
	tableName := fmt.Sprintf("subjects_%s", year)
	fmt.Println("Querying table:", tableName)
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
	defer rows.Close() // Ensure rows are closed after function returns

	// Slice to hold the movies
	var subjects []*Subject

	// Iterate over rows
	for rows.Next() {
		var subject Subject
		err := rows.Scan(
			&subject.ID,
			&subject.SubjectName,
			&subject.SubjectNameEnglish,
			&subject.Stage,
			&subject.Semester,
			&subject.Department,
			&subject.MaxTheoryMark,
			&subject.MaxLabMark,
			&subject.MaxSemesterMark,
			&subject.MaxFinalExam,
			&subject.Credits,
			&subject.Active,
			&subject.Ministerial,
			&subject.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, &subject)
	}

	// Check for any iteration errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m SubjectModel) Get(id int64) (*Subject, error) {
	// The PostgreSQL bigserial type that we're using for the movie ID starts
	// auto-incrementing at 1 by default, so we know that no movies will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	query := `
	SELECT subject_id, subject_name, subject_name_english, stage, semester, department, max_theory_mark, max_lab_mark, max_semester_mark, max_final_exam, credits, active, ministerial
	FROM subjects
	WHERE subject_id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var subject Subject
	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&subject.ID,
		&subject.SubjectName,
		&subject.SubjectNameEnglish,
		&subject.Stage,
		&subject.Semester,
		&subject.Department,
		&subject.MaxTheoryMark,
		&subject.MaxLabMark,
		&subject.MaxSemesterMark,
		&subject.MaxFinalExam,
		&subject.Credits,
		&subject.Active,
		&subject.Ministerial,
	)
	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &subject, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m SubjectModel) Update(subject *Subject) error {
	query := `
	UPDATE subjects
	SET subject_id = $1, subject_name = $2, subject_name_english = $3, stage = $4, semester = $5, department = $6, max_theory_mark = $7, max_lab_mark = $8, max_semester_mark = $9, max_final_exam = $10, credits = $11, active = $12, ministerial = $13
	WHERE subject_id = $1`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		subject.ID,
		subject.SubjectName,
		subject.SubjectNameEnglish,
		subject.Stage,
		subject.Semester,
		subject.Department,
		subject.MaxTheoryMark,
		subject.MaxLabMark,
		subject.MaxSemesterMark,
		subject.MaxFinalExam,
		subject.Credits,
		subject.Active,
		subject.Ministerial,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m SubjectModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
	DELETE FROM subjects
	WHERE subject_id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
