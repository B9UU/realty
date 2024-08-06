package main

import (
	"fmt"
	"net/http"

	"github.com/b9uu/realty/internal/data"
)

func (app *application) addRealty(w http.ResponseWriter, r *http.Request) {
	var realty struct {
		ListingType    string `json:"listing_type"`
		PromoType      string `json:"promo_type,omitempty"`
		URL            string `json:"url"`
		ProjectName    string `json:"project_name"`
		DisplayAddress string `json:"display_address"`
	}
	err := app.ReadJSON(w, r, &realty)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	newRealty := &data.Realty{
		ListingType:    realty.ListingType,
		PromoType:      realty.PromoType,
		URL:            realty.URL,
		ProjectName:    realty.ProjectName,
		DisplayAddress: realty.DisplayAddress,
	}
	// insert to db
	err = app.models.Realty.Insert(newRealty)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("v1/realty/%d", newRealty.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"realty": newRealty}, header)
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
