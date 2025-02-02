package main

import (
	"context"
	"net/http"
)

type contextKey string

const isLoggedInContextKey = contextKey("isLoggedIn")
const userModelContextKey = contextKey("userStruct")
const yearContextKey = contextKey("year")
const stageContextKey = contextKey("stage")

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

func (app *application) subjectsAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		tableName := "subjects" + year
		user, err := app.getUserFromContext(r)
		privilege, err := app.models.Privileges.CheckAccess(int(user.ID), tableName, stage)
		if err != nil {
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
