// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/ent/predicate"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
)

// ReviewUpdate is the builder for updating Review entities.
type ReviewUpdate struct {
	config
	hooks    []Hook
	mutation *ReviewMutation
}

// Where appends a list predicates to the ReviewUpdate builder.
func (ru *ReviewUpdate) Where(ps ...predicate.Review) *ReviewUpdate {
	ru.mutation.Where(ps...)
	return ru
}

// SetEvent sets the "event" field.
func (ru *ReviewUpdate) SetEvent(r review.Event) *ReviewUpdate {
	ru.mutation.SetEvent(r)
	return ru
}

// SetComment sets the "comment" field.
func (ru *ReviewUpdate) SetComment(s string) *ReviewUpdate {
	ru.mutation.SetComment(s)
	return ru
}

// SetNillableComment sets the "comment" field if the given value is not nil.
func (ru *ReviewUpdate) SetNillableComment(s *string) *ReviewUpdate {
	if s != nil {
		ru.SetComment(*s)
	}
	return ru
}

// ClearComment clears the value of the "comment" field.
func (ru *ReviewUpdate) ClearComment() *ReviewUpdate {
	ru.mutation.ClearComment()
	return ru
}

// SetUpdateTime sets the "update_time" field.
func (ru *ReviewUpdate) SetUpdateTime(t time.Time) *ReviewUpdate {
	ru.mutation.SetUpdateTime(t)
	return ru
}

// SetPaymentID sets the "payment_id" field.
func (ru *ReviewUpdate) SetPaymentID(i int64) *ReviewUpdate {
	ru.mutation.SetPaymentID(i)
	return ru
}

// SetPayment sets the "payment" edge to the Payment entity.
func (ru *ReviewUpdate) SetPayment(p *Payment) *ReviewUpdate {
	return ru.SetPaymentID(p.ID)
}

// Mutation returns the ReviewMutation object of the builder.
func (ru *ReviewUpdate) Mutation() *ReviewMutation {
	return ru.mutation
}

