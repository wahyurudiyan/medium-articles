package biz

import (
	"go-products/internal/biz/product"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(product.NewProductUsecase)
