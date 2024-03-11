package productsUsecases

import (
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/products"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
}

type productsUsecase struct {
	productsRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productsRepositories productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productsRepository: productsRepositories,
	}
}

func (u *productsUsecase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.productsRepository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}
