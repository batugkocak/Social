package main

import (
	"net/http"

	"github.com/batugkocak/social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Hard coded id until auth
	feed, feedErr := app.store.Posts.GetUserFeed(r.Context(), int64(42), fq)
	if feedErr != nil {
		app.internalServerError(w, r, feedErr)
		return
	}

	responseErr := app.jsonResponse(w, http.StatusOK, feed)
	if responseErr != nil {
		app.internalServerError(w, r, responseErr)
		return
	}
}
