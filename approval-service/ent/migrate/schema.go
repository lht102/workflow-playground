// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// PaymentsColumns holds the columns for the "payments" table.
	PaymentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "request_id", Type: field.TypeUUID},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"PENDING", "APPROVED", "REJECTED"}, Default: "PENDING"},
		{Name: "remark", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime(6)"}},
		{Name: "update_time", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime(6)"}},
	}
	// PaymentsTable holds the schema information for the "payments" table.
	PaymentsTable = &schema.Table{
		Name:       "payments",
		Columns:    PaymentsColumns,
		PrimaryKey: []*schema.Column{PaymentsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "payment_request_id",
				Unique:  true,
				Columns: []*schema.Column{PaymentsColumns[1]},
			},
			{
				Name:    "payment_create_time_id",
				Unique:  false,
				Columns: []*schema.Column{PaymentsColumns[4], PaymentsColumns[0]},
				Annotation: &entsql.IndexAnnotation{
					DescColumns: map[string]bool{
						PaymentsColumns[4].Name: true,

						PaymentsColumns[0].Name: true,
					},
				},
			},
		},
	}
	// ReviewsColumns holds the columns for the "reviews" table.
	ReviewsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "event", Type: field.TypeEnum, Enums: []string{"APPROVE", "REJECT"}},
		{Name: "reviewer_id", Type: field.TypeString},
		{Name: "comment", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime(6)"}},
		{Name: "update_time", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime(6)"}},
		{Name: "payment_id", Type: field.TypeInt64},
	}
	// ReviewsTable holds the schema information for the "reviews" table.
	ReviewsTable = &schema.Table{
		Name:       "reviews",
		Columns:    ReviewsColumns,
		PrimaryKey: []*schema.Column{ReviewsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "reviews_payments_reviews",
				Columns:    []*schema.Column{ReviewsColumns[6]},
				RefColumns: []*schema.Column{PaymentsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "review_payment_id_reviewer_id",
				Unique:  true,
				Columns: []*schema.Column{ReviewsColumns[6], ReviewsColumns[2]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		PaymentsTable,
		ReviewsTable,
	}
)

func init() {
	ReviewsTable.ForeignKeys[0].RefTable = PaymentsTable
}