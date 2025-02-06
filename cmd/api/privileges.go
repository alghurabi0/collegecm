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
		UserId    int    `json:"user_id"`
		Year      string `json:"year"`
		TableName string `json:"table_name"`
		Stage     string `json:"stage"`
		SubjectId *int   `json:"subject_id"`
		CanRead   bool   `json:"can_read"`
		CanWrite  bool   `json:"can_write"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	privilege := &data.Privilege{
		UserId:   input.UserId,
		Year:     input.Year,
		Stage:    input.Stage,
		CanRead:  input.CanRead,
		CanWrite: input.CanWrite,
	}
	if input.TableName == "" {
		app.serverErrorResponse(w, r, errors.New("empty table name"))
		return
	} else if input.TableName == "all" {
		privilege.TableId = -1
	} else {
		table, err := app.models.Tables.GetByName(input.TableName, privilege.Year)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		privilege.TableId = int(table.ID)
	}
	if input.SubjectId != nil {
		privilege.SubjectId = *input.SubjectId
	} else {
		privilege.SubjectId = -1
	}
	v := validator.New()
	if data.ValidatePrivilege(v, privilege); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	if input.TableName == "users" || input.TableName == "privileges" || input.TableName == "years" {
		privilege.SubjectId = -1
		privilege.Stage = "all"
		privilege.Year = "all"
	}
	err = app.models.Privileges.Insert(privilege)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if input.TableName == "users" || input.TableName == "privileges" || input.TableName == "years" {
		privilege.TableName = input.TableName
	} else {
		privilege.TableName = input.TableName + "_" + input.Year
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"privilege": privilege}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePrivilege(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId  int    `json:"user_id"`
		Year    string `json:"year"`
		TableId int    `json:"table_id"`
		Stage   string `json:"stage"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	privilege := &data.Privilege{
		UserId:  input.UserId,
		Year:    input.Year,
		TableId: input.TableId,
		Stage:   input.Stage,
	}
	err = app.models.Privileges.Delete(privilege)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
