package main

import (
	"errors"
	"net/http"

	"collegecm.hamid.net/internal/data"
)

func (app *application) getStudentData(w http.ResponseWriter, r *http.Request) {
	year, id, err := app.readIdYearParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	studentData, err := app.models.Customs.GetStudentData(year, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"student_data": studentData}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
