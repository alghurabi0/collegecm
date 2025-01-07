package data

import (
	"database/sql"
	"time"

	"collegecm.hamid.net/internal/validator"
)

type Subject struct {
	ID                 int       `json:"subject_id"`
	SubjectName        string    `json:"subject_name"`
	SubjectNameEnglish string    `json:"subject_name_english"`
	Stage              string    `json:"stage"`
	Semester           string    `json:"semester"`
	Department         string    `json:"department"`
	MaxTheoryMark      int       `json:"max_theory_mark"`
	MaxLabMark         int       `json:"max_lab_mark"`
	MaxSemesterMark    int       `json:"max_semester_mark"`
	MaxFinalExam       int       `json:"max_final_exam"`
	Credits            int       `json:"credits"`
	Active             string    `json:"active"`
	Ministerial        string    `json:"ministerial"`
	CreatedAt          time.Time `json:"-"`
}

func ValidateSubject(v *validator.Validator, subject *Subject) {
	v.Check(subject.SubjectName != "", "subject_name", "must be provided")
	v.Check(len(subject.SubjectName) <= 50, "subject_name", "must not be more than 50 charachters")
	v.Check(subject.SubjectNameEnglish != "", "subject_name_english", "must be provided")
	v.Check(len(subject.SubjectNameEnglish) <= 50, "subject_name_english", "must not be more than 50 charachters")
	v.Check(subject.Stage != "", "stage", "must be provided")
	v.Check(len(subject.Stage) <= 10, "stage", "must not be more than 10 charachters")
	v.Check(subject.Semester != "", "semester", "must be provided")
	v.Check(len(subject.Semester) <= 10, "semester", "must not be more than 10 charachters")
	v.Check(subject.Department != "", "department", "must be provided")
	v.Check(len(subject.Department) <= 50, "department", "must not be more than 50 charachters")
	v.Check(subject.MaxTheoryMark != 0, "max_theory_mark", "must not be empty or zero")
	v.Check(subject.MaxLabMark != 0, "max_lab_mark", "must not be empty or zero")
	v.Check(subject.MaxSemesterMark != 0, "max_semester_mark", "must not be empty or zero")
	v.Check(subject.MaxFinalExam != 0, "max_final_exam", "must not be empty or zero")
	v.Check(subject.Credits != 0, "credits", "must not be empty or zero")
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
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&subject.CreatedAt)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m SubjectModel) Get(id int64) (*Subject, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m SubjectModel) Update(movie *Subject) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m SubjectModel) Delete(id int64) error {
	return nil
}
