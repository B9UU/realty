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
	var realty data.Realty
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
		switch {
		case errors.Is(err, data.ErrDuplicateId):
			v.AddError("id", "a listing with this id already exists")
			app.failedValidationRespone(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
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

func (app *application) Realties(w http.ResponseWriter, r *http.Request) {

	var realty struct {
		City string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	realty.City = app.readString(qs, "city", "")
	realty.Filters.Page = app.readInt(qs, "page", 1, v)
	realty.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	realty.Filters.Sort = app.readString(qs, "sort", "id")
	realty.Filters.SortSafeList = []string{
		"id", "updated",
		"-id", "-updated",
	}
	if data.ValidateFilters(v, realty.Filters); !v.Valid() {
		app.failedValidationRespone(w, r, v.Errors)
		return

	}
	realties, metadata, err := app.models.Realty.GetAll(realty.City, realty.Filters)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, data.Envelope{"realties": realties, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) Realty(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundErrorResponse(w, r)
		return
	}
	fmt.Println(id)
	realty, err := app.models.Realty.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, data.Envelope{"realty": realty}, nil)
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
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, data.Envelope{"result": results}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
