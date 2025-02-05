package main

import (
	"errors"
	"net/http"

	"collegecm.hamid.net/internal/data"
	"collegecm.hamid.net/internal/validator"
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

func (app *application) createYear(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Year string `json:"year"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	year := &data.Year{
		Year: input.Year,
	}
	v := validator.New()
	if data.ValidateYear(v, year); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Years.Insert(year)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"year": year}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
