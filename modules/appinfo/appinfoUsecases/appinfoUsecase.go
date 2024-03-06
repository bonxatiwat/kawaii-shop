package appinfoUsecases

import "github.com/bonxatiwat/kawaii-shop-tutortial/modules/appinfo/appinfoRepositories"

type IAppinfUsecase interface {
}

type appinfoUsecase struct {
	appinfoRepository appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecases(appinfoRepository appinfoRepositories.IAppinfoRepository) IAppinfUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}
