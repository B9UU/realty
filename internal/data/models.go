package data

import "database/sql"

type Models struct {
	Realty RealtyInterface
}

func NewModels(db *sql.DB) Models {
	return Models{
		Realty: RealtyModel{DB: db},
	}
}
