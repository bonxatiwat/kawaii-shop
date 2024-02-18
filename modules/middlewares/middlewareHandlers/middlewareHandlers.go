package middlewareHandlers

import (
	"github.com/bonxatiwat/kawaii-shop-tutortial/config"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/entities"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/middlewares/middlewareUsecases"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middleware-001"
)

type IMiddelwaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
}

type middlewaresHandler struct {
	cfg                config.IConfig
	middlewareUsecases middlewareUsecases.IMiddelwaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewareUsecases middlewareUsecases.IMiddelwaresUsecase) IMiddelwaresHandler {
	return &middlewaresHandler{
		cfg:                cfg,
		middlewareUsecases: middlewareUsecases,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"routter not found",
		).Res()
	}
}
