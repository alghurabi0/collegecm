package main

import (
	"errors"
	"fmt"
	"net/http"

	"collegecm.hamid.net/internal/data"
	"collegecm.hamid.net/internal/validator"
)

func (app *application) getSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := app.models.Subjects.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	header := make(http.Header)
	header.Set("Access-Control-Allow-Origin", "*")                            // Allow the Vue.js app to make requests
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")      // Allow these HTTP methods
	header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Allow these headers
	header.Set("Access-Control-Allow-Credentials", "true")
	err = app.writeJSON(w, http.StatusOK, envelope{"subjects": subjects}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) getSubjectHandler(w http.ResponseWriter, r *http.Request) {
	//id
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	subject, err := app.models.Subjects.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"subject": subject}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createSubjectHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Movie struct that we created earlier). This struct will be our *target
	// decode destination*.
	var input struct {
		ID                 int    `json:"subject_id"`
		SubjectName        string `json:"subject_name"`
		SubjectNameEnglish string `json:"subject_name_english"`
		Stage              string `json:"stage"`
		Semester           string `json:"semester"`
		Department         string `json:"department"`
		MaxTheoryMark      int    `json:"max_theory_mark"`
		MaxLabMark         int    `json:"max_lab_mark"`
		MaxSemesterMark    int    `json:"max_semester_mark"`
		MaxFinalExam       int    `json:"max_final_exam"`
		Credits            int    `json:"credits"`
		Active             string `json:"active"`
		Ministerial        string `json:"ministerial"`
	}
	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use the Decode() method to decode the body contents into the input struct.
	// Importantly, notice that when we call Decode() we pass a *pointer* to the input
	// struct as the target decode destination. If there was an error during decoding,
	// we also use our generic errorResponse() helper to send the client a 400 Bad
	// Request response containing the error message.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		fmt.Println(err)
		return
	}
	subject := &data.Subject{
		ID:                 input.ID,
		SubjectName:        input.SubjectName,
		SubjectNameEnglish: input.SubjectNameEnglish,
		Stage:              input.Stage,
		Semester:           input.Semester,
		Department:         input.Department,
		MaxTheoryMark:      input.MaxTheoryMark,
		MaxLabMark:         input.MaxLabMark,
		MaxSemesterMark:    input.MaxSemesterMark,
		MaxFinalExam:       input.MaxFinalExam,
		Credits:            input.Credits,
		Active:             input.Active,
		Ministerial:        input.Ministerial,
	}
	// Initialize a new Validator.
	v := validator.New()
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateSubject(v, subject); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Subjects.Insert(subject)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	//headers := make(http.Header)
	//headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
	// Dump the contents of the input struct in a HTTP response.
	header := make(http.Header)
	header.Set("Access-Control-Allow-Origin", "*")                                       // Allow the Vue.js app to make requests
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS") // Allow these HTTP methods
	header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")            // Allow these headers
	header.Set("Access-Control-Allow-Credentials", "true")
	err = app.writeJSON(w, http.StatusCreated, envelope{"subject": subject}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateSubject(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	subject, err := app.models.Subjects.Get(id)
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
		ID                 *int    `json:"subject_id"`
		SubjectName        *string `json:"subject_name"`
		SubjectNameEnglish *string `json:"subject_name_english"`
		Stage              *string `json:"stage"`
		Semester           *string `json:"semester"`
		Department         *string `json:"department"`
		MaxTheoryMark      *int    `json:"max_theory_mark"`
		MaxLabMark         *int    `json:"max_lab_mark"`
		MaxSemesterMark    *int    `json:"max_semester_mark"`
		MaxFinalExam       *int    `json:"max_final_exam"`
		Credits            *int    `json:"credits"`
		Active             *string `json:"active"`
		Ministerial        *string `json:"ministerial"`
	}
	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.ID != nil {
		subject.ID = *input.ID
	}
	if input.SubjectName != nil {
		subject.SubjectName = *input.SubjectName
	}
	if input.SubjectNameEnglish != nil {
		subject.SubjectNameEnglish = *input.SubjectNameEnglish
	}
	if input.Stage != nil {
		subject.Stage = *input.Stage
	}
	if input.Semester != nil {
		subject.Semester = *input.Semester
	}
	if input.Department != nil {
		subject.Department = *input.Department
	}
	if input.MaxTheoryMark != nil {
		subject.MaxTheoryMark = *input.MaxTheoryMark
	}
	if input.MaxLabMark != nil {
		subject.MaxLabMark = *input.MaxLabMark
	}
	if input.MaxSemesterMark != nil {
		subject.MaxSemesterMark = *input.MaxSemesterMark
	}
	if input.MaxFinalExam != nil {
		subject.MaxFinalExam = *input.MaxFinalExam
	}
	if input.Credits != nil {
		subject.Credits = *input.Credits
	}
	if input.Active != nil {
		subject.Active = *input.Active
	}
	if input.Ministerial != nil {
		subject.Ministerial = *input.Ministerial
	}

	v := validator.New()
	if data.ValidateSubject(v, subject); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Subjects.Update(subject)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Write the updated movie record in a JSON response.
	header := make(http.Header)
	header.Set("Access-Control-Allow-Origin", "*")                                       // Allow the Vue.js app to make requests
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS") // Allow these HTTP methods
	header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")            // Allow these headers
	header.Set("Access-Control-Allow-Credentials", "true")
	err = app.writeJSON(w, http.StatusOK, envelope{"subject": subject}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteSubject(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the movie from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Subjects.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	header := make(http.Header)
	header.Set("Access-Control-Allow-Origin", "*")                                       // Allow the Vue.js app to make requests
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS") // Allow these HTTP methods
	header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")            // Allow these headers
	header.Set("Access-Control-Allow-Credentials", "true")
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "subject successfully deleted"}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
