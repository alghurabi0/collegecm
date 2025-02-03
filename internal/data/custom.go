package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Custom struct {
	Student    *Student     `json:"student"`
	Subjects   []*Subject   `json:"subjects"`
	Exempteds  []*Exempted  `json:"exempteds"`
	Carryovers []*Carryover `json:"carryovers"`
	Marks      []*Mark      `json:"marks"`
}

type CustomModel struct {
	DB *sql.DB
}

func (c CustomModel) GetStudentData(year string, id int64, privs *CustomPrivilegeAccess) (*Custom, error) {
	carryoverTablename := fmt.Sprintf("carryovers_%s", year)
	exemptedTablename := fmt.Sprintf("exempted_%s", year)
	marksTablename := fmt.Sprintf("marks_%s", year)
	subjectsTablename := fmt.Sprintf("subjects_%s", year)
	studentsTablename := fmt.Sprintf("students_%s", year)

	carryoverQ := fmt.Sprintf(`
	SELECT c.id, s.subject_name
	FROM %s c
	JOIN %s s ON c.subject_id = s.subject_id
	WHERE c.student_id = $1;`, carryoverTablename, subjectsTablename)
	exemptedQ := fmt.Sprintf(`
	SELECT e.id, s.subject_name
	FROM %s e
	JOIN %s s ON e.subject_id = s.subject_id
	WHERE e.student_id = $1;`, exemptedTablename, subjectsTablename)
	marksQ := fmt.Sprintf(`
	SELECT
	m.id, s.subject_name, s.max_semester_mark, m.semester_mark, s.max_final_exam, m.final_mark
	FROM %s m
	JOIN %s s ON m.subject_id = s.subject_id
	WHERE m.student_id = $1;`, marksTablename, subjectsTablename)
	subjectsByStageQ := fmt.Sprintf(`
	SELECT subject_id, subject_name
	FROM %s
	WHERE stage = (
		SELECT stage
		FROM %s
		WHERE student_id = $1
	);`, subjectsTablename, studentsTablename)
	// studentInfoQ := fmt.Sprintf(`
	// SELECT student_id, student_name, stage
	// FROM %s
	// WHERE student_id = $1;`, studentsTablename)
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var carryovers []*Carryover
	var exempteds []*Exempted
	var marks []*Mark
	var subjectsByStage []*Subject
	var err error
	// Fetch carryovers
	if privs.Carryovers {
		carryovers, err = fetchCarryovers(ctx, c.DB, carryoverQ, id)
		if err != nil {
			return nil, err
		}
	}
	// Fetch exempteds
	if privs.Exempted {
		exempteds, err = fetchExempteds(ctx, c.DB, exemptedQ, id)
		if err != nil {
			return nil, err
		}
	}
	// Fetch marks
	if privs.Marks {
		marks, err = fetchMarks(ctx, c.DB, marksQ, id)
		if err != nil {
			return nil, err
		}
	}
	// Fetch subjects by stage
	if privs.Subjects {
		subjectsByStage, err = fetchSubjectsByStage(ctx, c.DB, subjectsByStageQ, id)
		if err != nil {
			return nil, err
		}
	}
	// Fetch student information
	// studentInfo, err := fetchStudentInfo(ctx, c.DB, studentInfoQ, id)
	// if err != nil {
	// 	return nil, err
	// }

	// Combine into a custom struct
	custom := &Custom{
		Carryovers: carryovers,
		Exempteds:  exempteds,
		Marks:      marks,
		Subjects:   subjectsByStage,
		//Student:    studentInfo,
	}

	return custom, nil
}

// Helper function to fetch subjects by stage
func fetchSubjectsByStage(ctx context.Context, db *sql.DB, query string, id int64) ([]*Subject, error) {
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []*Subject
	for rows.Next() {
		var subject Subject
		if err := rows.Scan(&subject.ID, &subject.SubjectName); err != nil {
			return nil, err
		}
		subjects = append(subjects, &subject)
	}

	return subjects, rows.Err()
}

// Helper function to fetch student info
// func fetchStudentInfo(ctx context.Context, db *sql.DB, query string, id int64) (*Student, error) {
// 	row := db.QueryRowContext(ctx, query, id)

// 	var student Student
// 	err := row.Scan(&student.StudentId, &student.StudentName, &student.Stage)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &student, nil
// }

// Helper function to fetch carryovers
func fetchCarryovers(ctx context.Context, db *sql.DB, query string, id int64) ([]*Carryover, error) {
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carryovers []*Carryover
	for rows.Next() {
		var carryover Carryover
		if err := rows.Scan(&carryover.Id, &carryover.SubjectName); err != nil {
			return nil, err
		}
		carryovers = append(carryovers, &carryover)
	}

	return carryovers, rows.Err()
}

// Helper function to fetch exempteds
func fetchExempteds(ctx context.Context, db *sql.DB, query string, id int64) ([]*Exempted, error) {
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exempteds []*Exempted
	for rows.Next() {
		var exempted Exempted
		if err := rows.Scan(&exempted.Id, &exempted.SubjectName); err != nil {
			return nil, err
		}
		exempteds = append(exempteds, &exempted)
	}

	return exempteds, rows.Err()
}

// Helper function to fetch marks
func fetchMarks(ctx context.Context, db *sql.DB, query string, id int64) ([]*Mark, error) {
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var marks []*Mark
	for rows.Next() {
		var mark Mark
		if err := rows.Scan(&mark.Id, &mark.SubjectName, &mark.MaxSemesterMark, &mark.SemesterMark, &mark.MaxFinalExam, &mark.FinalMark); err != nil {
			return nil, err
		}
		marks = append(marks, &mark)
	}

	return marks, rows.Err()
}
