package data

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"collegecm.hamid.net/internal/validator"
)

type Year struct {
	Year string `json:"year"`
}

func ValidateYear(v *validator.Validator, year *Year) {
	// TODO - handle strings length with varchar
	v.Check(isValidAcademicYear(year.Year), "السنة الاكاديمية", "يحب ادخال سنة اكاديمية صحيحة")
}

func isValidAcademicYear(yearString string) bool {
	// Regex to match the format 'YYYY_YYYY'
	regex := regexp.MustCompile(`^\d{4}_\d{4}$`)

	// Check if the string matches the format
	if !regex.MatchString(yearString) {
		return false
	}

	// Split the string into the two years
	years := strings.Split(yearString, "_")
	year1, _ := strconv.Atoi(years[0])
	year2, _ := strconv.Atoi(years[1])

	// Check if the second year is exactly one year after the first
	return year2 == year1+1
}

type YearModel struct {
	DB *sql.DB
}

func (y YearModel) GetAll() ([]*Year, error) {
	q := `SELECT * FROM years;`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := y.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var years []*Year
	for rows.Next() {
		var year Year
		err := rows.Scan(
			&year.Year,
		)
		if err != nil {
			return nil, err
		}
		years = append(years, &year)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return years, nil
}

func (y YearModel) Insert(year *Year) error {
	studentsTablename := fmt.Sprintf("students_%s", year.Year)
	subjectsTablename := fmt.Sprintf("subjects_%s", year.Year)
	carryoverTablename := fmt.Sprintf("carryovers_%s", year.Year)
	exemptedTablename := fmt.Sprintf("exempted_%s", year.Year)
	marksTablename := fmt.Sprintf("marks_%s", year.Year)

	studentsQ := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	seq_in_college SERIAL,
    student_name VARCHAR(255) NOT NULL,
    stage VARCHAR(100) NOT NULL,
    student_id INTEGER NOT NULL PRIMARY KEY,
    state VARCHAR(100) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
	);`, studentsTablename)
	subjectsQ := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	subject_id           INTEGER NOT NULL PRIMARY KEY 
  	,subject_name         VARCHAR(100) NOT NULL
  	,subject_name_english VARCHAR(100) NOT NULL
  	,stage                VARCHAR(30) NOT NULL
  	,semester             VARCHAR(30) NOT NULL
  	,department           VARCHAR(100) NOT NULL
  	,max_theory_mark      INTEGER  NOT NULL
  	,max_lab_mark         INTEGER  NOT NULL
  	,max_semester_mark    INTEGER  NOT NULL
  	,max_final_exam       INTEGER  NOT NULL
  	,credits              INTEGER  NOT NULL
  	,active               VARCHAR(10) NOT NULL
  	,ministerial          VARCHAR(10) NOT NULL
  	,created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
	);`, subjectsTablename)
	carryoverQ := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES %s(student_id) NOT NULL ON DELETE CASCADE,
    subject_id INTEGER REFERENCES %s(subject_id) NOT NULL ON DELETE CASCADE,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (student_id, subject_id)
	);`, carryoverTablename, studentsTablename, subjectsTablename)
	exemptedQ := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES %s(student_id) NOT NULL ON DELETE CASCADE,
    subject_id INTEGER REFERENCES %s(subject_id) NOT NULL ON DELETE CASCADE,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (student_id, subject_id)
	);`, exemptedTablename, studentsTablename, subjectsTablename)
	marksQ := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES %s(student_id) NOT NULL ON DELETE CASCADE,
    subject_id INTEGER REFERENCES %s(subject_id) NOT NULL ON DELETE CASCADE,
	semester_mark INTEGER NOT NULL DEFAULT 0,
	final_mark INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (student_id, subject_id)
	);`, marksTablename, studentsTablename, subjectsTablename)
	q2 := `INSERT INTO tables (table_name) values ($1);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := createStudentsTable(ctx, y.DB, studentsQ, q2, studentsTablename)
	if err != nil {
		return err
	}
	err = createSubjectsTable(ctx, y.DB, subjectsQ, q2, subjectsTablename)
	if err != nil {
		return err
	}
	err = createCarryoversTable(ctx, y.DB, carryoverQ, q2, carryoverTablename)
	if err != nil {
		return err
	}
	err = createExemptedTable(ctx, y.DB, exemptedQ, q2, exemptedTablename)
	if err != nil {
		return err
	}
	err = createMarksTable(ctx, y.DB, marksQ, q2, marksTablename)
	if err != nil {
		return err
	}

	q := `INSERT INTO years (year) VALUES ($1);`
	args := []interface{}{
		year.Year,
	}
	_, err = y.DB.ExecContext(ctx, q, args...)
	return err
}

func createStudentsTable(ctx context.Context, db *sql.DB, query, q2, table string) error {
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	args := []interface{}{table}
	_, err = db.ExecContext(ctx, q2, args...)
	return err
}

func createSubjectsTable(ctx context.Context, db *sql.DB, query, q2, table string) error {
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	args := []interface{}{table}
	_, err = db.ExecContext(ctx, q2, args...)
	return err
}

func createCarryoversTable(ctx context.Context, db *sql.DB, query, q2, table string) error {
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	args := []interface{}{table}
	_, err = db.ExecContext(ctx, q2, args...)
	return err
}

func createExemptedTable(ctx context.Context, db *sql.DB, query, q2, table string) error {
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	args := []interface{}{table}
	_, err = db.ExecContext(ctx, q2, args...)
	return err
}

func createMarksTable(ctx context.Context, db *sql.DB, query, q2, table string) error {
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	args := []interface{}{table}
	_, err = db.ExecContext(ctx, q2, args...)
	return err
}

func (y YearModel) Delete(year string) error {
	studentsTable := fmt.Sprintf("students_%s", year)
	subjectsTable := fmt.Sprintf("subjects_%s", year)
	carryoversTable := fmt.Sprintf("carryovers_%s", year)
	exemptedTable := fmt.Sprintf("exempted_%s", year)
	marksTable := fmt.Sprintf("marks_%s", year)
	stq := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, studentsTable)
	suq := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, subjectsTable)
	cq := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, carryoversTable)
	eq := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, exemptedTable)
	mq := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, marksTable)
	q := `DELETE FROM tables WHERE table_name LIKE $1;`
	q2 := `DELETE FROM years WHERE year = $1;`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := y.DB.ExecContext(ctx, mq)
	if err != nil {
		return err
	}
	_, err = y.DB.ExecContext(ctx, eq)
	if err != nil {
		return err
	}
	_, err = y.DB.ExecContext(ctx, cq)
	if err != nil {
		return err
	}
	_, err = y.DB.ExecContext(ctx, stq)
	if err != nil {
		return err
	}
	_, err = y.DB.ExecContext(ctx, suq)
	if err != nil {
		return err
	}
	_, err = y.DB.ExecContext(ctx, q, "%"+year)
	if err != nil {
		return err
	}
	_, err = y.DB.ExecContext(ctx, q2, year)
	return err
}
