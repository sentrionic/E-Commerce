package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common/order"
)

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

func (Order) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.Enum("status").GoType(order.Created),
		field.Uint("price").Default(0),
		field.UUID("user_id", uuid.UUID{}),
		field.Uint("version").Default(0),
	}
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
	return nil
}
