package main

import (
	"context"
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

// CreatePost godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
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
		UserID:    int64(mockUserId),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	//TODO: Add Validation (Add to DB as varchar(x) too for the title and content)
	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// GetPost godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	store.Post
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	post.Comments = *comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// DeletePost godoc
//
//	@Summary		Deletes a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object} string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postIDString := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(postIDString, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	deleteErr := app.store.Posts.DeleteById(ctx, postID)
	if deleteErr != nil {
		fmt.Print(deleteErr)
		switch {
		case errors.Is(deleteErr, store.ErrNotFound):
			app.notFoundError(w, r, deleteErr)
			return
		default:
			app.internalServerError(w, r, deleteErr)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, "Post deleted!"); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type UpdatePostPayload struct {
	Content string   `json:"content" validate:"required,max=100"`
	Title   string   `json:"title" validate:"required,max=100"`
	Tags    []string `json:"tags"`
}

// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (app *application) patchPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(&payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.Content != "" {
		post.Content = payload.Content
	}
	if payload.Title != "" {
		post.Title = payload.Title
	}
	if len(payload.Tags) > 0 {
		post.Tags = payload.Tags
	}

	ctx := r.Context()

	if err := app.store.Posts.UpdateById(ctx, post); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// Middlewares & Helpers

type postKey string

const postCtx postKey = "post"

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
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

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
