package data

import (
	"database/sql"
	"time"
)

// cintains all models the application needs
type Models struct {
	Realty RealtyInterface
}
type RealtyResponse struct {
	ID              int64           `json:"id"`
	Name            string          `json:"name"`
	Address1        string          `json:"address1"`
	Address2        string          `json:"address2"`
	PostalCode      string          `json:"postal_code"`
	Lat             float64         `json:"lat"`
	Lng             float64         `json:"lng"`
	Title           string          `json:"title"`
	FeaturedStatus  string          `json:"featured_status"`
	CityName        string          `json:"city_name"`
	PhotoCount      int             `json:"photo_count"`
	PhotoURL        string          `json:"photo_url"`
	RawPropertyType string          `json:"raw_property_type"`
	PropertyType    string          `json:"property_type"`
	Updated         time.Time       `json:"updated"`
	RentRange       []sql.NullInt32 `json:"rent_range"`
	BedsRange       []sql.NullInt32 `json:"beds_range"`
	BathsRange      []sql.NullInt32 `json:"baths_range"`
	DimensionsRange []sql.NullInt32 `json:"dimensions_range"`
}
type RealtyInput struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Address1        string    `json:"address1"`
	Address2        string    `json:"address2"`
	PostalCode      string    `json:"postal_code"`
	Lat             float64   `json:"lat"`
	Lng             float64   `json:"lng"`
	Title           string    `json:"title"`
	FeaturedStatus  string    `json:"featured_status"`
	CityName        string    `json:"city_name"`
	PhotoCount      int       `json:"photo_count"`
	PhotoURL        string    `json:"photo_url"`
	RawPropertyType string    `json:"raw_property_type"`
	PropertyType    string    `json:"property_type"`
	Updated         time.Time `json:"updated"`
	RentRange       []int     `json:"rent_range"`
	BedsRange       []int     `json:"beds_range"`
	BathsRange      []int     `json:"baths_range"`
	DimensionsRange []int     `json:"dimensions_range"`
}

// initiate new models
func NewModels(db *sql.DB) Models {
	return Models{
		Realty: RealtyModel{DB: db},
	}
}
