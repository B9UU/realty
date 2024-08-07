package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// costume type for context key to avoide collisions
type contextKey string

const reqIdKey = contextKey("requestIdKey")

// wrapping http.ResponseWriter to capture the request status code with custom WriteHeader method
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// logging requests
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// new responseWriter wrapper
		rww := &statusRecorder{w, 1}

		reqId := uuid.New().String()
		ctx := context.WithValue(r.Context(), reqIdKey, reqId)

		r = r.WithContext(ctx)

		app.LogNewRequest(r, start)
		next.ServeHTTP(rww, r)
		app.LogEndRequest(r, start, fmt.Sprintf("%d", rww.status))

	})
}

// Log new request
func (app *application) LogNewRequest(r *http.Request, start time.Time) {
	uui := r.Context().Value(reqIdKey).(string)
	lo := map[string]string{
		"ID":     uui,
		"method": r.Method,
		"URI":    r.RequestURI,
	}
	app.logger.PrintInfo("New request", lo)
}

// log the end of the request
func (app *application) LogEndRequest(r *http.Request, start time.Time, status string) {
	lo := map[string]string{
		"ID":       r.Context().Value(reqIdKey).(string),
		"method":   r.Method,
		"URI":      r.RequestURI,
		"status":   status,
		"duration": fmt.Sprintf("%d ms", time.Since(start).Milliseconds()),
	}
	app.logger.PrintInfo("Request done", lo)
}
