// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/ent/predicate"
)

// PaymentDelete is the builder for deleting a Payment entity.
type PaymentDelete struct {
	config
	hooks    []Hook
	mutation *PaymentMutation
}

// Where appends a list predicates to the PaymentDelete builder.
func (pd *PaymentDelete) Where(ps ...predicate.Payment) *PaymentDelete {
	pd.mutation.Where(ps...)
	return pd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pd *PaymentDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, PaymentMutation](ctx, pd.sqlExec, pd.mutation, pd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pd *PaymentDelete) ExecX(ctx context.Context) int {
	n, err := pd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pd *PaymentDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: payment.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: payment.FieldID,
			},
		},
	}
	if ps := pd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pd.mutation.done = true
	return affected, err
}

// PaymentDeleteOne is the builder for deleting a single Payment entity.
type PaymentDeleteOne struct {
	pd *PaymentDelete
}

// Exec executes the deletion query.
func (pdo *PaymentDeleteOne) Exec(ctx context.Context) error {
	n, err := pdo.pd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{payment.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pdo *PaymentDeleteOne) ExecX(ctx context.Context) {
	pdo.pd.ExecX(ctx)
}
