package service

import (
	"context"
	"errors"

	pb "go-products/api/product"

	"go-products/internal/biz/product"
	"go-products/internal/conf"
	productRepo "go-products/internal/data/product"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
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
		log.Context(ctx).Errorw("msg", "Unable to add product(s)!", "error", err)
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

	if err := s.productUc.CreateProducts(ctx, products); err != nil {
		log.Context(ctx).Errorw("msg", "Cannot insert product(s)!", "error", err)
		return nil, err
	}

	return &pb.CreateProductsReply{
		Meta: &pb.CommonResponse{
			Success: true,
			Message: "created successfully",
		},
	}, nil
}
func (s *ProductsService) UpdateProducts(ctx context.Context, req *pb.UpdateProductsRequest) (*pb.UpdateProductsReply, error) {
	ctx, span := s.tracer.Start(ctx, "products.ProductService.UpdateProducts")
	defer span.End()

	if req.Product == nil || req.Product.SKU == "" {
		err := errors.New("product data or SKU is missing")
		log.Context(ctx).Errorw("msg", "Invalid update request", "error", err)
		return &pb.UpdateProductsReply{
			Meta: &pb.CommonResponse{
				Error: &pb.Error{Code: uuid.NewString(), Reason: err.Error()},
			},
		}, err
	}

	product := &productRepo.Product{
		Sku:      req.Product.SKU,
		Name:     req.Product.Name,
		Price:    req.Product.Price,
		Quantity: int64(req.Product.Quantity),
	}

	_, err := s.productUc.UpdateProduct(ctx, product)
	if err != nil {
		log.Context(ctx).Errorw("msg", "Failed to update product", "error", err)
		return &pb.UpdateProductsReply{
			Meta: &pb.CommonResponse{
				Error: &pb.Error{Code: uuid.NewString(), Reason: err.Error()},
			},
		}, err
	}

	return &pb.UpdateProductsReply{
		Meta: &pb.CommonResponse{
			Success: true,
			Message: "updated successfully",
		},
	}, nil
}

func (s *ProductsService) DeleteProducts(ctx context.Context, req *pb.DeleteProductsRequest) (*pb.DeleteProductsReply, error) {
	ctx, span := s.tracer.Start(ctx, "products.ProductService.DeleteProducts")
	defer span.End()

	if req.Sku == "" {
		err := errors.New("SKU is required for deletion")
		log.Context(ctx).Errorw("msg", "DeleteProducts error", "error", err)
		return &pb.DeleteProductsReply{
			Meta: &pb.CommonResponse{
				Error: &pb.Error{Code: uuid.NewString(), Reason: err.Error()},
			},
		}, err
	}

	if err := s.productUc.DeleteProduct(ctx, req.Sku); err != nil {
		log.Context(ctx).Errorw("msg", "Failed to delete product", "error", err)
		return &pb.DeleteProductsReply{
			Meta: &pb.CommonResponse{
				Error: &pb.Error{Code: uuid.NewString(), Reason: err.Error()},
			},
		}, err
	}

	return &pb.DeleteProductsReply{
		Meta: &pb.CommonResponse{
			Success: true,
			Message: "deleted succesfully",
		},
	}, nil
}

func (s *ProductsService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductReply, error) {
	ctx, span := s.tracer.Start(ctx, "products.ProductService.GetProducts")
	defer span.End()

	if req.Sku == "" && req.Id == 0 {
		err := errors.New("must provide SKU or ID")
		log.Context(ctx).Errorw("msg", "Invalid GetProducts request", "error", err)
		return &pb.GetProductReply{
			Meta: &pb.CommonResponse{
				Error: &pb.Error{Code: uuid.NewString(), Reason: err.Error()},
			},
		}, err
	}

	var (
		product *productRepo.Product
		err     error
	)

	if req.Sku != "" {
		product, err = s.productUc.FetchProductBySKU(ctx, req.Sku)
	} else {
		product, err = s.productUc.FetchProductByID(ctx, req.Id)
	}

	if err != nil {
		log.Context(ctx).Errorw("msg", "Failed to fetch product", "error", err)
		return &pb.GetProductReply{
			Meta: &pb.CommonResponse{
				Error: &pb.Error{Code: uuid.NewString(), Reason: err.Error()},
			},
		}, err
	}

	pbProduct := &pb.Product{
		ID:       product.ID,
		SKU:      product.Sku,
		Name:     product.Name,
		Price:    product.Price,
		Quantity: int32(product.Quantity),
	}

	return &pb.GetProductReply{
		Product: pbProduct,
		Meta: &pb.CommonResponse{
			Success: true,
			Message: "product found",
		},
	}, nil
}

func (s *ProductsService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsReply, error) {
	ctx, span := s.tracer.Start(ctx, "products.ProductService.ListProducts")
	defer span.End()

	products, err := s.productUc.FetchProducts(ctx, req.Page, req.PageSize)
	if err != nil {
		log.Context(ctx).Errorw("msg", "Failed to list products", "error", err)
		return &pb.ListProductsReply{
			Meta: &pb.CommonResponse{
				Error: &pb.Error{Code: uuid.NewString(), Reason: err.Error()},
			},
		}, err
	}

	var result []*pb.Product
	for _, p := range products {
		result = append(result, &pb.Product{
			ID:       p.ID,
			SKU:      p.Sku,
			Name:     p.Name,
			Price:    p.Price,
			Quantity: int32(p.Quantity),
		})
	}

	return &pb.ListProductsReply{
		Products: result,
		Meta: &pb.CommonResponse{
			Success: true,
			Message: "product found",
		},
	}, nil
}
