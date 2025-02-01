package main

import (
	"database/sql"
	"errors"
	"net/http"

	"collegecm.hamid.net/internal/data"
	"collegecm.hamid.net/internal/validator"
)

func (app *application) getPrivileges(w http.ResponseWriter, r *http.Request) {
	userId, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	privileges, err := app.models.Privileges.GetAll(int(userId))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"privileges": privileges}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createPrivilege(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId    int     `json:"user_id"`
		TableId   int     `json:"table_id"`
		Stage     *string `json:"stage"`
		SubjectId *int    `json:"subject_id"`
		CanRead   bool    `json:"can_read"`
		CanWrite  bool    `json:"can_write"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	privilege := &data.Privilege{
		UserId:    input.UserId,
		TableId:   input.TableId,
		Stage:     sql.NullString{String: "", Valid: false},
		SubjectId: sql.NullInt64{Int64: 0, Valid: false},
		CanRead:   input.CanRead,
		CanWrite:  input.CanWrite,
	}
	if input.Stage != nil {
		privilege.Stage = sql.NullString{String: *input.Stage, Valid: true}
	}
	if input.SubjectId != nil {
		privilege.SubjectId = sql.NullInt64{Int64: int64(*input.SubjectId), Valid: true}
	}
	v := validator.New()
	if data.ValidatePrivilege(v, privilege); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Privileges.Insert(privilege)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"privilege": privilege}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
