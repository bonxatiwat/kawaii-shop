package appinfoUsecases

import (
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/appinfo"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/appinfo/appinfoRepositories"
)

type IAppinfUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
}

type appinfoUsecase struct {
	appinfoRepository appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecases(appinfoRepository appinfoRepositories.IAppinfoRepository) IAppinfUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	category, err := u.appinfoRepository.FindCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error {
	if err := u.appinfoRepository.InsertCategory(req); err != nil {
		return err
	}
	return nil
}
