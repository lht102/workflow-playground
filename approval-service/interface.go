package approvalservice

import (
	"context"

	"github.com/lht102/workflow-playground/approval-service/ent"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, request *CreatePaymentRequest) (*ent.Payment, error)
	CreatePaymentReview(ctx context.Context, request *CreatePaymentReviewRequest) error
	GetPayment(ctx context.Context, request *GetPaymentRequest) (*ent.Payment, error)
	ListPayments(ctx context.Context, request *ListPaymentsRequest) (*ListPaymentsResponse, error)
}
