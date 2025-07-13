package mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type MetaTimeMixin struct {
	mixin.Schema
}

func (MetaTimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now).Comment(
			"CreatedAt is a timestamp that inform when the enties are inserted",
		),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).Comment(
			"UpdatedAt is a timestamp that inform when the enties are updated",
		),
	}
}
