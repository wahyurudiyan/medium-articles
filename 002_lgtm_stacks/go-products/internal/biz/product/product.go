package product

import (
	"context"
	"go-products/internal/conf"
	productRepo "go-products/internal/data/product"
)

type usecase struct {
	conf        *conf.Bootstrap
	productRepo *productRepo.Queries
}

type IProductUsecase interface {
	AddProducts(ctx context.Context, items []productRepo.CreateProductsParams) error
}

func NewProductUsecase(c *conf.Bootstrap, productRepo *productRepo.Queries) IProductUsecase {
	return &usecase{
		conf:        c,
		productRepo: productRepo,
	}
}

func (uc *usecase) AddProducts(ctx context.Context, items []productRepo.CreateProductsParams) error {
	if _, err := uc.productRepo.CreateProducts(ctx, items); err != nil {
		return err
	}

	return nil
}