// ClearPayment clears the "payment" edge to the Payment entity.
func (ru *ReviewUpdate) ClearPayment() *ReviewUpdate {
	ru.mutation.ClearPayment()
	return ru
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ru *ReviewUpdate) Save(ctx context.Context) (int, error) {
	ru.defaults()
	return withHooks[int, ReviewMutation](ctx, ru.sqlSave, ru.mutation, ru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ru *ReviewUpdate) SaveX(ctx context.Context) int {
	affected, err := ru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ru *ReviewUpdate) Exec(ctx context.Context) error {
	_, err := ru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ru *ReviewUpdate) ExecX(ctx context.Context) {
	if err := ru.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ru *ReviewUpdate) defaults() {
	if _, ok := ru.mutation.UpdateTime(); !ok {
		v := review.UpdateDefaultUpdateTime()
		ru.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ru *ReviewUpdate) check() error {
	if v, ok := ru.mutation.Event(); ok {
		if err := review.EventValidator(v); err != nil {
			return &ValidationError{Name: "event", err: fmt.Errorf(`ent: validator failed for field "Review.event": %w`, err)}
		}
	}
	if _, ok := ru.mutation.PaymentID(); ru.mutation.PaymentCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Review.payment"`)
	}
	return nil
}

func (ru *ReviewUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ru.check(); err != nil {
		return n, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   review.Table,
			Columns: review.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: review.FieldID,
			},
		},
	}
	if ps := ru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ru.mutation.Event(); ok {
		_spec.SetField(review.FieldEvent, field.TypeEnum, value)
	}
	if value, ok := ru.mutation.Comment(); ok {
		_spec.SetField(review.FieldComment, field.TypeString, value)
	}
	if ru.mutation.CommentCleared() {
		_spec.ClearField(review.FieldComment, field.TypeString)
	}
	if value, ok := ru.mutation.UpdateTime(); ok {
		_spec.SetField(review.FieldUpdateTime, field.TypeTime, value)
	}
	if ru.mutation.PaymentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   review.PaymentTable,
			Columns: []string{review.PaymentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: payment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.PaymentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   review.PaymentTable,
			Columns: []string{review.PaymentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: payment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{review.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ru.mutation.done = true
	return n, nil
}

// ReviewUpdateOne is the builder for updating a single Review entity.
type ReviewUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ReviewMutation
}

// SetEvent sets the "event" field.
func (ruo *ReviewUpdateOne) SetEvent(r review.Event) *ReviewUpdateOne {
	ruo.mutation.SetEvent(r)
	return ruo
}

// SetComment sets the "comment" field.
func (ruo *ReviewUpdateOne) SetComment(s string) *ReviewUpdateOne {
	ruo.mutation.SetComment(s)
	return ruo
}

// SetNillableComment sets the "comment" field if the given value is not nil.
func (ruo *ReviewUpdateOne) SetNillableComment(s *string) *ReviewUpdateOne {
	if s != nil {
		ruo.SetComment(*s)
	}
	return ruo
}

// ClearComment clears the value of the "comment" field.
func (ruo *ReviewUpdateOne) ClearComment() *ReviewUpdateOne {
	ruo.mutation.ClearComment()
	return ruo
}

// SetUpdateTime sets the "update_time" field.
func (ruo *ReviewUpdateOne) SetUpdateTime(t time.Time) *ReviewUpdateOne {
	ruo.mutation.SetUpdateTime(t)
	return ruo
}

// SetPaymentID sets the "payment_id" field.
func (ruo *ReviewUpdateOne) SetPaymentID(i int64) *ReviewUpdateOne {
	ruo.mutation.SetPaymentID(i)
	return ruo
}

// SetPayment sets the "payment" edge to the Payment entity.
func (ruo *ReviewUpdateOne) SetPayment(p *Payment) *ReviewUpdateOne {
	return ruo.SetPaymentID(p.ID)
}

// Mutation returns the ReviewMutation object of the builder.
func (ruo *ReviewUpdateOne) Mutation() *ReviewMutation {
	return ruo.mutation
}

// ClearPayment clears the "payment" edge to the Payment entity.
func (ruo *ReviewUpdateOne) ClearPayment() *ReviewUpdateOne {
	ruo.mutation.ClearPayment()
	return ruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ruo *ReviewUpdateOne) Select(field string, fields ...string) *ReviewUpdateOne {
	ruo.fields = append([]string{field}, fields...)
	return ruo
}

// Save executes the query and returns the updated Review entity.
func (ruo *ReviewUpdateOne) Save(ctx context.Context) (*Review, error) {
	ruo.defaults()
	return withHooks[*Review, ReviewMutation](ctx, ruo.sqlSave, ruo.mutation, ruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ruo *ReviewUpdateOne) SaveX(ctx context.Context) *Review {
	node, err := ruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ruo *ReviewUpdateOne) Exec(ctx context.Context) error {
	_, err := ruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruo *ReviewUpdateOne) ExecX(ctx context.Context) {
	if err := ruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ruo *ReviewUpdateOne) defaults() {
	if _, ok := ruo.mutation.UpdateTime(); !ok {
		v := review.UpdateDefaultUpdateTime()
		ruo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ruo *ReviewUpdateOne) check() error {
	if v, ok := ruo.mutation.Event(); ok {
		if err := review.EventValidator(v); err != nil {
			return &ValidationError{Name: "event", err: fmt.Errorf(`ent: validator failed for field "Review.event": %w`, err)}
		}
	}
	if _, ok := ruo.mutation.PaymentID(); ruo.mutation.PaymentCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Review.payment"`)
	}
	return nil
}

func (ruo *ReviewUpdateOne) sqlSave(ctx context.Context) (_node *Review, err error) {
	if err := ruo.check(); err != nil {
		return _node, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   review.Table,
			Columns: review.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: review.FieldID,
			},
		},
	}
	id, ok := ruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Review.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, review.FieldID)
		for _, f := range fields {
			if !review.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != review.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ruo.mutation.Event(); ok {
		_spec.SetField(review.FieldEvent, field.TypeEnum, value)
	}
	if value, ok := ruo.mutation.Comment(); ok {
		_spec.SetField(review.FieldComment, field.TypeString, value)
	}
	if ruo.mutation.CommentCleared() {
		_spec.ClearField(review.FieldComment, field.TypeString)
	}
	if value, ok := ruo.mutation.UpdateTime(); ok {
		_spec.SetField(review.FieldUpdateTime, field.TypeTime, value)
	}
	if ruo.mutation.PaymentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   review.PaymentTable,
			Columns: []string{review.PaymentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: payment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.PaymentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   review.PaymentTable,
			Columns: []string{review.PaymentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: payment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Review{config: ruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{review.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ruo.mutation.done = true
	return _node, nil
}