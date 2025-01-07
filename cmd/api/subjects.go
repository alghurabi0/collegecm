package main

import (
	"fmt"
	"net/http"
	"time"

	"collegecm.hamid.net/internal/data"
	"collegecm.hamid.net/internal/validator"
)

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
	// Otherwise, interpolate the movie ID in a placeholder response.
	fmt.Fprintf(w, "show the details of movie %d\n", id)
}

func (app *application) createSubjectHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Movie struct that we created earlier). This struct will be our *target
	// decode destination*.
	var input struct {
		ID                 int       `json:"subject_id"`
		SubjectName        string    `json:"subject_name"`
		SubjectNameEnglish string    `json:"subject_name_english"`
		Stage              string    `json:"stage"`
		Semester           string    `json:"semester"`
		Department         string    `json:"department"`
		MaxTheoryMark      int       `json:"max_theory_mark"`
		MaxLabMark         int       `json:"max_lab_mark"`
		MaxSemesterMark    int       `json:"max_semester_mark"`
		MaxFinalExam       int       `json:"max_final_exam"`
		Credits            int       `json:"credits"`
		Active             string    `json:"active"`
		Ministerial        string    `json:"ministerial"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updatedAt"`
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
	err = app.writeJSON(w, http.StatusCreated, envelope{"subject": subject}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
