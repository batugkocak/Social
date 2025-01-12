package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/batugkocak/social/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, paramErr := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if paramErr != nil {
		app.badRequestError(w, r, paramErr)
		return
	}

	ctx := r.Context()
	user, err := app.store.Users.GetById(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestError(w, r, err)
			return
		default:
			app.internalServerError(w, r, paramErr)
			return
		}
	}

	writeJsonErr := app.jsonResponse(w, http.StatusOK, user)
	if writeJsonErr != nil {
		app.internalServerError(w, r, writeJsonErr)
		return
	}
}
