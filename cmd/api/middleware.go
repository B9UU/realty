package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/google/uuid"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

// costume type for context key to avoide collisions
type contextKey string

const reqIdKey = contextKey("requestIdKey")

// wrapping http.ResponseWriter to capture the request status code with custom WriteHeader method
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func (app *application) rateLimiter(next http.Handler) http.Handler {
	type newClient struct {
		Limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*newClient)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {

			ip := realip.FromRequest(r)

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &newClient{
					Limiter: rate.NewLimiter(
						rate.Limit(app.config.limiter.rps),
						app.config.limiter.burst,
					)}
			}

			clients[ip].lastSeen = time.Now()
			if !clients[ip].Limiter.Allow() {
				mu.Unlock()
				app.rateLimiterExceededResponse(w, r)
				return
			}
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
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

// metrics
func (app *application) metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponseSent := expvar.NewInt("total_response_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_micro_seconds")
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totalRequestsReceived.Add(1)

		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalResponseSent.Add(1)
		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())
		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}
