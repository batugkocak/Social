package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/batugkocak/social/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	writeJsonErr := app.jsonResponse(w, http.StatusOK, user)
	if writeJsonErr != nil {
		app.internalServerError(w, r, writeJsonErr)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
}
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
}

// Middlewares & Helpers

type userKey string

const userCtx userKey = "user"

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
