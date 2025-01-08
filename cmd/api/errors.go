package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal Server Error: %s Path: %s Error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusInternalServerError, "The Server Encountered a Problem.")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad Request Error: %s Path: %s Error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not Found Error: %s Path: %s Error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusNotFound, "Requested entity not found.")
}
