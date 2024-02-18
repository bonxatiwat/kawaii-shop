package usersHandles

import (
	"github.com/bonxatiwat/kawaii-shop-tutortial/config"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/users/usersUsecases"
)

type IUsersHandle interface {
}

type usersHandle struct {
	cfg          *config.IConfig
	usersUsecase usersUsecases.IUsersUsecase
}

func UsersHandle(cfg *config.IConfig, usersUsecase usersUsecases.IUsersUsecase) IUsersHandle {
	return &usersHandle{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}
