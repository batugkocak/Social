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

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)

	// TODO: Revert back to auth userID from ctx
	var payload FollowUser
	readJSONErr := readJSON(w, r, &payload)
	if readJSONErr != nil {
		app.badRequestError(w, r, readJSONErr)
	}

	ctx := r.Context()
	followError := app.store.Followers.Follow(ctx, followerUser.ID, payload.UserID)
	if followError != nil {
		switch followError {
		case store.ErrConflict:
			app.conflictError(w, r, followError)
			return
		default:
			app.internalServerError(w, r, followError)
		}
		return
	}

	writeJsonErr := app.jsonResponse(w, http.StatusNoContent, nil)
	if writeJsonErr != nil {
		app.internalServerError(w, r, writeJsonErr)
		return
	}
}
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := getUserFromCtx(r)

	// TODO: Revert back to auth userID from ctx
	var payload FollowUser
	readJSONErr := readJSON(w, r, &payload)
	if readJSONErr != nil {
		app.badRequestError(w, r, readJSONErr)
	}

	ctx := r.Context()
	followError := app.store.Followers.Unfollow(ctx, unfollowedUser.ID, payload.UserID)
	if followError != nil {
		app.internalServerError(w, r, followError)
		return
	}

	writeJsonErr := app.jsonResponse(w, http.StatusNoContent, nil)
	if writeJsonErr != nil {
		app.internalServerError(w, r, writeJsonErr)
		return
	}
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
