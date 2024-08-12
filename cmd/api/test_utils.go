package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/b9uu/realty/internal/data"
	"github.com/b9uu/realty/internal/mocks"
	"github.com/b9uu/realty/jsonlog"
)

// returns new test application
func newTestApp(rmock []*data.RealtyResponse) *application {
	return &application{
		models: data.Models{Realty: mocks.RealtyModelM{MockRealtyData: rmock}},
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}
}

// returns new test server
func newTestServer(t *testing.T, h http.Handler) *httptest.Server {
	ts := httptest.NewServer(h)
	return ts

}
