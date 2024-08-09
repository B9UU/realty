package data

import "database/sql"

// cintains all models the application needs
type Models struct {
	Realty RealtyInterface
}

// initiate new models
func NewModels(db *sql.DB) Models {
	return Models{
		Realty: RealtyModel{DB: db},
	}
}
