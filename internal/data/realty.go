package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var (
	ErrDuplicateId = errors.New("duplicate id")
)

type Realty struct {
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

type RealtyModel struct {
	DB *sql.DB
}

type RealtyInterface interface {
	Insert(realty *Realty) error
	GetAll() ([]*Realty, error)
}

func (m RealtyModel) Insert(realty *Realty) error {
	query := `
			INSERT INTO realty (
				id, name, address1, address2, postal_code, lat, lng, title,
				featured_status, city_name, photo_count, photo_url, raw_property_type,
				property_type, updated, rent_range, beds_range, baths_range, dimensions_range
			)
			VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
				$11, $12, $13, $14, $15, $16, $17, $18, $19
			)`

	args := []interface{}{
		realty.ID, realty.Name, realty.Address1, realty.Address2, realty.PostalCode,
		realty.Lat, realty.Lng, realty.Title, realty.FeaturedStatus, realty.CityName,
		realty.PhotoCount, realty.PhotoURL, realty.RawPropertyType, realty.PropertyType,
		realty.Updated, pq.Array(realty.RentRange), pq.Array(realty.BedsRange),
		pq.Array(realty.BathsRange), pq.Array(realty.DimensionsRange),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		// check if the id already exists
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("23505") {
			return ErrDuplicateId
		}
		return err
	}
	return nil
}
func (m RealtyModel) GetAll() ([]*Realty, error) {
	query := "SELECT * FROM realty LIMIT 200"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	realties := []*Realty{}
	for rows.Next() {
		var realty Realty
		err := rows.Scan(
			&realty.ID, &realty.Name, &realty.Address1, &realty.Address2, &realty.PostalCode,
			&realty.Lat, &realty.Lng, &realty.Title, &realty.FeaturedStatus, &realty.CityName,
			&realty.PhotoCount, &realty.PhotoURL, &realty.RawPropertyType, &realty.PropertyType,
			&realty.Updated, pq.Array(&realty.RentRange), pq.Array(&realty.BedsRange),
			pq.Array(&realty.BathsRange), pq.Array(&realty.DimensionsRange),
		)
		if err != nil {
			return nil, err
		}
		realties = append(realties, &realty)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return realties, nil
}

// Validate Realty inpute

// func ValidateRealty(v *validator.Validator, realty *Realty) {
// 	v.Check(realty.Title != "", "title", "Must be provided")
// 	v.Check(len(realty.Title) <= 500, "title", "Must not be more than 500 bytes long")
//
// 	v.Check(movie.Year != 0, "year", "Must be provided")
// 	v.Check(movie.Year >= 1888, "year", "Must be greater than 1888")
// 	v.Check(movie.Year <= int32(time.Now().Year()), "year", "Must not be in the future")
//
// 	v.Check(movie.Runtime != 0, "runtime", "Must be provided")
// 	v.Check(movie.Runtime > 0, "runtime", "Must be a positive integer")
//
// 	v.Check(movie.Genres != nil, "genres", "Must be provided")
// 	v.Check(len(movie.Genres) >= 1, "genres", "Must contain at least 1")
// 	v.Check(len(movie.Genres) <= 5, "genres", "Must not contain more than 5")
// 	v.Check(validator.Unique(movie.Genres), "genres", "Must not contain duplicate values")
// }
