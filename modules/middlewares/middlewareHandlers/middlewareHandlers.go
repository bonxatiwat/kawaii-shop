package middlewareHandlers

import "github.com/bonxatiwat/kawaii-shop-tutortial/modules/middlewares/middlewareUsecases"

type IMiddelwaresHandler interface {
}

type middlewaresHandler struct {
	middlewareUsecases middlewareUsecases.IMiddelwaresUsecase
}

func MiddlewaresHandler(middlewareUsecases middlewareUsecases.IMiddelwaresUsecase) IMiddelwaresHandler {
	return &middlewaresHandler{
		middlewareUsecases: middlewareUsecases,
	}
}
