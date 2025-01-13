package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: Pagination, filters and sorting etc.

	// Hard coded id until auth
	feed, feedErr := app.store.Posts.GetUserFeed(r.Context(), int64(42))
	if feedErr != nil {
		app.internalServerError(w, r, feedErr)
	}

	responseErr := app.jsonResponse(w, http.StatusOK, feed)
	if responseErr != nil {
		app.internalServerError(w, r, responseErr)
	}
}
