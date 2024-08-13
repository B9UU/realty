package mocks

import (
	"database/sql"
	"strings"
	"time"

	"github.com/b9uu/realty/internal/data"
)

// RealtyModel Mock that implements Realty Interface
type RealtyModelM struct {
	MockRealtyData []*data.Realties
	MockCities     []string
}

func (m RealtyModelM) Insert(realty *data.RealtyInput) error {
	realty.ID = 1
	return nil
}

func (m RealtyModelM) GetAll(city string, filters data.Filters) ([]*data.Realties, data.Metadata, error) {
	return m.MockRealtyData, data.Metadata{}, nil

}

// mock for autocomplete method
func (m RealtyModelM) AutoComplete(sub string) ([]string, error) {

	var results []string
	for _, s := range m.MockCities {
		if len(results) <= 10 {
			if strings.Contains(s, sub) {
				results = append(results, s)
			}
		} else {
			break
		}
	}
	return results, nil
}

var MockCities = []string{
	"Toronto",
	"Montreal",
	"Vancouver",
	"Calgary",
	"Edmonton",
	"Ottawa",
	"Winnipeg",
	"Quebec City",
	"Hamilton",
	"Kitchener",
	"London",
	"Victoria",
	"Halifax",
	"Oshawa",
	"Windsor",
	"Saskatoon",
	"St. Catharines",
	"Regina",
	"St. John's",
	"Barrie",
	"Kelowna",
	"Abbotsford",
	"Greater Sudbury",
	"Kingston",
	"Saguenay",
	"Trois-Rivières",
	"Guelph",
	"Moncton",
	"Brantford",
	"Thunder Bay",
}

var MockRealties = []data.Realties{
	{
		ID:           1,
		Name:         "Modern Apartment",
		Address1:     "123 Main St",
		Address2:     "Apt 4B",
		PostalCode:   "12345",
		CityName:     "San Francisco",
		PropertyType: "Residential",
		Updated:      time.Now(),
	},
	{
		ID:           2,
		Name:         "Cozy Cottage",
		Address1:     "456 Elm St",
		Address2:     "",
		PostalCode:   "67890",
		CityName:     "New York",
		PropertyType: "Residential",
		Updated:      time.Now(),
	},
}
var MockRealtiesResponse = []data.RealtyResponse{
	{
		ID:              1,
		Name:            "Modern Apartment",
		Address1:        "123 Main St",
		Address2:        "Apt 4B",
		PostalCode:      "12345",
		Lat:             37.7749,
		Lng:             -122.4194,
		Title:           "Beautiful Modern Apartment in Downtown",
		FeaturedStatus:  "Featured",
		CityName:        "San Francisco",
		PhotoCount:      10,
		PhotoURL:        "http://example.com/photo.jpg",
		RawPropertyType: "Apartment",
		PropertyType:    "Residential",
		Updated:         time.Now(),
		RentRange: []sql.NullInt32{
			{Int32: 1500, Valid: true},
			{Int32: 1800, Valid: true},
		},
		BedsRange: []sql.NullInt32{
			{Int32: 2, Valid: true},
			{Int32: 3, Valid: true},
		},
		BathsRange: []sql.NullInt32{
			{Int32: 1, Valid: true},
			{Int32: 2, Valid: true},
		},
		DimensionsRange: []sql.NullInt32{
			{Int32: 850, Valid: true},
			{Int32: 1050, Valid: true},
		},
	},
	{
		ID:              2,
		Name:            "Cozy Cottage",
		Address1:        "456 Elm St",
		Address2:        "",
		PostalCode:      "67890",
		Lat:             40.7128,
		Lng:             -74.0060,
		Title:           "Charming Cottage in the Suburbs",
		FeaturedStatus:  "Not Featured",
		CityName:        "New York",
		PhotoCount:      5,
		PhotoURL:        "http://example.com/cottage.jpg",
		RawPropertyType: "Cottage",
		PropertyType:    "Residential",
		Updated:         time.Now(),
		RentRange: []sql.NullInt32{
			{Int32: 1200, Valid: true},
			{Int32: 1500, Valid: true},
		},
		BedsRange: []sql.NullInt32{
			{Int32: 3, Valid: true},
			{Int32: 5, Valid: true},
		},
		BathsRange: []sql.NullInt32{
			{Int32: 2, Valid: true},
			{Int32: 3, Valid: true},
		},
		DimensionsRange: []sql.NullInt32{
			{Int32: 1000, Valid: true},
			{Int32: 1200, Valid: true},
		},
	},
}
