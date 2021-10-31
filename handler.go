package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"encoding/json"
)

// AppHandler is a type that defines a standard handler with return type AppError
type AppHandler func(http.ResponseWriter, *http.Request) *AppError

// AppError is the type of error that all handlers will return
// This error will be used in the ServeHTTP function
// Good to set this equal to the actix way of handling errors
type AppError struct {
	Error   error
	Reason  string
	Message string
	Code    int
}

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Errorf("Handler error: status code: %d, error type: %s,\n\t message: %s,\n\t underlying err: %#v",
			e.Code, e.Reason, e.Message, e.Error)

		errorResponse := struct {
			Message string `json:"message"`
			Reason  string `json:"reason"`
		}{
			Message: e.Message,
			Reason:  e.Reason,
		}

		// The error is returned in a json, let the client know
		w.Header().Set("Content-Type", "application/json")
		// Set status code
		w.WriteHeader(e.Code)

		// Return a JSON object
		err := json.NewEncoder(w).Encode(errorResponse)
		if err != nil {
			http.Error(w, "Unable to marshal error object, status code still valid", e.Code)
		}
		return
	}

	// This might cause a superfluous writeHeader call
	if r.Method == "POST" {
		w.WriteHeader(201)
	}
}

// AppErrorf returns a type AppError
func AppErrorf(errorCode int, err error, reason string, format string, v ...interface{}) *AppError {
	return &AppError{
		Error:   err,
		Reason:  reason,
		Message: fmt.Sprintf(format, v...),
		Code:    errorCode,
	}
}
