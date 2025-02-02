package main

import (
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
	user, err := app.models.Users.Get(userId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"privileges": privileges, "user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createPrivilege(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId    int     `json:"user_id"`
		Year      string  `json:"year"`
		TableName *string `json:"table_name"`
		Stage     *string `json:"stage"`
		SubjectId *int    `json:"subject_id"`
		CanRead   *bool   `json:"can_read"`
		CanWrite  *bool   `json:"can_write"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	privilege := &data.Privilege{
		UserId: input.UserId,
		Year:   input.Year,
	}
	v := validator.New()
	if data.ValidatePrivilege(v, privilege); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	if input.Stage != nil {
		privilege.Stage = *input.Stage
	} else {
		privilege.Stage = "none"
	}
	if input.SubjectId != nil {
		privilege.SubjectId = *input.SubjectId
	} else {
		privilege.SubjectId = 0
	}
	if input.CanRead != nil {
		privilege.CanRead = *input.CanRead
	} else {
		privilege.CanRead = false
	}
	if input.CanWrite != nil {
		privilege.CanWrite = *input.CanWrite
	} else {
		privilege.CanWrite = false
	}
	if input.TableName != nil {
		table, err := app.models.Tables.GetByName(*input.TableName, privilege.Year)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		privilege.TableId = int(table.ID)
	} else {
		privilege.TableId = 0
	}
	v = validator.New()
	if data.ValidatePrivilegeFull(v, privilege); !v.Valid() {
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
