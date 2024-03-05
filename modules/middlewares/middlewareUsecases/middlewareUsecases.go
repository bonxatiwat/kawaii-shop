package middlewareUsecases

import (
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/middlewares"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/middlewares/middlewareRepositories"
)

type IMiddelwaresUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresUsecase struct {
	middlewareReposiroty middlewareRepositories.IMiddelwaresRespository
}

func MiddlewaresUsecase(middlewareReposiroty middlewareRepositories.IMiddelwaresRespository) IMiddelwaresUsecase {
	return &middlewaresUsecase{
		middlewareReposiroty: middlewareReposiroty,
	}
}

func (u *middlewaresUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewareReposiroty.FindAccessToken(userId, accessToken)
}

func (u *middlewaresUsecase) FindRole() ([]*middlewares.Role, error) {
	roles, err := u.middlewareReposiroty.FindRole()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
