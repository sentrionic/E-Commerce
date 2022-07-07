package schema

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	gen "github.com/sentrionic/ecommerce/products/ent"
	"github.com/sentrionic/ecommerce/products/ent/hook"
)

// Product holds the schema definition for the Product entity.
type Product struct {
	ent.Schema
}

func (Product) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Product.
func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("title").
			NotEmpty(),
		field.Int("price").
			NonNegative(),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("order_id", uuid.UUID{}).Optional().Nillable(),
		field.Uint("version").Default(0),
	}
}

// Edges of the Product.
func (Product) Edges() []ent.Edge {
	return nil
}

func (Product) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(func(next ent.Mutator) ent.Mutator {
			return hook.ProductFunc(func(ctx context.Context, mutation *gen.ProductMutation) (gen.Value, error) {
				version, _ := mutation.OldVersion(ctx)
				mutation.SetVersion(version + 1)
				return next.Mutate(ctx, mutation)
			})
		}, ent.OpUpdateOne),
	}
}
