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
	router.HandleFunc("GET /v1/subjects/{year}/{stage}", app.getSubjects)
	router.HandleFunc("GET /v1/subject/{year}/{id}", app.getSubjectHandler)
	router.HandleFunc("POST /v1/subjects/{year}", app.createSubjectHandler)
	router.HandleFunc("POST /v1/subjects/import", app.importSubjects)
	router.HandleFunc("PATCH /v1/subjects/{year}/{id}", app.updateSubject)
	router.HandleFunc("DELETE /v1/subjects/{year}/{id}", app.deleteSubject)
	// students
	router.HandleFunc("GET /v1/students/{year}/{stage}", app.getStudents)
	router.HandleFunc("GET /v1/student/{year}/{id}", app.getStudent)
	router.HandleFunc("POST /v1/students/{year}", app.createStudent)
	router.HandleFunc("POST /v1/students/import", app.importstudents)
	router.HandleFunc("PATCH /v1/students/{year}/{id}", app.updateStudent)
	router.HandleFunc("DELETE /v1/students/{year}/{id}", app.deleteStudent)
	// carryovers
	router.HandleFunc("GET /v1/carryovers/{year}/{stage}", app.getCarryovers)
	router.HandleFunc("GET /v1/carryover/{year}/{id}", app.getCarryover)
	router.HandleFunc("GET /v1/carryovers/find/{year}/{student_id}/{subject_id}", app.findCarryover)
	router.HandleFunc("GET /v1/carryovers/subjects/{year}/{id}", app.getSubjectsCarryovers)
	router.HandleFunc("GET /v1/carryovers/students/{year}/{id}", app.getStudentsCarryovers)
	router.HandleFunc("POST /v1/carryovers/{year}", app.createCarryover)
	router.HandleFunc("DELETE /v1/carryovers/{year}/{id}", app.deleteCarryover)
	// exempteds
	router.HandleFunc("GET /v1/exempteds/{year}/{stage}", app.getExempteds)
	router.HandleFunc("GET /v1/exempted/{year}/{id}", app.getExempted)
	router.HandleFunc("GET /v1/exempteds/find/{student_id}/{subject_id}", app.findExempted)
	router.HandleFunc("GET /v1/exempteds/subjects/{year}/{id}", app.getSubjectsExempteds)
	router.HandleFunc("GET /v1/exempteds/students/{year}/{id}", app.getStudentsExempteds)
	router.HandleFunc("POST /v1/exempteds/{year}", app.createExempted)
	router.HandleFunc("DELETE /v1/exempteds/{year}/{id}", app.deleteExempted)
	// marks
	router.HandleFunc("GET /v1/marks/{year}/{stage}", app.getMarks)
	router.HandleFunc("GET /v1/mark/{year}/{id}", app.getMark)
	router.HandleFunc("POST /v1/marks/{year}", app.createMark)
	router.HandleFunc("PATCH /v1/marks/{year}/{id}", app.updateMark)
	router.HandleFunc("DELETE /v1/marks/{year}/{id}", app.deleteMark)
	// custom
	router.HandleFunc("GET /v1/custom/{yaer}/{id}", app.getStudentData)
	// general
	router.HandleFunc("GET /v1/years", app.getYears)
	// Return the httprouter instance.
	return headers.Then(router)
}
