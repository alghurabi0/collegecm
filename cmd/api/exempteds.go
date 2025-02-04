package main

import (
	"errors"
	"net/http"

	"collegecm.hamid.net/internal/data"
	"collegecm.hamid.net/internal/validator"
)

func (app *application) getExempteds(w http.ResponseWriter, r *http.Request) {
	year, err := app.getYearFromContext(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	stage, err := app.getStageFromContext(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	exempteds, err := app.models.Exempteds.GetAll(year, stage)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"exempteds": exempteds}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// func (app *application) getExempted(w http.ResponseWriter, r *http.Request) {
// 	//id
// 	year, id, err := app.readIdYearParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	exempted, err := app.models.Exempteds.Get(year, id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	err = app.writeJSON(w, http.StatusOK, envelope{"exempted": exempted}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) findExempted(w http.ResponseWriter, r *http.Request) {
// 	student_idStr := r.PathValue("student_id")
// 	if strings.TrimSpace(student_idStr) == "" {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	student_id, err := strconv.ParseInt(student_idStr, 10, 64)
// 	if err != nil || student_id < 1 {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	subject_idStr := r.PathValue("subject_id")
// 	if strings.TrimSpace(subject_idStr) == "" {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	subject_id, err := strconv.ParseInt(subject_idStr, 10, 64)
// 	if err != nil || subject_id < 1 {
// 		app.notFoundResponse(w, r)
// 		return
// 	}

// 	exempted, err := app.models.Exempteds.Find(student_id, subject_id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	err = app.writeJSON(w, http.StatusOK, envelope{"exempted": exempted}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) getSubjectsExempteds(w http.ResponseWriter, r *http.Request) {
// 	year, id, err := app.readIdYearParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	exempteds, err := app.models.Exempteds.GetSubjects(year, id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	err = app.writeJSON(w, http.StatusOK, envelope{"subjects_exempteds": exempteds}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) getStudentsExempteds(w http.ResponseWriter, r *http.Request) {
// 	year, id, err := app.readIdYearParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	exempteds, err := app.models.Exempteds.GetStudents(year, id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	err = app.writeJSON(w, http.StatusOK, envelope{"students_exempteds": exempteds}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

func (app *application) createExempted(w http.ResponseWriter, r *http.Request) {
	year, err := app.getYearFromContext(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	var input struct {
		StudentId int64 `json:"student_id"`
		SubjectId int64 `json:"subject_id"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// privilege check
	stage, err := app.models.Students.GetStage(input.StudentId, year)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user, err := app.getUserFromContext(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	hasAccess, err := app.models.Privileges.CheckWriteAccess(int(user.ID), "exempted_"+year, stage)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !hasAccess {
		app.unauthorized(w, r)
		return
	}

	exempted := &data.Exempted{
		StudentId: input.StudentId,
		SubjectId: input.SubjectId,
	}
	v := validator.New()
	if data.ValidateExempted(v, exempted); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Exempteds.Insert(year, exempted)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	student, err := app.models.Students.Get(year, exempted.StudentId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	subject, err := app.models.Subjects.Get(year, exempted.SubjectId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	exempted.StudentName = student.StudentName
	exempted.SubjectName = subject.SubjectName
	err = app.writeJSON(w, http.StatusCreated, envelope{"exempted": exempted}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteExempted(w http.ResponseWriter, r *http.Request) {
	id, err := app.getIdFromContext(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	year, err := app.getYearFromContext(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Exempteds.Delete(year, id)
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
