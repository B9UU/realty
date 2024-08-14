package mocks

import (
	"strings"
	"time"

	"github.com/b9uu/realty/internal/data"
)

// RealtyModel Mock that implements Realty Interface
type RealtyModelM struct {
	MockRealtyData []*data.Realties
	MockCities     []string
}

func (m RealtyModelM) Insert(realty *data.Realty) error {
	realty.ID = 1
	return nil
}

func (m RealtyModelM) GetAll(city string, filters data.Filters) ([]*data.Realties, data.Metadata, error) {
	return m.MockRealtyData, data.Metadata{}, nil

}
func (m RealtyModelM) Get(id int64) (*data.Realty, error) {
	return &data.Realty{}, nil
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
var MockRealtiesResponse = []data.Realty{
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
		RentRange: []int32{
			1500,
			1800,
		},
		BedsRange: []int32{
			3, 2,
		},
		BathsRange: []int32{
			1, 2,
		},
		DimensionsRange: []int32{
			850,
			1050,
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
		RentRange: []int32{
			1200,
			1500,
		},
		BedsRange: []int32{
			3,
			5,
		},
		BathsRange: []int32{
			2,
			3,
		},
		DimensionsRange: []int32{
			1000,
			1200,
		},
	},
}
