package schema

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	gen "github.com/sentrionic/ecommerce/auth/ent"
	"github.com/sentrionic/ecommerce/auth/ent/hook"
	"regexp"
)

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("username").
			Unique().
			MinLen(3).
			MaxLen(30).
			NotEmpty(),
		field.String("email").
			Match(emailRegexp).
			NotEmpty().
			Unique(),
		field.Text("password").
			MinLen(6).
			NotEmpty().
			Sensitive(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(func(next ent.Mutator) ent.Mutator {
			return hook.UserFunc(func(ctx context.Context, mutation *gen.UserMutation) (gen.Value, error) {
				if password, ok := mutation.Password(); ok {
					hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
					if err != nil {
						return nil, err
					}

					mutation.SetPassword(hash)
				}

				return next.Mutate(ctx, mutation)
			})
		}, ent.OpCreate|ent.OpUpdateOne),
	}
}
