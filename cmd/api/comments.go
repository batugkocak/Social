package main

import (
	"net/http"

	"github.com/batugkocak/social/internal/store"
)

// TODO: Validate if the post exists
type CreateCommentPayload struct {
	PostID  int64  `json:"post_id" validate:"required"`
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,max=100"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
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
		PostID:  payload.PostID,
		UserID:  payload.UserID,
		Content: payload.Content,
	}
	ctx := r.Context()

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
