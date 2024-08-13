package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/b9uu/realty/internal/data"
	"github.com/b9uu/realty/internal/validator"
)

// TODO: only admin can access this handler
func (app *application) addRealty(w http.ResponseWriter, r *http.Request) {
	var realty data.RealtyInput
	v := validator.New()
	err := app.ReadJSON(w, r, &realty)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// insert to db
	err = app.models.Realty.Insert(&realty)
	if err != nil {
		// check if id is unique
		if errors.Is(err, data.ErrDuplicateId) {
			v.AddError("id", "a listing with this id already exists")
			app.failedValidationRespone(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("v1/realty/%d", realty.ID))
	err = app.writeJSON(w, http.StatusCreated, data.Envelope{"realty": realty}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getRealties(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city != "" {
		v := validator.New()
		if data.ValidateCity(v, city); !v.Valid() {
			app.failedValidationRespone(w, r, v.Errors)
			return
		}
	}
	app.logger.PrintInfo(city, nil)
	realties, err := app.models.Realty.GetAll(city)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.notFoundErrorResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	err = app.writeJSON(w, http.StatusOK, data.Envelope{"realties": realties}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) autoComplete(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	v := validator.New()
	if data.ValidateQuery(v, q); !v.Valid() {
		app.failedValidationRespone(w, r, v.Errors)
		return
	}
	results, err := app.models.Realty.AutoComplete(q)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.notFoundErrorResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	err = app.writeJSON(w, http.StatusOK, data.Envelope{"result": results}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
