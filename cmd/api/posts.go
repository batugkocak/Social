package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/batugkocak/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Content string   `json:"content" validate:"required,max=100"`
	Title   string   `json:"title" validate:"required,max=100"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	mockUserId := 1 // Fake User

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		//TODO: Change after auth
		UserId:    int64(mockUserId),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	//TODO: Add Validation (Add to DB as varchar(x) too for the title and content)
	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	fmt.Println(postID)
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	post, err := app.store.Posts.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
