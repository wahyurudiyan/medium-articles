package interceptor

import (
	"context"
	"fmt"
	"time"

	"go-products/ent"
	internalMixin "go-products/ent/custom/mixin"
	"go-products/ent/hook"
	"go-products/ent/intercept"
	_ "go-products/ent/runtime"

	"entgo.io/ent/dialect/sql"
)

type softDelete struct {
	sd internalMixin.SoftDeleteMixin
}

type ISoftDelete interface {
	Interceptors() []ent.Interceptor
}

func NewSoftDelete(sd internalMixin.SoftDeleteMixin) softDelete {
	return softDelete{sd: sd}
}

// Interceptors of the SoftDeleteMixin.
func (d softDelete) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
			// Skip soft-delete, means include soft-deleted entities.
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}
			d.P(q)
			return nil
		}),
	}
}

// Hooks of the SoftDeleteMixin.
func (d softDelete) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
					// Skip soft-delete, means delete the entity permanently.
					if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
						return next.Mutate(ctx, m)
					}
					mx, ok := m.(interface {
						SetOp(ent.Op)
						Client() *ent.Client
						SetDeleteTime(time.Time)
						WhereP(...func(*sql.Selector))
					})
					if !ok {
						return nil, fmt.Errorf("unexpected mutation type %T", m)
					}
					d.P(mx)
					mx.SetOp(ent.OpUpdate)
					mx.SetDeleteTime(time.Now())
					return mx.Client().Mutate(ctx, m)
				})
			},
			ent.OpDeleteOne|ent.OpDelete,
		),
	}
}

// P adds a storage-level predicate to the queries and mutations.
func (d softDelete) P(w interface{ WhereP(...func(*sql.Selector)) }) {
	w.WhereP(
		sql.FieldIsNull(d.sd.Fields()[0].Descriptor().Name),
	)
}

type softDeleteKey struct{}

// SkipSoftDelete returns a new context that skips the soft-delete interceptor/mutators.
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}
