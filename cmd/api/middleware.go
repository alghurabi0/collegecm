package main

import (
	"context"
	"net/http"
	"strings"

	"collegecm.hamid.net/internal/data"
)

type contextKey string

const isLoggedInContextKey = contextKey("isLoggedIn")
const userModelContextKey = contextKey("userStruct")
const yearContextKey = contextKey("year")
const stageContextKey = contextKey("stage")
const idContextKey = contextKey("id")
const customPrivsContextKey = contextKey("custom_privs")
const studentContextKey = contextKey("student")
const subjectContextKey = contextKey("subject")
const markContextKey = contextKey("mark")

//const stagesContextKey = contextKey("stages")

func (app *application) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := map[string]bool{
			"http://localhost:5173":            true,
			"https://collegecm-vue.vercel.app": true,
		}
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

/* middleware to check if user is authenticated, if not return unauthorized
 * else save user struct in context */
func (app *application) isLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := app.sessionManager.GetInt(r.Context(), "userID")
		if userId == 0 {
			app.unauthorized(w, r)
			return
		}
		user, err := app.models.Users.Get(int64(userId))
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), isLoggedInContextKey, true)
		ctx = context.WithValue(ctx, userModelContextKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (app *application) getAllAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		parts := strings.Split(path, "/")
		cat := parts[2]
		year, err := app.readYearParam(r)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		stage, err := app.readStageParam(r)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		tableName := cat + "_" + year
		user, err := app.getUserFromContext(r)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		privilege, err := app.models.Privileges.CheckAccess(int(user.ID), tableName, stage)
		if err != nil {
			if err == data.ErrRecordNotFound {
				app.unauthorized(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}
		if !privilege.CanRead {
			app.unauthorized(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), yearContextKey, year)
		ctx = context.WithValue(ctx, stageContextKey, stage)
		//ctx = context.WithValue(ctx, stagesContextKey, stages)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// func (app *application) createAccess(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		path := r.URL.Path
// 		parts := strings.Split(path, "/")
// 		cat := parts[2]
// 		year, err := app.readYearParam(r)
// 		if err != nil {
// 			app.notFoundResponse(w, r)
// 			return
// 		}
// 		tableName := cat + "_" + year
// 		user, err := app.getUserFromContext(r)
// 		if err != nil {
// 			app.serverErrorResponse(w, r, err)
// 			return
// 		}
// 		hasAccess, err := app.models.Privileges.CheckCreateAccess(int(user.ID), tableName)
// 		if err != nil {
// 			app.serverErrorResponse(w, r, err)
// 			return
// 		}
// 		if !hasAccess {
// 			app.unauthorized(w, r)
// 			return
// 		}
// 		ctx := context.WithValue(r.Context(), yearContextKey, year)
// 		r = r.WithContext(ctx)
// 		next.ServeHTTP(w, r)
// 	})
// }

func (app *application) writeAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := app.getUserFromContext(r)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		year, err := app.readYearParam(r)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		id, err := app.readIdParam(r)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), yearContextKey, year)
		ctx = context.WithValue(ctx, idContextKey, id)
		path := r.URL.Path
		parts := strings.Split(path, "/")
		cat := parts[2]
		tableName := cat + "_" + year

		var student *data.Student
		var subject *data.Subject
		var mark *data.Mark
		var stage string
		switch cat {
		case "students":
			student, err = app.models.Students.Get(year, id)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			stage = student.Stage
			ctx = context.WithValue(ctx, studentContextKey, student)
		case "subjects":
			subject, err = app.models.Subjects.Get(year, id)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			stage = subject.Stage
			ctx = context.WithValue(ctx, subjectContextKey, subject)
		case "carryovers":
			stage, err = app.models.Carryovers.GetStage(id, year)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			if stage == "" {
				app.serverErrorResponse(w, r, err)
				return
			}
		case "exempteds":
			tableName = "exempted_" + year
			stage, err = app.models.Exempteds.GetStage(id, year)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			if stage == "" {
				app.serverErrorResponse(w, r, err)
				return
			}
		case "marks":
			mark, err = app.models.Marks.GetRaw(year, id)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			ctx = context.WithValue(ctx, markContextKey, mark)
			stage, err = app.models.Marks.GetStage(mark.StudentId, year)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			if stage == "" {
				app.serverErrorResponse(w, r, err)
				return
			}
		default:
			app.notFoundResponse(w, r)
			return
		}

		hasAccess, err := app.models.Privileges.CheckWriteAccess(int(user.ID), tableName, stage)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if !hasAccess {
			app.unauthorized(w, r)
			return
		}
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (app *application) customAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		year, err := app.readYearParam(r)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		id, err := app.readIdParam(r)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		tableName := "students_" + year
		student, err := app.models.Students.GetCustom(tableName, id)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		stage := student.Stage
		user, err := app.getUserFromContext(r)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		privileges, err := app.models.Privileges.CheckCustomAccess(int(user.ID), year, stage)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), yearContextKey, year)
		ctx = context.WithValue(ctx, idContextKey, id)
		ctx = context.WithValue(ctx, customPrivsContextKey, privileges)
		ctx = context.WithValue(ctx, studentContextKey, student)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (app *application) userReadAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := app.getUserFromContext(r)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		path := r.URL.Path
		parts := strings.Split(path, "/")
		table := parts[2]
		hasAccess, err := app.models.Privileges.CheckUserReadAccess(int(user.ID), table)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if !hasAccess {
			app.unauthorized(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) userWriteAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := app.getUserFromContext(r)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		path := r.URL.Path
		parts := strings.Split(path, "/")
		table := parts[2]
		hasAccess, err := app.models.Privileges.CheckUserWriteAccess(int(user.ID), table)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if !hasAccess {
			app.unauthorized(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
