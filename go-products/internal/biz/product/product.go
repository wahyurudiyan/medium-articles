package product

import (
	"context"
	"encoding/json"
	"go-products/internal/conf"
	productRepo "go-products/internal/data/product"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/trace"
)

type usecase struct {
	conf        *conf.Bootstrap
	tracer      trace.Tracer
	productRepo *productRepo.Queries
}

type IProductUsecase interface {
	CreateProducts(ctx context.Context, items []productRepo.CreateProductsParams) error
	FetchProducts(ctx context.Context, page int32, limit int32) ([]productRepo.Product, error)
	FetchProductByID(ctx context.Context, id int64) (*productRepo.Product, error)
	FetchProductBySKU(ctx context.Context, sku string) (*productRepo.Product, error)
	UpdateProduct(ctx context.Context, product *productRepo.Product) (*productRepo.Product, error)
	DeleteProduct(ctx context.Context, sku string) error
}

func NewProductUsecase(c *conf.Bootstrap, tracer trace.Tracer, productRepo *productRepo.Queries) IProductUsecase {
	return &usecase{
		conf:        c,
		tracer:      tracer,
		productRepo: productRepo,
	}
}

func (uc *usecase) CreateProducts(ctx context.Context, products []productRepo.CreateProductsParams) error {
	ctx, span := uc.tracer.Start(ctx, "biz.IProductUsecase.CreateProducts")
	defer span.End()

	log.Context(ctx).Debugw("msg", "Creating products with batch")

	batchSize := 5
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}
		batchItems := products[i:end]

		data, _ := json.Marshal(batchItems)
		log.Context(ctx).Debugw("msg", "Creating products", map[string]interface{}{
			"batch.id": i, "products": string(data),
		})
		if _, err := uc.productRepo.CreateProducts(ctx, batchItems); err != nil {
			return err
		}
	}

	return nil
}

func (uc *usecase) FetchProducts(ctx context.Context, page int32, limit int32) ([]productRepo.Product, error) {
	ctx, span := uc.tracer.Start(ctx, "biz.IProductUsecase.FetchProducts")
	defer span.End()

	log.Context(ctx).Debugw("msg", "Fetching list of products", map[string]interface{}{
		"page": page, "page_size": limit,
	})

	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	products, err := uc.productRepo.ReadListProducts(ctx, productRepo.ReadListProductsParams{Offset: int32(offset), Limit: limit})
	if err != nil {
		log.Context(ctx).Errorw("msg", "Error occur when updating product by ID!", map[string]interface{}{
			"page": page, "page_size": limit,
		})
		return nil, err
	}

	log.Context(ctx).Debugw("msg", "Products successfully loaded", map[string]interface{}{
		"page": page, "page_size": limit,
	})

	return products, nil
}

func (uc *usecase) FetchProductByID(ctx context.Context, id int64) (*productRepo.Product, error) {
	ctx, span := uc.tracer.Start(ctx, "biz.IProductUsecase.FetchProductByID")
	defer span.End()

	log.Context(ctx).Debugw("msg", "Fetching product (by ID)", map[string]interface{}{
		"product.id": id,
	})

	product, err := uc.productRepo.ReadProductByID(ctx, id)
	if err != nil {
		log.Context(ctx).Errorw("msg", "Error occur when reading product by ID!", map[string]interface{}{
			"product.sku": product.Sku, "product.name": product.Name,
		})
		return nil, err
	}

	log.Context(ctx).Debugw("msg", "Product found!", map[string]interface{}{
		"product.id": id,
	})

	return &product, nil
}

func (uc *usecase) FetchProductBySKU(ctx context.Context, sku string) (*productRepo.Product, error) {
	ctx, span := uc.tracer.Start(ctx, "biz.IProductUsecase.FetchProductBySKU")
	defer span.End()

	log.Context(ctx).Debugw("msg", "Fetching product (by SKU)", map[string]interface{}{
		"product.sku": sku,
	})

	product, err := uc.productRepo.ReadProductBySKU(ctx, sku)
	if err != nil {
		log.Context(ctx).Errorw("msg", "Error occur when reading product by SKU!", map[string]interface{}{
			"product.sku": product.Sku, "product.name": product.Name,
		})
		return nil, err
	}

	log.Context(ctx).Debugw("msg", "Product found!", map[string]interface{}{
		"product.sku": sku,
	})

	return &product, nil
}

func (uc *usecase) UpdateProduct(ctx context.Context, product *productRepo.Product) (*productRepo.Product, error) {
	ctx, span := uc.tracer.Start(ctx, "biz.IProductUsecase.UpdateProduct")
	defer span.End()

	log.Context(ctx).Debugw("msg", "Updating product", map[string]interface{}{
		"product.sku": product.Sku, "product.name": product.Name,
	})

	updatedProduct, err := uc.productRepo.UpdateProduct(ctx, productRepo.UpdateProductParams{
		Sku:      product.Sku,
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	})
	if err != nil {
		log.Context(ctx).Errorw("msg", "Error occur when updating product!", map[string]interface{}{
			"product.sku": product.Sku, "product.name": product.Name,
		})
		return nil, err
	}

	log.Context(ctx).Debugw("msg", "Product updated successfully!", map[string]interface{}{
		"product.sku": product.Sku, "product.name": product.Name,
	})

	return &updatedProduct, nil
}

func (uc *usecase) DeleteProduct(ctx context.Context, sku string) error {
	ctx, span := uc.tracer.Start(ctx, "biz.IProductUsecase.DeleteProduct")
	defer span.End()

	log.Context(ctx).Debugw("msg", "Process deleting product", map[string]interface{}{
		"sku": sku,
	})

	_, err := uc.productRepo.SoftDeleteProduct(ctx, sku)
	if err != nil {
		log.Context(ctx).Errorw("msg", "Error occur when deleting product!", map[string]interface{}{
			"error": err, "product.sku": sku,
		})
		return err
	}

	log.Context(ctx).Debugw("msg", "Product deleted successfully!", map[string]interface{}{
		"sku": sku,
	})

	return nil
}
