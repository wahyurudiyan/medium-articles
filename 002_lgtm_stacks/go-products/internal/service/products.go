package service

import (
	"context"
	"errors"

	pb "go-products/api/product"
	"go-products/ent"
	"go-products/internal/biz/product"
	"go-products/internal/conf"
)

type ProductsService struct {
	conf      *conf.Data
	productUc product.IProductUsecase

	pb.UnimplementedProductsServer
}

func NewProductsService(c *conf.Data, productUc product.IProductUsecase) *ProductsService {
	return &ProductsService{
		conf:      c,
		productUc: productUc,
	}
}

func (s *ProductsService) CreateProducts(ctx context.Context, req *pb.CreateProductsRequest) (*pb.CreateProductsReply, error) {
	if len(req.Products) == 0 || req.Products[0].Name == "" {
		return nil, errors.New("product(s) is empty")
	}

	var products ent.Products
	for _, p := range req.Products {
		products = append(products, &ent.Product{
			Name:     p.Name,
			Price:    p.Price,
			Quantity: int64(p.Quantity),
		})
	}

	if err := s.productUc.AddProducts(ctx, products); err != nil {
		return nil, err
	}

	return &pb.CreateProductsReply{}, nil
}
func (s *ProductsService) UpdateProducts(ctx context.Context, req *pb.UpdateProductsRequest) (*pb.UpdateProductsReply, error) {
	return &pb.UpdateProductsReply{}, nil
}
func (s *ProductsService) DeleteProducts(ctx context.Context, req *pb.DeleteProductsRequest) (*pb.DeleteProductsReply, error) {
	return &pb.DeleteProductsReply{}, nil
}
func (s *ProductsService) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsReply, error) {
	return &pb.GetProductsReply{}, nil
}
func (s *ProductsService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsReply, error) {
	return &pb.ListProductsReply{}, nil
}
