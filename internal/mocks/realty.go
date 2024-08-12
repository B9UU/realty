package mocks

import (
	"database/sql"
	"strings"
	"time"

	"github.com/b9uu/realty/internal/data"
)

// RealtyModel Mock that implements Realty Interface
type RealtyModelM struct {
	MockRealtyData []*data.RealtyResponse
}

func (m RealtyModelM) Insert(realty *data.RealtyInput) error {
	realty.ID = 1
	return nil
}

func (m RealtyModelM) GetAll() ([]*data.RealtyResponse, error) {
	return m.MockRealtyData, nil

}

// mock for autocomplete method
func (m RealtyModelM) AutoComplete(sub string) ([]string, error) {
	cities := []string{
		"Vancouver",
		"Montreal",
	}

	var results []string
	for _, s := range cities {
		if strings.Contains(s, sub) {
			results = append(results, s)
		}
	}
	return results, nil
}

var MockRealties = []data.RealtyResponse{
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
		RentRange:       []sql.NullInt32{sql.NullInt32{Int32: 1500, Valid: true}},
		BedsRange:       []sql.NullInt32{sql.NullInt32{Int32: 2, Valid: true}},
		BathsRange:      []sql.NullInt32{sql.NullInt32{Int32: 1, Valid: true}},
		DimensionsRange: []sql.NullInt32{sql.NullInt32{Int32: 850, Valid: true}},
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
		RentRange:       []sql.NullInt32{sql.NullInt32{Int32: 1200, Valid: true}},
		BedsRange:       []sql.NullInt32{sql.NullInt32{Int32: 3, Valid: true}},
		BathsRange:      []sql.NullInt32{sql.NullInt32{Int32: 2, Valid: true}},
		DimensionsRange: []sql.NullInt32{sql.NullInt32{Int32: 1000, Valid: true}},
	},
	// Add more mock entries as needed
}
