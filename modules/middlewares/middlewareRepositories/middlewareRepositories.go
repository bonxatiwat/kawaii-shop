package middlewareRepositories

import "github.com/jmoiron/sqlx"

type IMiddelwaresRespository interface {
}

type middlewaresRespository struct {
	db *sqlx.DB
}

func MiddlewaresRespository(db *sqlx.DB) IMiddelwaresRespository {
	return &middlewaresRespository{
		db: db,
	}
}
