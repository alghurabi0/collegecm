package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"collegecm.hamid.net/internal/data"
	"collegecm.hamid.net/internal/validator"
)

func (app *application) getStudents(w http.ResponseWriter, r *http.Request) {
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
	students, err := app.models.Students.GetAll(year, stage)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"students": students}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// func (app *application) getStudent(w http.ResponseWriter, r *http.Request) {
// 	year, id, err := app.readIdYearParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	student, err := app.models.Students.Get(year, id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	err = app.writeJSON(w, http.StatusOK, envelope{"student": student}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

func (app *application) createStudent(w http.ResponseWriter, r *http.Request) {
	year, err := app.getYearFromContext(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	var input struct {
		StudentName string `json:"student_name"`
		Stage       string `json:"stage"`
		StudentId   int    `json:"student_id"`
		State       string `json:"state"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	student := &data.Student{
		StudentName: input.StudentName,
		Stage:       input.Stage,
		StudentId:   input.StudentId,
		State:       input.State,
	}
	v := validator.New()
	if data.ValidateStudent(v, student); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Students.Insert(year, student)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"student": student}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateStudent(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	year, err := app.getYearFromContext(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	student, err := app.models.Students.Get(year, id)
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
		StudentName *string `json:"student_name"`
		Stage       *string `json:"stage"`
		StudentId   *int    `json:"student_id"`
		State       *string `json:"state"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.StudentName != nil {
		student.StudentName = *input.StudentName
	}
	if input.Stage != nil {
		student.Stage = *input.Stage
	}
	if input.StudentId != nil {
		student.StudentId = *input.StudentId
	}
	if input.State != nil {
		student.State = *input.State
	}
	v := validator.New()
	if data.ValidateStudent(v, student); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Students.Update(year, student)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"student": student}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteStudent(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	year, err := app.getYearFromContext(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Students.Delete(year, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "تم حذف الطالب بنجاح"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) importstudents(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "الحد الاقصى لحجم الملف هو mb 10 ")
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "لم يتم ارفاق ملف")
		return
	}
	defer file.Close()
	students := []*data.Student{}
	err = app.processFile(&file, &students)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "حدث خطأ, يرجى التواصل مع الدعم")
		fmt.Println(err)
		return
	}
	allErrors := make(map[string]string)
	v := validator.New()
	for i, student := range students {
		// validate
		v.Errors = make(map[string]string)
		if data.ValidateStudent(v, student); !v.Valid() {
			var errorMsgs []string
			for key, msg := range v.Errors {
				errorMsgs = append(errorMsgs, key+": "+msg)
			}
			allErrors[fmt.Sprintf("row-%d", i+1)] = strings.Join(errorMsgs, ", ")
			continue
		}
		err = app.models.Students.Insert("", student)
		if err != nil {
			allErrors[fmt.Sprintf("row-%d", i+1)] = "رقم الطالب مكرر او حدث خطأ"
		}
	}
	// get all subjects or redirect
	allStudents, err := app.models.Students.GetAll("", "")
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if len(allErrors) > 0 {
		err = app.writeJSON(w, http.StatusOK, envelope{"students": allStudents, "errors": allErrors}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	} else {
		err = app.writeJSON(w, http.StatusOK, envelope{"students": allStudents}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
