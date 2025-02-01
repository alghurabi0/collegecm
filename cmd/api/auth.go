package main

import (
	"errors"
	"fmt"
	"net/http"

	"collegecm.hamid.net/internal/data"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Username == "" || input.Password == "" {
		app.badRequestResponse(w, r, errors.New("invalid username or password"))
		return
	}
	user, err := app.models.Users.GetByUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if user.Password != input.Password {
		app.badRequestResponse(w, r, errors.New("invalid username or password"))
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "userID", int(user.ID))
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "userID")
	w.WriteHeader(http.StatusOK)
}

func (app *application) authStatus(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), "userID")
	if userId == 0 {
		fmt.Println("user not logged in")
		app.notFoundResponse(w, r)
		return
	}
	user, err := app.models.Users.Get(int64(userId))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
