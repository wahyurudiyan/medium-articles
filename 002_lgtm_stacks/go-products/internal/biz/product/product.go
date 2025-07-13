package product

import (
	"context"
	"go-products/ent"
	"go-products/internal/conf"
	"go-products/internal/data"
)

type usecase struct {
	conf        *conf.Data
	productRepo data.IProductRepo
}

type IProductUsecase interface {
	AddProducts(ctx context.Context, items ent.Products) error
}

func NewProductUsecase(c *conf.Data, productRepo data.IProductRepo) IProductUsecase {
	return &usecase{
		conf:        c,
		productRepo: productRepo,
	}
}

func (uc *usecase) AddProducts(ctx context.Context, items ent.Products) error {
	if err := uc.productRepo.CreateProducts(ctx, items); err != nil {
		return err
	}
	return nil
}
