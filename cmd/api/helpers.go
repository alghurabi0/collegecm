package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"collegecm.hamid.net/internal/data"
	"github.com/gocarina/gocsv"
	"github.com/xuri/excelize/v2"
)

type envelope map[string]interface{}

// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.
// func (app *application) readIdYearParam(r *http.Request) (string, int64, error) {
// 	year := r.PathValue("year")
// 	if strings.TrimSpace(year) == "" {
// 		return "", -1, errors.New("empty year parameter")
// 	}
// 	params := r.PathValue("id")
// 	id, err := strconv.ParseInt(params, 10, 64)
// 	if err != nil || id < 0 {
// 		return "", -1, errors.New("invalid id parameter")
// 	}
// 	return year, id, nil
// }

func (app *application) readIdParam(r *http.Request) (int64, error) {
	params := r.PathValue("id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil || id < 0 {
		return -1, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) readYearParam(r *http.Request) (string, error) {
	year := r.PathValue("year")
	if strings.TrimSpace(year) == "" {
		return "", errors.New("empty year parameter")
	}
	return year, nil
}

func (app *application) readStageParam(r *http.Request) (string, error) {
	param2 := r.PathValue("stage")
	switch strings.TrimSpace(param2) {
	case "1":
		param2 = "الاولى"
	case "2":
		param2 = "الثانية"
	case "3":
		param2 = "الثالثة"
	case "4":
		param2 = "الرابعة"
	case "5":
		param2 = "الخامسة"
	case "6":
		param2 = "السادسة"
	case "all":
		param2 = "all"
	default:
		param2 = ""
	}
	if strings.TrimSpace(param2) == "" {
		return "", errors.New("empty stage parameter")
	}
	return param2, nil
}

// func (app *application) readParams(r *http.Request) (string, string, error) {
// 	param1 := r.PathValue("year")
// 	if strings.TrimSpace(param1) == "" {
// 		return "", "", errors.New("empty year parameter")
// 	}
// 	param2 := r.PathValue("stage")
// 	switch strings.TrimSpace(param2) {
// 	case "1":
// 		param2 = "الاولى"
// 	case "2":
// 		param2 = "الثانية"
// 	case "3":
// 		param2 = "الثالثة"
// 	case "4":
// 		param2 = "الرابعة"
// 	case "5":
// 		param2 = "الخامسة"
// 	case "6":
// 		param2 = "السادسة"
// 	case "all":
// 		param2 = "all"
// 	default:
// 		param2 = ""
// 	}
// 	if strings.TrimSpace(param2) == "" {
// 		return "", "", errors.New("empty stage parameter")
// 	}
// 	return param1, param2, nil
// }

// Define a writeJSON() helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')
	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	// Note that it's OK if the provided header map is nil. Go doesn't throw an error
	// if you try to range over (or generally, read from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}
	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	// Decode the request body into the target destination.
	// Decode the request body to the destination.
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message. Note that there's an open
		// issue at https://github.com/golang/go/issues/29035 regarding turning this
		// into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		// If the request body exceeds 1MB in size the decode will now fail with the
		// error "http: request body too large". There is an open issue about turning
		// this into a distinct error type at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) processFile(file *multipart.File, data interface{}) error {
	err := gocsv.UnmarshalMultipartFile(file, data)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) saveFile(file multipart.File, filePath string) error {
	// Create a new file on the server
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the uploaded file's content to the new file
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}

	return nil
}

// remove a file
func (app *application) removeFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) readExcel(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1") // Assuming data is in Sheet1
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// func (app *application) isLoggedInCheck(r *http.Request) bool {
// 	isLoggedIn, ok := r.Context().Value(isLoggedInContextKey).(bool)
// 	if !ok {
// 		return false
// 	}
// 	return isLoggedIn
// }

func (app *application) getUserFromContext(r *http.Request) (*data.User, error) {
	user, ok := r.Context().Value(userModelContextKey).(*data.User)
	if !ok {
		return nil, errors.New("can't get user object from context")
	}
	return user, nil
}

func (app *application) getYearFromContext(r *http.Request) (string, error) {
	year, ok := r.Context().Value(yearContextKey).(string)
	if !ok {
		return "", errors.New("can't get year from context")
	}
	return year, nil
}

func (app *application) getStageFromContext(r *http.Request) (string, error) {
	stage, ok := r.Context().Value(stageContextKey).(string)
	if !ok {
		return "", errors.New("can't get stage from context")
	}
	return stage, nil
}

func (app *application) getIdFromContext(r *http.Request) (int64, error) {
	id, ok := r.Context().Value(idContextKey).(int64)
	if !ok {
		return -1, errors.New("can't get id from context")
	}
	return id, nil
}

func (app *application) getCustomPrivsFromContext(r *http.Request) (*data.CustomPrivilegeAccess, error) {
	customPrivilegeAccess, ok := r.Context().Value(customPrivsContextKey).(*data.CustomPrivilegeAccess)
	if !ok {
		return nil, errors.New("can't get custom privilege access from context")
	}
	return customPrivilegeAccess, nil
}

func (app *application) getStudentFromContext(r *http.Request) (*data.Student, error) {
	student, ok := r.Context().Value(studentContextKey).(*data.Student)
	if !ok {
		return nil, errors.New("can't get student from context")
	}
	return student, nil
}

func (app *application) getSubjectFromContext(r *http.Request) (*data.Subject, error) {
	subject, ok := r.Context().Value(subjectContextKey).(*data.Subject)
	if !ok {
		return nil, errors.New("can't get subject from context")
	}
	return subject, nil
}

func (app *application) getMarkFromContext(r *http.Request) (*data.Mark, error) {
	mark, ok := r.Context().Value(markContextKey).(*data.Mark)
	if !ok {
		return nil, errors.New("can't get mark from context")
	}
	return mark, nil
}

// get stages from context, an array of strings
//func (app *application) getStagesFromContext(r *http.Request) ([]string, error) {
//	stages, ok := r.Context().Value(stagesContextKey).([]string)
//	if !ok {
//		return nil, errors.New("can't get stages from context")
//	}
//	return stages, nil
//}
