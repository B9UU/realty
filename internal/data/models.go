package data

import "database/sql"

type Models struct {
	Realty RealtyModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Realty: RealtyModel{DB: db},
	}
}
