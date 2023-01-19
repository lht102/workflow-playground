package approvalservice

import (
	"github.com/google/uuid"
	"github.com/lht102/workflow-playground/approval-service/ent"
	"github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
	"gopkg.in/guregu/null.v4"
)

type CreatePaymentRequest struct {
	RequestID *uuid.UUID
}

type GetPaymentRequest struct {
	PaymentID int64
}

type ListPaymentsRequest struct {
	BeforeCreateTime  null.Time
	LessThanPaymentID null.Int
	PageSize          null.Int

	PaymentStatuses *[]payment.Status
}

type ListPaymentsResponse struct {
	Payments ent.Payments
	HasNext  bool
}

type UpdatePaymentRequest struct {
	PaymentID     int64
	PaymentStatus *payment.Status
	Remark        null.String
}

type CreatePaymentReviewRequest struct {
	PaymentID   int64
	ReviewEvent review.Event
	ReviewerID  string
	Comment     null.String
}
