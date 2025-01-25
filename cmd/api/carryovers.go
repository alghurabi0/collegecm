package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"collegecm.hamid.net/internal/data"
	"collegecm.hamid.net/internal/validator"
)

func (app *application) getCarryovers(w http.ResponseWriter, r *http.Request) {
	year, stage, err := app.readParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	carryovers, err := app.models.Carryovers.GetAll(year, stage)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"carryovers": carryovers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCarryover(w http.ResponseWriter, r *http.Request) {
	//id
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	carryover, err := app.models.Carryovers.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"carryover": carryover}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) findCarryover(w http.ResponseWriter, r *http.Request) {
	student_idStr := r.PathValue("student_id")
	if strings.TrimSpace(student_idStr) == "" {
		app.notFoundResponse(w, r)
		return
	}
	student_id, err := strconv.ParseInt(student_idStr, 10, 64)
	if err != nil || student_id < 1 {
		app.notFoundResponse(w, r)
		return
	}
	subject_idStr := r.PathValue("subject_id")
	if strings.TrimSpace(subject_idStr) == "" {
		app.notFoundResponse(w, r)
		return
	}
	subject_id, err := strconv.ParseInt(subject_idStr, 10, 64)
	if err != nil || subject_id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	carryover, err := app.models.Carryovers.Find(student_id, subject_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"carryover": carryover}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getSubjectsCarryovers(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	carryovers, err := app.models.Carryovers.GetSubjects(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"subjects_carryovers": carryovers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getStudentsCarryovers(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	carryovers, err := app.models.Carryovers.GetStudents(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"students_carryovers": carryovers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createCarryover(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StudentId int64 `json:"student_id"`
		SubjectId int64 `json:"subject_id"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	carryover := &data.Carryover{
		StudentId: input.StudentId,
		SubjectId: input.SubjectId,
	}
	v := validator.New()
	if data.ValidateCarryover(v, carryover); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Carryovers.Insert(carryover)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	student, err := app.models.Students.Get(carryover.StudentId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	subject, err := app.models.Subjects.Get(carryover.SubjectId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	carryover.StudentName = student.StudentName
	carryover.SubjectName = subject.SubjectName
	err = app.writeJSON(w, http.StatusCreated, envelope{"carryover": carryover}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCarryover(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Carryovers.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "تم الحذف بنجاح"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
