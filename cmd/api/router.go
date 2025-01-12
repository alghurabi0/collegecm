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
	router.HandleFunc("GET /v1/subjects", app.getSubjects)
	router.HandleFunc("GET /v1/subjects/{id}", app.getSubjectHandler)
	router.HandleFunc("POST /v1/subjects", app.createSubjectHandler)
	router.HandleFunc("POST /v1/subjects/import", app.importSubjects)
	router.HandleFunc("PATCH /v1/subjects/{id}", app.updateSubject)
	router.HandleFunc("DELETE /v1/subjects/{id}", app.deleteSubject)
	// Return the httprouter instance.
	return headers.Then(router)
}
