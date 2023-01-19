package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Payment holds the schema definition for the Payment entity.
type Payment struct {
	ent.Schema
}

// Fields of the Payment.
func (Payment) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.UUID("request_id", uuid.UUID{}).
			Immutable().
			Default(uuid.New),
		field.Enum("status").
			Values("PENDING", "APPROVED", "REJECTED").
			Default("PENDING"),
		field.String("remark").
			NotEmpty().
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
	}
}

// Edges of the Payment.
func (Payment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("reviews", Review.Type),
	}
}

func (Payment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("request_id").
			Unique(),
		index.Fields("create_time", "id").
			Annotations(entsql.DescColumns("create_time", "id")),
	}
}
