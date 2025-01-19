package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal Error", "method", r.Method, "Path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusInternalServerError, "The Server Encountered a Problem.")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Bad Request Error", "method", r.Method, "Path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Not Found Error", "method", r.Method, "Path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusNotFound, "Requested entity not found.")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Conflict Error", "method", r.Method, "Path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusConflict, err.Error())
}
