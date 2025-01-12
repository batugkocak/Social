package main

import (
	"net/http"
	"strconv"

	"github.com/batugkocak/social/internal/store"
	"github.com/go-chi/chi/v5"
)

// TODO: Validate if the post exists
type CreateCommentPayload struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,max=100"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	postIDString := chi.URLParam(r, "postID")
	postID, parseErr := strconv.ParseInt(postIDString, 10, 64)
	if parseErr != nil {
		app.internalServerError(w, r, parseErr)
		return
	}

	ctx := r.Context()
	_, postErr := app.store.Posts.GetById(ctx, postID)
	if postErr != nil {
		app.badRequestError(w, r, store.ErrNotFound)
		return
	}

	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comment := &store.Comment{
		PostID:  postID,
		UserID:  payload.UserID,
		Content: payload.Content,
	}

	err := app.store.Comments.Create(ctx, comment)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	responseErr := app.jsonResponse(w, http.StatusCreated, comment)
	if responseErr != nil {
		app.internalServerError(w, r, responseErr)
		return
	}

}
