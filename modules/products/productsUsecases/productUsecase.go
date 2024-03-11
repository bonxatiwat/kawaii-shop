package productsUsecases

import "github.com/bonxatiwat/kawaii-shop-tutortial/modules/products/productsRepositories"

type IProductsUsecase interface {
}

type productsUsecase struct {
	productsRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productsRepositories productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productsRepository: productsRepositories,
	}
}
