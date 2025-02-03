package main

import (
	"errors"
	"net/http"

	"collegecm.hamid.net/internal/data"
)

func (app *application) getStudentData(w http.ResponseWriter, r *http.Request) {
	year, err := app.getYearFromContext(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	id, err := app.getIdFromContext(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	privileges, err := app.getCustomPrivsFromContext(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	studentData, err := app.models.Customs.GetStudentData(year, id, privileges)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	student, err := app.getStudentFromContext(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	studentData.Student = student
	err = app.writeJSON(w, http.StatusOK, envelope{"student_data": studentData}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
