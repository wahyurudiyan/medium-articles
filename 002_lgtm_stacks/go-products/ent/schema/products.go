package schema

import (
	"go-products/ent/custom/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Products holds the schema definition for the Products entity.
type Product struct {
	ent.Schema
}

func (Product) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.MetaTimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}

// Fields of the Products.
func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StructTag(`json:"id,omitempty"`).Comment(
			"id is an incremental value for database primary key",
		),
		field.UUID("sku", uuid.UUID{}).
			Unique().Optional().Default(uuid.New).
			StructTag(`json:"sku,omitempty"`).Comment(
			"SKU is Stock Keeping Unit that become a unique id for products",
		),
		field.String("name").
			MaxLen(120).StructTag(`json:"name"`).Comment(
			"Name field contain name of products",
		),
		field.Int64("quantity").
			StructTag(`json:"qty"`).Comment(
			"Quantity is amount of products that available in inventory",
		),
		field.Int64("price").StructTag(`json:"price"`).Comment(
			"Price is the cost of the product in the smallest currency unit",
		),
	}
}

func (Product) Indexes() []ent.Index {
	return []ent.Index{}
}

// Edges of the Products.
func (Product) Edges() []ent.Edge {
	return nil
}
