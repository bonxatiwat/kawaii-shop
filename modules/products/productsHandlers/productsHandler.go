package productsHandlers

import (
	"strings"

	"github.com/bonxatiwat/kawaii-shop-tutortial/config"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/entities"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/files/filesUsecases"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)

type productsHandlersErrCode string

const (
	findOneProductErr productsHandlersErrCode = "products-001"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
}

type productsHandler struct {
	cfg             config.IConfig
	productsUsecase productsUsecases.IProductsUsecase
	filesUsecases   filesUsecases.IFilesUsecase
}

func ProductsHandler(cfg config.IConfig, productsUsecase productsUsecases.IProductsUsecase, filesUsecases filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHandler{
		cfg:             cfg,
		productsUsecase: productsUsecase,
		filesUsecases:   filesUsecases,
	}
}

func (h *productsHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), "")

	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}
