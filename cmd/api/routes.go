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
	router.HandlerFunc(http.MethodGet, "/realty", app.getRealties)
	return app.logRequest(router)
}
