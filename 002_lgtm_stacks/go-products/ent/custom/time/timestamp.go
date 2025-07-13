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

func (m *MetaTimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Default(time.Now().UTC()).Comment(
			"CreatedAt is a metadata that inform when the enties are inserted",
		),
		field.Time("updated_at").Default(time.Now().UTC()).Comment(
			"UpdatedAt is a metadata that inform when the enties are updated",
		),
	}
}
