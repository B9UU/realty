package data

import (
	"database/sql"
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
	return m.DB.QueryRow(query, args...).Scan(&realty.ID)
}
