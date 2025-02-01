package main

import (
	"errors"
	"fmt"
	"net/http"

	"collegecm.hamid.net/internal/data"
	"collegecm.hamid.net/internal/validator"
)

func (app *application) getMarks(w http.ResponseWriter, r *http.Request) {
	year, stage, err := app.readParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	marks, err := app.models.Marks.GetAll(year, stage)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"marks": marks}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getMark(w http.ResponseWriter, r *http.Request) {
	//id
	year, id, err := app.readIdYearParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	mark, err := app.models.Marks.Get(year, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"mark": mark}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createMark(w http.ResponseWriter, r *http.Request) {
	year, err := app.readYearParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	var input struct {
		StudentId    int64 `json:"student_id"`
		SubjectId    int64 `json:"subject_id"`
		SemesterMark *int  `json:"semester_mark"`
		FinalMark    *int  `json:"final_mark"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	mark := &data.Mark{
		StudentId: input.StudentId,
		SubjectId: input.SubjectId,
	}
	if input.SemesterMark != nil {
		mark.SemesterMark = *input.SemesterMark
	} else {
		mark.SemesterMark = 0
	}
	if input.FinalMark != nil {
		mark.FinalMark = *input.FinalMark
	} else {
		mark.FinalMark = 0
	}
	subject, err := app.models.Subjects.Get(year, mark.SubjectId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateMark(v, mark, subject.MaxSemesterMark, subject.MaxFinalExam); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Marks.Insert(year, mark)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	student, err := app.models.Students.Get(year, mark.StudentId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	mark.StudentName = student.StudentName
	mark.SubjectName = subject.SubjectName
	mark.MaxSemesterMark = subject.MaxSemesterMark
	mark.MaxFinalExam = subject.MaxFinalExam
	err = app.writeJSON(w, http.StatusCreated, envelope{"mark": mark}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMark(w http.ResponseWriter, r *http.Request) {
	year, id, err := app.readIdYearParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		fmt.Println(err)
		return
	}
	mark, err := app.models.Marks.GetRaw(year, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		SemesterMark *int `json:"semester_mark"`
		FinalMark    *int `json:"final_mark"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.SemesterMark != nil {
		mark.SemesterMark = *input.SemesterMark
	}
	if input.FinalMark != nil {
		mark.FinalMark = *input.FinalMark
	}
	subject, err := app.models.Subjects.Get(year, mark.SubjectId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateMark(v, mark, subject.MaxSemesterMark, subject.MaxFinalExam); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Marks.Update(year, mark)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	newMark, err := app.models.Marks.Get(year, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"mark": newMark}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMark(w http.ResponseWriter, r *http.Request) {
	year, id, err := app.readIdYearParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Marks.Delete(year, id)
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
