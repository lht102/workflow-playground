package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Review holds the schema definition for the Review entity.
type Review struct {
	ent.Schema
}

// Fields of the Review.
func (Review) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Enum("event").
			Values("APPROVE", "REJECT"),
		field.String("reviewer_id").
			Immutable().
			NotEmpty(),
		field.String("comment").
			Optional().
			Nillable(),
		field.Time("create_time").
			Immutable().
			Default(time.Now).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime(6)",
			}),
		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime(6)",
			}),
		field.Int64("payment_id"),
	}
}

// Edges of the Review.
func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("payment", Payment.Type).
			Ref("reviews").
			Required().
			Field("payment_id").
			Unique(),
	}
}

func (Review) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("payment_id", "reviewer_id").
			Unique(),
	}
}
