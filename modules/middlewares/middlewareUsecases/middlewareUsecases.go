package middlewareUsecases

import "github.com/bonxatiwat/kawaii-shop-tutortial/modules/middlewares/middlewareRepositories"

type IMiddelwaresUsecase interface {
}

type middlewaresUsecase struct {
	middlewareReposiroty middlewareRepositories.IMiddelwaresRespository
}

func MiddlewaresUsecase(middlewareReposiroty middlewareRepositories.IMiddelwaresRespository) IMiddelwaresUsecase {
	return &middlewaresUsecase{
		middlewareReposiroty: middlewareReposiroty,
	}
}
