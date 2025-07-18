package service

import (
	"context"
	"errors"

	pb "go-products/api/product"

	"go-products/internal/biz/product"
	"go-products/internal/conf"
	productRepo "go-products/internal/data/product"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/trace"
)

type ProductsService struct {
	conf      *conf.Bootstrap
	tracer    trace.Tracer
	productUc product.IProductUsecase

	pb.UnimplementedProductsServer
}

func NewProductsService(c *conf.Bootstrap, tracer trace.Tracer, productUc product.IProductUsecase) *ProductsService {
	return &ProductsService{
		conf:      c,
		tracer:    tracer,
		productUc: productUc,
	}
}

func (s *ProductsService) CreateProducts(ctx context.Context, req *pb.CreateProductsRequest) (*pb.CreateProductsReply, error) {
	ctx, span := s.tracer.Start(ctx, "products.ProductService.CreateProducts")
	defer span.End()

	if len(req.Products) == 0 || req.Products[0].Name == "" {
		err := errors.New("product(s) is empty")
		log.Context(ctx).Errorw("Unable to add product(s)!", map[string]interface{}{
			"error": err,
		})
		return nil, err
	}

	var products []productRepo.CreateProductsParams
	for _, p := range req.Products {
		products = append(products, productRepo.CreateProductsParams{
			Sku:      p.SKU,
			Name:     p.Name,
			Price:    p.Price,
			Quantity: int64(p.Quantity),
		})
	}

	if err := s.productUc.AddProducts(ctx, products); err != nil {
		log.Context(ctx).Errorw("Cannot insert product(s)!", map[string]interface{}{
			"error": err,
		})
		return nil, err
	}

	log.Context(ctx).Info("Product saved!")

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
