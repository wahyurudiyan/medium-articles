package data

import (
	"context"
	"go-products/ent"
)

type IProductRepo interface {
	CreateProducts(ctx context.Context, items ent.Products) error
}

func NewProductRepo(data *Data) IProductRepo {
	return data
}

func (d *Data) CreateProducts(ctx context.Context, items ent.Products) error {
	var products []*ent.ProductCreate
	for _, item := range items {
		create := d.dbCli.Product.Create()
		products = append(products,
			create.
				SetName(item.Name).
				SetQuantity(item.Quantity).
				SetPrice(item.Price),
		)
	}

	if _, err := d.dbCli.Product.CreateBulk(products...).Save(ctx); err != nil {
		return err
	}

	return nil
}
