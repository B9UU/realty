package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// handle routes
func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.Handler(http.MethodGet, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	router.HandlerFunc(http.MethodPost, "/realty", app.addRealty)
	router.HandlerFunc(http.MethodGet, "/realty/:id", app.Realty)
	router.HandlerFunc(http.MethodGet, "/realties", app.Realties)
	router.HandlerFunc(http.MethodGet, "/auto-complete", app.autoComplete)

	router.HandlerFunc(http.MethodPut, "/users/activated", app.activateUser)

	router.HandlerFunc(http.MethodPost, "/usersActivated", app.registerUserActivated)
	router.HandlerFunc(http.MethodPost, "/users", app.registerUser)
	router.HandlerFunc(http.MethodPost, "/login", app.AuthToken)
	return app.logRequest(app.rateLimiter(router))
}
