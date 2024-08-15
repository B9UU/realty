package data

import (
	"database/sql"
	"time"
)

// type to wrape the response
type Envelope map[string]interface{}

// cintains all models the application needs
type Models struct {
	Realty RealtyInterface
	User   UserInterface
	Token  TokenInterface
}

// initiate new models
func NewModels(db *sql.DB) Models {
	return Models{
		Realty: RealtyModel{DB: db},
		User:   UserModel{DB: db},
		Token:  TokenModel{DB: db},
	}
}

type Realty struct {
	ID              int64     `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Address1        string    `json:"address1,omitempty"`
	Address2        string    `json:"address2,omitempty"`
	PostalCode      string    `json:"postal_code,omitempty"`
	Lat             float64   `json:"lat,omitempty"`
	Lng             float64   `json:"lng,omitempty"`
	Title           string    `json:"title,omitempty"`
	FeaturedStatus  string    `json:"featured_status,omitempty"`
	CityName        string    `json:"city_name,omitempty"`
	PhotoCount      int       `json:"photo_count,omitempty"`
	PhotoURL        string    `json:"photo_url,omitempty"`
	RawPropertyType string    `json:"raw_property_type,omitempty"`
	PropertyType    string    `json:"property_type,omitempty"`
	Updated         time.Time `json:"updated,omitempty"`
	RentRange       []int32   `json:"rent_range,omitempty"`
	BedsRange       []int32   `json:"beds_range,omitempty"`
	BathsRange      []int32   `json:"baths_range,omitempty"`
	DimensionsRange []int32   `json:"dimensions_range,omitempty"`
}

type Realties struct {
	ID           int64     `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Address1     string    `json:"address1,omitempty"`
	Address2     string    `json:"address2,omitempty"`
	PostalCode   string    `json:"postal_code,omitempty"`
	CityName     string    `json:"city_name,omitempty"`
	PropertyType string    `json:"property_type,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
}
