package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/b9uu/realty/internal/data"
	"github.com/b9uu/realty/internal/validator"
)

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
	err = app.writeJSON(w, http.StatusCreated, envelope{"realty": realty}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getRealties(w http.ResponseWriter, r *http.Request) {
	realties, err := app.models.Realty.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"realties": realties}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
