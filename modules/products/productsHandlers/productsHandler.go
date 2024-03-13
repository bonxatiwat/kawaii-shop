package productsHandlers

import (
	"fmt"
	"strings"

	"github.com/bonxatiwat/kawaii-shop-tutortial/config"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/appinfo"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/entities"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/files"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/files/filesUsecases"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/products"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)

type productsHandlersErrCode string

const (
	findOneProductErr productsHandlersErrCode = "products-001"
	findProductsErr   productsHandlersErrCode = "products-002"
	insertProductsErr productsHandlersErrCode = "products-003"
	deleteProductsErr productsHandlersErrCode = "products-004"
	updateProductsErr productsHandlersErrCode = "products-005"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindProducts(c *fiber.Ctx) error
	AddProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
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

func (h *productsHandler) FindProducts(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(findProductsErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}

	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products := h.productsUsecase.FindProducts(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

func (h *productsHandler) AddProduct(c *fiber.Ctx) error {
	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductsErr),
			err.Error(),
		).Res()
	}
	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductsErr),
			"category id is invalid",
		).Res()
	}

	product, err := h.productsUsecase.AddProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProductsErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productsHandler) UpdateProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProductsErr),
			err.Error(),
		).Res()
	}
	req.Id = productId

	product, err := h.productsUsecase.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateProductsErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productsHandler) DeleteProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductsErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range product.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("images/products/%s", p.FileName),
		})
	}
	if err := h.filesUsecases.DeleteFileOnGCP(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductsErr),
			err.Error(),
		).Res()
	}

	if err := h.productsUsecase.DeleteProduct(productId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductsErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
