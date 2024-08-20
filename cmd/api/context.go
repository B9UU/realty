package main

import (
	"context"
	"net/http"

	"github.com/b9uu/realty/internal/data"
)

// costume type for context key to avoide collisions
type contextKey string

var userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)

	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("something went wrong with context value")
	}
	return user
}
