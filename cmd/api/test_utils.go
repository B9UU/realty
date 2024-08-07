package main

import (
	"github.com/b9uu/realty/internal/data"
	"github.com/b9uu/realty/internal/mocks"
)

func newTestApp() *application {
	return &application{
		models: data.Models{Realty: mocks.RealtyModelM{}},
	}

}
