package productsRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/bonxatiwat/kawaii-shop-tutortial/config"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/entities"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/files/filesUsecases"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/products"
	"github.com/jmoiron/sqlx"
)

type IProductsRepository interface {
	FindOneProduct(productId string) (*products.Product, error)
}

type productsRepository struct {
	db            *sqlx.DB
	cfg           config.IConfig
	filesUsecases filesUsecases.IFilesUsecase
}

func ProductsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IProductsRepository {
	return &productsRepository{
		db:            db,
		cfg:           cfg,
		filesUsecases: filesUsecase,
	}
}

func (r *productsRepository) FindOneProduct(productId string) (*products.Product, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
			SELECT
				"p"."id",
				"p"."title",
				"p"."description",
				"p"."price",
				(
					SELECT
						to_jsonb("ct")
					FROM (
						SELECT
							"c"."id",
							"c"."title"
						FROM "categories" "c"
							LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
						WHERE "pc"."product_id" = "p"."id"
					) AS "ct"
				) AS "category",
				"p"."created_at",
				"p"."updated_at",
				(
					SELECT
						COALESCE(array_to_json(array_agg("it")), '[]'::json)					
					FROM (
						SELECT
							"i"."id",
							"i"."filename",
							"i"."url"
						FROM "images" "i"
						WHERE "i"."product_id" = "p"."id"
					) AS "it"
				) AS "images"
			FROM "products" "p"
			WHERE "p"."id" = $1
			LIMIT 1
	) AS "t"`

	productBytes := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&productBytes, query, productId); err != nil {
		return nil, fmt.Errorf("get product failed: %v", err)
	}
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product failed: %v", err)
	}
	return product, nil
}
