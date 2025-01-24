package main

import (
	"errors"
	"net/http"

	"collegecm.hamid.net/internal/data"
)

func (app *application) getYears(w http.ResponseWriter, r *http.Request) {
	years, err := app.models.Years.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"years": years}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
