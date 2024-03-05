package middlewareRepositories

import (
	"fmt"

	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/middlewares"
	"github.com/jmoiron/sqlx"
)

type IMiddelwaresRespository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresRespository struct {
	db *sqlx.DB
}

func MiddlewaresRespository(db *sqlx.DB) IMiddelwaresRespository {
	return &middlewaresRespository{
		db: db,
	}
}

func (r *middlewaresRespository) FindAccessToken(userId, accessToken string) bool {
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "oauth"
	WHERE "user_id" = $1
	AND "access_token" = $2;`

	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}
	return true
}

func (r *middlewaresRespository) FindRole() ([]*middlewares.Role, error) {
	query := `
	SELECT
		"id",
		"title"
	FROM "roles"
	ORDER BY "id" DESC;`

	roles := make([]*middlewares.Role, 0)
	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empty")
	}
	return roles, nil
}
