package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/ent"
	"github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/ent/predicate"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
)

func createPayment(
	ctx context.Context,
	entClient *ent.PaymentClient,
	request *approvalservice.CreatePaymentRequest,
) (*ent.Payment, error) {
	builder := entClient.
		Create().
		SetStatus(payment.StatusPENDING)

	if request.RequestID != nil {
		builder = builder.SetRequestID(*request.RequestID)
	}

	entityPayment, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	return entityPayment, nil
}

func getPaymentByID(
	ctx context.Context,
	entClient *ent.PaymentClient,
	paymentID int64,
) (*ent.Payment, error) {
	item, err := entClient.
		Query().
		Where(payment.IDEQ(paymentID)).
		WithReviews().
		Only(ctx)

	if err != nil {
		var notFoundErr *ent.NotFoundError
		if errors.As(err, &notFoundErr) {
			return nil, approvalservice.ErrNotFound
		}

		return nil, fmt.Errorf("get payment by id: %w", err)
	}

	return item, nil
}

func getPaymentStatusByID(
	ctx context.Context,
	entClient *ent.PaymentClient,
	paymentID int64,
) (payment.Status, error) {
	paymentStatus, err := entClient.
		Query().
		Where(payment.IDEQ(paymentID)).
		Select(payment.FieldStatus).
		String(ctx)
	if err != nil {
		var notFoundErr *ent.NotFoundError
		if errors.As(err, &notFoundErr) {
			return "", approvalservice.ErrNotFound
		}

		return "", fmt.Errorf("get payment status by id: %w", err)
	}

	return payment.Status(paymentStatus), nil
}

func upsertPaymentReviewExec(
	ctx context.Context,
	entClient *ent.ReviewClient,
	request *approvalservice.CreatePaymentReviewRequest,
) error {
	builder := entClient.
		Create().
		SetPaymentID(request.PaymentID).
		SetEvent(request.ReviewEvent).
		SetReviewerID(request.ReviewerID).
		OnConflict().
		UpdateNewValues()

	if request.Comment.Valid {
		builder = builder.SetComment(request.Comment.String)
	}

	if err := builder.Exec(ctx); err != nil {
		return fmt.Errorf("upsert payment review: %w", err)
	}

	return nil
}

func getNumberOfApprovalsByPaymentID(
	ctx context.Context,
	entClient *ent.ReviewClient,
	paymentID int64,
) (int, error) {
	numOfApprovals, err := entClient.
		Query().
		Where(
			review.And(
				review.PaymentIDEQ(paymentID),
				review.EventEQ(review.EventAPPROVE),
			),
		).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("get number of approvals by payment id: %w", err)
	}

	return numOfApprovals, nil
}

func updatePayment(
	ctx context.Context,
	entClient *ent.PaymentClient,
	request *approvalservice.UpdatePaymentRequest,
) error {
	builder := entClient.
		Update().
		SetNillableStatus(request.PaymentStatus).
		Where(payment.IDEQ(request.PaymentID))

	if request.Remark.Valid {
		builder = builder.SetRemark(request.Remark.String)
	}

	if err := builder.Exec(ctx); err != nil {
		return fmt.Errorf("update payment: %w", err)
	}

	return nil
}

func listPayments(
	ctx context.Context,
	entClient *ent.PaymentClient,
	request *approvalservice.ListPaymentsRequest,
) (*approvalservice.ListPaymentsResponse, error) {
	beforeCreateTime := time.Now()
	if request.BeforeCreateTime.Valid {
		beforeCreateTime = request.BeforeCreateTime.Time
	}

	predicates := []predicate.Payment{
		payment.CreateTimeLT(beforeCreateTime),
	}

	if request.LessThanPaymentID.Valid {
		idPredicate := payment.IDLT(request.LessThanPaymentID.Int64)
		predicates = append(predicates, idPredicate)
	}

	if request.PaymentStatuses != nil {
		predicates = append(predicates, payment.StatusIn(*request.PaymentStatuses...))
	}

	pageSize := defaultPageSize + 1
	if request.PageSize.Valid {
		pageSize = int(request.PageSize.Int64) + 1
	}

	payments, err := entClient.Query().
		Where(
			payment.And(predicates...),
		).
		Order(ent.Desc(payment.FieldCreateTime, payment.FieldID)).
		Limit(pageSize).
		WithReviews().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query payments with filtering: %w", err)
	}

	hasNext := false

	if len(payments) == pageSize {
		payments = payments[:len(payments)-1]
		hasNext = true
	}

	return &approvalservice.ListPaymentsResponse{
		Payments: payments,
		HasNext:  hasNext,
	}, nil
}
