package main

import "net/http"

// Shortened HTTP constants
const (
	NotFound            = http.StatusNotFound
	BadRequest          = http.StatusBadRequest
	OK                  = http.StatusOK
	InternalServerError = http.StatusInternalServerError
	NoContent           = http.StatusNoContent
	Created             = http.StatusCreated
)
