package data

import (
	"context"
	"database/sql"
	"time"
)

type Realty struct {
	ID             int    `json:"id"`
	ListingType    string `json:"listing_type"`
	PromoType      string `json:"promo_type,omitempty"`
	URL            string `json:"url"`
	ProjectName    string `json:"project_name"`
	DisplayAddress string `json:"display_address"`
}

type RealtyModel struct {
	DB *sql.DB
}

func (m RealtyModel) Insert(realty *Realty) error {
	query := `
        INSERT INTO realty (listing_type, promo_type, url, project_name, display_address)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	args := []interface{}{
		realty.ListingType, realty.PromoType, realty.URL,
		realty.ProjectName, realty.DisplayAddress}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&realty.ID)
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
		err := rows.Scan(&realty.ID, &realty.ListingType, &realty.PromoType,
			&realty.URL, &realty.ProjectName, &realty.DisplayAddress)
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
