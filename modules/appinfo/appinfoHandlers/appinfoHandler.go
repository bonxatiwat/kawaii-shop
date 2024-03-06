package appinfoHandlers

import (
	"github.com/bonxatiwat/kawaii-shop-tutortial/config"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/appinfo"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/appinfo/appinfoUsecases"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/entities"
	"github.com/bonxatiwat/kawaii-shop-tutortial/pkg/kawaiiauth"
	"github.com/gofiber/fiber/v2"
)

type appinfoHandlerErrCode string

const (
	generateApiKeyErr appinfoHandlerErrCode = "appinfo-001"
	findCategoryErr   appinfoHandlerErrCode = "appinfo-002"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg             config.IConfig
	appinfoUsecases appinfoUsecases.IAppinfUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecase appinfoUsecases.IAppinfUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:             cfg,
		appinfoUsecases: appinfoUsecase,
	}
}
func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := kawaiiauth.NewKawaiiAuth(
		kawaiiauth.ApiKey,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	catogory, err := h.appinfoUsecases.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, catogory).Res()
}
