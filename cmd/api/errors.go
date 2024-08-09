package main

import (
	"fmt"
	"net/http"
)

// logs the error with r method and url
func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

// writes message and status to w.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request,
	message interface{}, status int) {
	err := app.writeJSON(w, status, envelope{"error": message}, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// writes 500 status code with helpful message and logs the error
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	app.logError(r, err)
	app.errorResponse(w, r, message, http.StatusInternalServerError)
}

// writes 404 with helpful message
func (app *application) notFoundErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "The request resource could not be found"
	app.errorResponse(w, r, message, http.StatusNotFound)
}

// writes 429 with helpful message
func (app *application) rateLimiterExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, message, http.StatusTooManyRequests)
}

// writes 401 with helpful message
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid credentials"
	app.errorResponse(w, r, message, http.StatusUnauthorized)
}

// writes 401 with helpful message
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, message, http.StatusUnauthorized)
}

// writes 405 with helpful message
func (app *application) methodNotAllowedErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, message, http.StatusMethodNotAllowed)
}

// writes 400 with helpful message
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, err.Error(), http.StatusBadRequest)
}

// writes 400 with helpful message
func (app *application) failedValidationRespone(w http.ResponseWriter, r *http.Request,
	errors map[string]string) {
	app.errorResponse(w, r, errors, http.StatusBadRequest)
}
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, message, http.StatusUnauthorized)
}
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(w, r, message, http.StatusForbidden)
}
func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to acces this resource"
	app.errorResponse(w, r, message, http.StatusForbidden)
}
