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
	standard := alice.New(app.secureHeaders)
	auth := alice.New(app.sessionManager.LoadAndSave, app.isLoggedIn)
	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) { fmt.Println("options req") })
	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	// subjects
	router.Handle("GET /v1/subjects/{year}/{stage}", auth.ThenFunc(app.getSubjects))
	router.Handle("GET /v1/subject/{year}/{id}", auth.ThenFunc(app.getSubjectHandler))
	router.Handle("POST /v1/subjects/{year}", auth.ThenFunc(app.createSubjectHandler))
	router.Handle("POST /v1/subjects/import", auth.ThenFunc(app.importSubjects))
	router.Handle("PATCH /v1/subjects/{year}/{id}", auth.ThenFunc(app.updateSubject))
	router.Handle("DELETE /v1/subjects/{year}/{id}", auth.ThenFunc(app.deleteSubject))
	// students
	router.Handle("GET /v1/students/{year}/{stage}", auth.ThenFunc(app.getStudents))
	router.Handle("GET /v1/student/{year}/{id}", auth.ThenFunc(app.getStudent))
	router.Handle("POST /v1/students/{year}", auth.ThenFunc(app.createStudent))
	router.Handle("POST /v1/students/import", auth.ThenFunc(app.importstudents))
	router.Handle("PATCH /v1/students/{year}/{id}", auth.ThenFunc(app.updateStudent))
	router.Handle("DELETE /v1/students/{year}/{id}", auth.ThenFunc(app.deleteStudent))
	// carryovers
	router.Handle("GET /v1/carryovers/{year}/{stage}", auth.ThenFunc(app.getCarryovers))
	router.Handle("GET /v1/carryover/{year}/{id}", auth.ThenFunc(app.getCarryover))
	router.Handle("GET /v1/carryovers/find/{year}/{student_id}/{subject_id}", auth.ThenFunc(app.findCarryover))
	router.Handle("GET /v1/carryovers/subjects/{year}/{id}", auth.ThenFunc(app.getSubjectsCarryovers))
	router.Handle("GET /v1/carryovers/students/{year}/{id}", auth.ThenFunc(app.getStudentsCarryovers))
	router.Handle("POST /v1/carryovers/{year}", auth.ThenFunc(app.createCarryover))
	router.Handle("DELETE /v1/carryovers/{year}/{id}", auth.ThenFunc(app.deleteCarryover))
	// exempteds
	router.Handle("GET /v1/exempteds/{year}/{stage}", auth.ThenFunc(app.getExempteds))
	router.Handle("GET /v1/exempted/{year}/{id}", auth.ThenFunc(app.getExempted))
	router.Handle("GET /v1/exempteds/find/{student_id}/{subject_id}", auth.ThenFunc(app.findExempted))
	router.Handle("GET /v1/exempteds/subjects/{year}/{id}", auth.ThenFunc(app.getSubjectsExempteds))
	router.Handle("GET /v1/exempteds/students/{year}/{id}", auth.ThenFunc(app.getStudentsExempteds))
	router.Handle("POST /v1/exempteds/{year}", auth.ThenFunc(app.createExempted))
	router.Handle("DELETE /v1/exempteds/{year}/{id}", auth.ThenFunc(app.deleteExempted))
	// marks
	router.Handle("GET /v1/marks/{year}/{stage}", auth.ThenFunc(app.getMarks))
	router.Handle("GET /v1/mark/{year}/{id}", auth.ThenFunc(app.getMark))
	router.Handle("POST /v1/marks/{year}", auth.ThenFunc(app.createMark))
	router.Handle("PATCH /v1/marks/{year}/{id}", auth.ThenFunc(app.updateMark))
	router.Handle("DELETE /v1/marks/{year}/{id}", auth.ThenFunc(app.deleteMark))
	// users
	router.Handle("GET /v1/users", auth.ThenFunc(app.getUsers))
	router.Handle("GET /v1/users/{id}", auth.ThenFunc(app.getUser))
	router.Handle("POST /v1/users", auth.ThenFunc(app.createUser))
	router.Handle("PATCH /v1/users/{id}", auth.ThenFunc(app.updateUser))
	router.Handle("DELETE /v1/users/{id}", auth.ThenFunc(app.deleteUser))
	// privileges
	router.Handle("GET /v1/privileges/{id}", auth.ThenFunc(app.getPrivileges))
	router.Handle("POST /v1/privileges", auth.ThenFunc(app.createPrivilege))
	// auth
	router.HandleFunc("GET /v1/auth/status", app.authStatus)
	router.HandleFunc("POST /v1/login", app.login)
	router.Handle("POST /v1/logout", auth.ThenFunc(app.logout))
	// custom
	router.Handle("GET /v1/custom/{year}/{id}", auth.ThenFunc(app.getStudentData))
	// general
	router.Handle("GET /v1/years", auth.ThenFunc(app.getYears))
	// Return the httprouter instance.
	return standard.Then(router)
}
