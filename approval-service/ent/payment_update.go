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

// PaymentUpdate is the builder for updating Payment entities.
type PaymentUpdate struct {
	config
	hooks    []Hook
	mutation *PaymentMutation
}

// Where appends a list predicates to the PaymentUpdate builder.
func (pu *PaymentUpdate) Where(ps ...predicate.Payment) *PaymentUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetStatus sets the "status" field.
func (pu *PaymentUpdate) SetStatus(pa payment.Status) *PaymentUpdate {
	pu.mutation.SetStatus(pa)
	return pu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (pu *PaymentUpdate) SetNillableStatus(pa *payment.Status) *PaymentUpdate {
	if pa != nil {
		pu.SetStatus(*pa)
	}
	return pu
}

// SetRemark sets the "remark" field.
func (pu *PaymentUpdate) SetRemark(s string) *PaymentUpdate {
	pu.mutation.SetRemark(s)
	return pu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (pu *PaymentUpdate) SetNillableRemark(s *string) *PaymentUpdate {
	if s != nil {
		pu.SetRemark(*s)
	}
	return pu
}

// ClearRemark clears the value of the "remark" field.
func (pu *PaymentUpdate) ClearRemark() *PaymentUpdate {
	pu.mutation.ClearRemark()
	return pu
}

// SetUpdateTime sets the "update_time" field.
func (pu *PaymentUpdate) SetUpdateTime(t time.Time) *PaymentUpdate {
	pu.mutation.SetUpdateTime(t)
	return pu
}

// AddReviewIDs adds the "reviews" edge to the Review entity by IDs.
func (pu *PaymentUpdate) AddReviewIDs(ids ...int64) *PaymentUpdate {
	pu.mutation.AddReviewIDs(ids...)
	return pu
}

// AddReviews adds the "reviews" edges to the Review entity.
func (pu *PaymentUpdate) AddReviews(r ...*Review) *PaymentUpdate {
	ids := make([]int64, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return pu.AddReviewIDs(ids...)
}

// Mutation returns the PaymentMutation object of the builder.
func (pu *PaymentUpdate) Mutation() *PaymentMutation {
	return pu.mutation
}

// ClearReviews clears all "reviews" edges to the Review entity.
func (pu *PaymentUpdate) ClearReviews() *PaymentUpdate {
	pu.mutation.ClearReviews()
	return pu
}

// RemoveReviewIDs removes the "reviews" edge to Review entities by IDs.
func (pu *PaymentUpdate) RemoveReviewIDs(ids ...int64) *PaymentUpdate {
	pu.mutation.RemoveReviewIDs(ids...)
	return pu
}

// RemoveReviews removes "reviews" edges to Review entities.
func (pu *PaymentUpdate) RemoveReviews(r ...*Review) *PaymentUpdate {
	ids := make([]int64, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return pu.RemoveReviewIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *PaymentUpdate) Save(ctx context.Context) (int, error) {
	pu.defaults()
	return withHooks[int, PaymentMutation](ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *PaymentUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *PaymentUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *PaymentUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pu *PaymentUpdate) defaults() {
	if _, ok := pu.mutation.UpdateTime(); !ok {
		v := payment.UpdateDefaultUpdateTime()
		pu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pu *PaymentUpdate) check() error {
	if v, ok := pu.mutation.Status(); ok {
		if err := payment.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Payment.status": %w`, err)}
		}
	}
	if v, ok := pu.mutation.Remark(); ok {
		if err := payment.RemarkValidator(v); err != nil {
			return &ValidationError{Name: "remark", err: fmt.Errorf(`ent: validator failed for field "Payment.remark": %w`, err)}
		}
	}
	return nil
}

func (pu *PaymentUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pu.check(); err != nil {
		return n, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   payment.Table,
			Columns: payment.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: payment.FieldID,
			},
		},
	}
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.Status(); ok {
		_spec.SetField(payment.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := pu.mutation.Remark(); ok {
		_spec.SetField(payment.FieldRemark, field.TypeString, value)
	}
	if pu.mutation.RemarkCleared() {
		_spec.ClearField(payment.FieldRemark, field.TypeString)
	}
	if value, ok := pu.mutation.UpdateTime(); ok {
		_spec.SetField(payment.FieldUpdateTime, field.TypeTime, value)
	}
	if pu.mutation.ReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   payment.ReviewsTable,
			Columns: []string{payment.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: review.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedReviewsIDs(); len(nodes) > 0 && !pu.mutation.ReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   payment.ReviewsTable,
			Columns: []string{payment.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: review.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.ReviewsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   payment.ReviewsTable,
			Columns: []string{payment.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: review.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{payment.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// PaymentUpdateOne is the builder for updating a single Payment entity.
type PaymentUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *PaymentMutation
}

// SetStatus sets the "status" field.
func (puo *PaymentUpdateOne) SetStatus(pa payment.Status) *PaymentUpdateOne {
	puo.mutation.SetStatus(pa)
	return puo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (puo *PaymentUpdateOne) SetNillableStatus(pa *payment.Status) *PaymentUpdateOne {
	if pa != nil {
		puo.SetStatus(*pa)
	}
	return puo
}

// SetRemark sets the "remark" field.
func (puo *PaymentUpdateOne) SetRemark(s string) *PaymentUpdateOne {
	puo.mutation.SetRemark(s)
	return puo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (puo *PaymentUpdateOne) SetNillableRemark(s *string) *PaymentUpdateOne {
	if s != nil {
		puo.SetRemark(*s)
	}
	return puo
}

// ClearRemark clears the value of the "remark" field.
func (puo *PaymentUpdateOne) ClearRemark() *PaymentUpdateOne {
	puo.mutation.ClearRemark()
	return puo
}

// SetUpdateTime sets the "update_time" field.
func (puo *PaymentUpdateOne) SetUpdateTime(t time.Time) *PaymentUpdateOne {
	puo.mutation.SetUpdateTime(t)
	return puo
}

// AddReviewIDs adds the "reviews" edge to the Review entity by IDs.
func (puo *PaymentUpdateOne) AddReviewIDs(ids ...int64) *PaymentUpdateOne {
	puo.mutation.AddReviewIDs(ids...)
	return puo
}

// AddReviews adds the "reviews" edges to the Review entity.
func (puo *PaymentUpdateOne) AddReviews(r ...*Review) *PaymentUpdateOne {
	ids := make([]int64, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return puo.AddReviewIDs(ids...)
}

// Mutation returns the PaymentMutation object of the builder.
func (puo *PaymentUpdateOne) Mutation() *PaymentMutation {
	return puo.mutation
}

// ClearReviews clears all "reviews" edges to the Review entity.
func (puo *PaymentUpdateOne) ClearReviews() *PaymentUpdateOne {
	puo.mutation.ClearReviews()
	return puo
}

// RemoveReviewIDs removes the "reviews" edge to Review entities by IDs.
func (puo *PaymentUpdateOne) RemoveReviewIDs(ids ...int64) *PaymentUpdateOne {
	puo.mutation.RemoveReviewIDs(ids...)
	return puo
}

// RemoveReviews removes "reviews" edges to Review entities.
func (puo *PaymentUpdateOne) RemoveReviews(r ...*Review) *PaymentUpdateOne {
	ids := make([]int64, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return puo.RemoveReviewIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *PaymentUpdateOne) Select(field string, fields ...string) *PaymentUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Payment entity.
func (puo *PaymentUpdateOne) Save(ctx context.Context) (*Payment, error) {
	puo.defaults()
	return withHooks[*Payment, PaymentMutation](ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *PaymentUpdateOne) SaveX(ctx context.Context) *Payment {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *PaymentUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *PaymentUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (puo *PaymentUpdateOne) defaults() {
	if _, ok := puo.mutation.UpdateTime(); !ok {
		v := payment.UpdateDefaultUpdateTime()
		puo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (puo *PaymentUpdateOne) check() error {
	if v, ok := puo.mutation.Status(); ok {
		if err := payment.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Payment.status": %w`, err)}
		}
	}
	if v, ok := puo.mutation.Remark(); ok {
		if err := payment.RemarkValidator(v); err != nil {
			return &ValidationError{Name: "remark", err: fmt.Errorf(`ent: validator failed for field "Payment.remark": %w`, err)}
		}
	}
	return nil
}

func (puo *PaymentUpdateOne) sqlSave(ctx context.Context) (_node *Payment, err error) {
	if err := puo.check(); err != nil {
		return _node, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   payment.Table,
			Columns: payment.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: payment.FieldID,
			},
		},
	}
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Payment.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, payment.FieldID)
		for _, f := range fields {
			if !payment.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != payment.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := puo.mutation.Status(); ok {
		_spec.SetField(payment.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := puo.mutation.Remark(); ok {
		_spec.SetField(payment.FieldRemark, field.TypeString, value)
	}
	if puo.mutation.RemarkCleared() {
		_spec.ClearField(payment.FieldRemark, field.TypeString)
	}
	if value, ok := puo.mutation.UpdateTime(); ok {
		_spec.SetField(payment.FieldUpdateTime, field.TypeTime, value)
	}
	if puo.mutation.ReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   payment.ReviewsTable,
			Columns: []string{payment.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: review.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedReviewsIDs(); len(nodes) > 0 && !puo.mutation.ReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   payment.ReviewsTable,
			Columns: []string{payment.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: review.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.ReviewsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   payment.ReviewsTable,
			Columns: []string{payment.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: review.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Payment{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{payment.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}