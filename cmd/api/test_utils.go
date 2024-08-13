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

func toPointer(drr []data.Realties) []*data.Realties {

	list := []*data.Realties{}
	for _, rr := range drr {
		list = append(list, &rr)
	}
	return list
}

// returns new test application
func newTestApp() *application {
	return &application{
		models: data.Models{
			Realty: mocks.RealtyModelM{
				MockCities:     mocks.MockCities,
				MockRealtyData: toPointer(mocks.MockRealties),
			},
		},
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}
}

// returns new test server
func newTestServer(t *testing.T, h http.Handler) *httptest.Server {
	ts := httptest.NewServer(h)
	return ts
}
