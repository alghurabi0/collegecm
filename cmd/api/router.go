package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := http.NewServeMux()
	// middleware chain
	headers := alice.New(app.secureHeaders)
	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) { fmt.Println("options req") })
	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	// subjects
	router.HandleFunc("GET /v1/subjects", app.getSubjects)
	router.HandleFunc("GET /v1/subjects/{id}", app.getSubjectHandler)
	router.HandleFunc("POST /v1/subjects", app.createSubjectHandler)
	router.HandleFunc("POST /v1/subjects/import", app.importSubjects)
	router.HandleFunc("PATCH /v1/subjects/{id}", app.updateSubject)
	router.HandleFunc("DELETE /v1/subjects/{id}", app.deleteSubject)
	// students
	router.HandleFunc("GET /v1/students", app.getStudents)
	router.HandleFunc("GET /v1/students/{id}", app.getStudent)
	router.HandleFunc("POST /v1/students", app.createStudent)
	router.HandleFunc("POST /v1/students/import", app.importstudents)
	router.HandleFunc("PATCH /v1/students/{id}", app.updateStudent)
	router.HandleFunc("DELETE /v1/students/{id}", app.deleteStudent)
	// carryovers
	router.HandleFunc("GET /v1/carryovers", app.getCarryovers)
	router.HandleFunc("GET /v1/carryovers/{id}", app.getCarryover)
	router.HandleFunc("GET /v1/carryovers/{student_id}/{subject_id}", app.findCarryover)
	router.HandleFunc("GET /v1/carryovers/subjects/{id}", app.getSubjectsCarryovers)
	router.HandleFunc("GET /v1/carryovers/students/{id}", app.getStudentsCarryovers)
	router.HandleFunc("POST /v1/carryovers", app.createCarryover)
	router.HandleFunc("DELETE /v1/carryovers/{id}", app.deleteCarryover)
	// exempteds
	router.HandleFunc("GET /v1/exempteds", app.getExempteds)
	router.HandleFunc("GET /v1/exempteds/{id}", app.getExempted)
	router.HandleFunc("GET /v1/exempteds/{student_id}/{subject_id}", app.findExempted)
	router.HandleFunc("GET /v1/exempteds/subjects/{id}", app.getSubjectsExempteds)
	router.HandleFunc("GET /v1/exempteds/students/{id}", app.getStudentsExempteds)
	router.HandleFunc("POST /v1/exempteds", app.createExempted)
	router.HandleFunc("DELETE /v1/exempteds/{id}", app.deleteExempted)
	// marks
	router.HandleFunc("GET /v1/marks", app.getMarks)
	router.HandleFunc("GET /v1/marks/{id}", app.getMark)
	router.HandleFunc("POST /v1/marks", app.createMark)
	router.HandleFunc("PATCH /v1/marks/{id}", app.updateMark)
	router.HandleFunc("DELETE /v1/marks/{id}", app.deleteMark)
	// Return the httprouter instance.
	return headers.Then(router)
}
