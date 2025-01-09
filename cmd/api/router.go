package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := http.NewServeMux()
	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	router.HandleFunc("GET /v1/subjects", app.getSubjects)
	router.HandleFunc("GET /v1/subjects/{id}", app.getSubjectHandler)
	router.HandleFunc("POST /v1/subjects", app.createSubjectHandler)
	router.HandleFunc("PATCH /v1/subjects/{id}", app.updateSubject)
	router.HandleFunc("DELETE /v1/subjects/{id}", app.deleteSubject)

	router.HandleFunc("OPTIONS /v1/subjects/{id}", app.subjectsPreflightHandler)
	router.HandleFunc("OPTIONS /v1/subjects", app.subjectsPreflightHandler)

	// Return the httprouter instance.
	return router
}

func (app *application) subjectsPreflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}
