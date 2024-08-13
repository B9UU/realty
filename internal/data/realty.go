package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/b9uu/realty/internal/validator"
	"github.com/lib/pq"
)

var (
	ErrDuplicateId = errors.New("duplicate id")
	ErrNotFound    = errors.New("record not found")
)

type RealtyModel struct {
	DB *sql.DB
}

type RealtyInterface interface {
	Insert(realty *RealtyInput) error
	GetAll(city string) ([]*RealtyResponse, error)
	AutoComplete(q string) ([]string, error)
}

// inserts into db
func (m RealtyModel) Insert(realty *RealtyInput) error {
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

// gets all realties from db
func (m RealtyModel) GetAll(city string) ([]*RealtyResponse, error) {
	query := "SELECT * FROM realty WHERE LOWER(city_name) = LOWER($1) OR $1 = ''"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// args := []interface{}{city}
	rows, err := m.DB.QueryContext(ctx, query, city)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	realties := []*RealtyResponse{}
	var hasRows bool
	for rows.Next() {
		hasRows = true
		var realty RealtyResponse
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
	if !hasRows {
		return nil, ErrNotFound
	}
	return realties, nil
}

func (m RealtyModel) AutoComplete(q string) ([]string, error) {
	query := `SELECT DISTINCT city_name FROM "realty" WHERE city_name ILIKE '%' || $1 || '%'`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, q)
	if err != nil {
		return nil, err
	}
	var hasRows bool
	results := []string{}
	for rows.Next() {
		hasRows = true
		var city string
		err := rows.Scan(&city)
		if err != nil {
			return nil, err
		}
		results = append(results, city)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !hasRows {
		return nil, ErrNotFound
	}
	return results, nil
}

// Validate AutoComplete city input
func ValidateQuery(v *validator.Validator, q string) {
	v.Check(q != "", "q", "must be provided")
	v.Check(len(q) <= 100, "q", "must not be more than 100 bytes long")
	v.Check(len(q) >= 3, "q", "must not be less than 3 characters long")
}
func ValidateCity(v *validator.Validator, city string) {
	v.Check(len(city) <= 100, "city", "must not be more than 100 bytes long")
	v.Check(len(city) >= 3, "city", "must not be less than 3 characters long")
}
