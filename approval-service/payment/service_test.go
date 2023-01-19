package payment_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	approvalservice "github.com/lht102/workflow-playground/approval-service"
	entpayment "github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
	"github.com/lht102/workflow-playground/approval-service/payment"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/mocks"
	"gopkg.in/guregu/null.v4"
)

func (s *PaymentTestSuite) TestCreatePayment() {
	workflowClientMock := &mocks.Client{}
	service := payment.NewService(s.entClient, workflowClientMock)

	workflowClientMock.
		On(
			"ExecuteWorkflow",
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("func(internal.Context, *payment.ApprovalWorkflowRequest) error"),
			mock.Anything,
		).
		Return(&mocks.WorkflowRun{}, nil).
		Once()

	requestID := uuid.New()
	entityPayment, err := service.CreatePayment(context.TODO(), &approvalservice.CreatePaymentRequest{
		RequestID: &requestID,
	})
	s.NoError(err)
	s.Equal(entpayment.StatusPENDING, entityPayment.Status)
	s.Equal(requestID, entityPayment.RequestID)

	_, err = getPaymentByRequestID(s.entClient, requestID)
	s.NoError(err)
}

func (s *PaymentTestSuite) TestCreatePaymentReview() {
	workflowClientMock := &mocks.Client{}
	service := payment.NewService(s.entClient, workflowClientMock)

	s.Run("Approve payment", func() {
		paymentID := int64(1)
		s.entClient.Payment.
			Create().
			SetID(paymentID).
			SetStatus(entpayment.StatusPENDING).
			ExecX(context.TODO())

		reviewerID := "bob@abc.com"
		comment := "This is another comment"

		workflowClientMock.
			On(
				"SignalWithStartWorkflow",
				mock.Anything,
				mock.Anything,
				payment.SignalChannelApprovePaymentChannel,
				&payment.ApprovePaymentSignal{
					ReviewerID: reviewerID,
					Comment:    null.StringFrom(comment),
				},
				mock.Anything,
				mock.AnythingOfType("func(internal.Context, *payment.ApprovalWorkflowRequest) error"),
				&payment.ApprovalWorkflowRequest{
					PaymentID: paymentID,
				},
			).
			Return(&mocks.WorkflowRun{}, nil).
			Once()

		err := service.CreatePaymentReview(context.TODO(), &approvalservice.CreatePaymentReviewRequest{
			PaymentID:   paymentID,
			ReviewEvent: review.EventAPPROVE,
			ReviewerID:  reviewerID,
			Comment:     null.StringFrom(comment),
		})
		s.NoError(err)
	})

	s.Run("Reject payment", func() {
		paymentID := int64(2)
		s.entClient.Payment.
			Create().
			SetID(paymentID).
			SetStatus(entpayment.StatusPENDING).
			ExecX(context.TODO())

		reviewerID := "alice@abc.com"
		comment := "This is a comment"

		workflowClientMock.
			On(
				"SignalWithStartWorkflow",
				mock.Anything,
				mock.Anything,
				payment.SignalChannelRejectPaymentChannel,
				&payment.RejectPaymentSignal{
					ReviewerID: reviewerID,
					Comment:    null.StringFrom(comment),
				},
				mock.Anything,
				mock.AnythingOfType("func(internal.Context, *payment.ApprovalWorkflowRequest) error"),
				&payment.ApprovalWorkflowRequest{
					PaymentID: paymentID,
				},
			).
			Return(&mocks.WorkflowRun{}, nil).
			Once()

		err := service.CreatePaymentReview(context.TODO(), &approvalservice.CreatePaymentReviewRequest{
			PaymentID:   paymentID,
			ReviewEvent: review.EventREJECT,
			ReviewerID:  reviewerID,
			Comment:     null.StringFrom(comment),
		})
		s.NoError(err)
	})

	s.Run("No signal for non pending payment", func() {
		paymentID := int64(3)
		s.entClient.Payment.
			Create().
			SetID(paymentID).
			SetStatus(entpayment.StatusREJECTED).
			ExecX(context.TODO())

		reviewerID := "alice@abc.com"
		comment := "This is a comment"

		err := service.CreatePaymentReview(context.TODO(), &approvalservice.CreatePaymentReviewRequest{
			PaymentID:   paymentID,
			ReviewEvent: review.EventREJECT,
			ReviewerID:  reviewerID,
			Comment:     null.StringFrom(comment),
		})
		s.Error(err)
	})
}

func (s *PaymentTestSuite) TestGetPayment() {
	workflowClientMock := &mocks.Client{}
	service := payment.NewService(s.entClient, workflowClientMock)

	paymentID := int64(1)
	s.entClient.Payment.
		Create().
		SetID(paymentID).
		SetStatus(entpayment.StatusPENDING).
		ExecX(context.TODO())
	s.entClient.Review.
		Create().
		SetPaymentID(paymentID).
		SetEvent(review.EventAPPROVE).
		SetReviewerID("alice@abc.com").
		ExecX(context.TODO())

	entityPayment, err := service.GetPayment(context.TODO(), &approvalservice.GetPaymentRequest{
		PaymentID: paymentID,
	})
	s.NoError(err)
	s.Equal(entpayment.StatusPENDING, entityPayment.Status)
	s.Len(entityPayment.Edges.Reviews, 1)
	s.Equal("alice@abc.com", entityPayment.Edges.Reviews[0].ReviewerID)
	s.Equal(review.EventAPPROVE, entityPayment.Edges.Reviews[0].Event)
}

func (s *PaymentTestSuite) TestListPayments() {
	workflowClientMock := &mocks.Client{}
	service := payment.NewService(s.entClient, workflowClientMock)

	s.entClient.Payment.
		Create().
		SetID(1).
		SetStatus(entpayment.StatusPENDING).
		SetCreateTime(time.Date(1970, 1, 1, 1, 1, 0, 0, time.UTC)).
		ExecX(context.TODO())
	s.entClient.Review.
		Create().
		SetPaymentID(1).
		SetEvent(review.EventAPPROVE).
		SetReviewerID("alice@abc.com").
		ExecX(context.TODO())
	s.entClient.Payment.
		Create().
		SetID(2).
		SetStatus(entpayment.StatusAPPROVED).
		SetCreateTime(time.Date(1970, 1, 1, 1, 2, 0, 0, time.UTC)).
		ExecX(context.TODO())
	s.entClient.Payment.
		Create().
		SetID(3).
		SetStatus(entpayment.StatusREJECTED).
		SetCreateTime(time.Date(1970, 1, 1, 1, 3, 0, 0, time.UTC)).
		ExecX(context.TODO())
	s.entClient.Payment.
		Create().
		SetID(4).
		SetStatus(entpayment.StatusPENDING).
		SetCreateTime(time.Date(1970, 1, 1, 1, 4, 0, 0, time.UTC)).
		ExecX(context.TODO())

	s.Run("No filtering", func() {
		resp, err := service.ListPayments(context.TODO(), &approvalservice.ListPaymentsRequest{})
		s.NoError(err)
		s.False(resp.HasNext)

		entityPayments := resp.Payments
		s.Len(entityPayments, 4)
		s.Equal(int64(1), entityPayments[3].ID)
		s.Len(entityPayments[3].Edges.Reviews, 1)
		s.Equal("alice@abc.com", entityPayments[3].Edges.Reviews[0].ReviewerID)
		s.Equal(review.EventAPPROVE, entityPayments[3].Edges.Reviews[0].Event)
		s.Equal(int64(2), entityPayments[2].ID)
		s.Equal(int64(3), entityPayments[1].ID)
		s.Equal(int64(4), entityPayments[0].ID)
	})

	s.Run("Filter by approved payment", func() {
		paymentStatuses := []entpayment.Status{entpayment.StatusAPPROVED}

		resp, err := service.ListPayments(context.TODO(), &approvalservice.ListPaymentsRequest{
			PaymentStatuses: &paymentStatuses,
		})
		s.NoError(err)
		s.False(resp.HasNext)

		entityPayments := resp.Payments
		s.Len(entityPayments, 1)
		s.Equal(int64(2), entityPayments[0].ID)
	})

	s.Run("Filter by pending payment and before create time", func() {
		paymentStatuses := []entpayment.Status{entpayment.StatusAPPROVED}
		beforeCreateTime := time.Date(1970, 1, 1, 1, 4, 0, 0, time.UTC)

		resp, err := service.ListPayments(context.TODO(), &approvalservice.ListPaymentsRequest{
			BeforeCreateTime: null.TimeFrom(beforeCreateTime),
			PaymentStatuses:  &paymentStatuses,
		})
		s.NoError(err)
		s.False(resp.HasNext)

		entityPayments := resp.Payments
		s.Len(entityPayments, 1)
		s.Equal(int64(2), entityPayments[0].ID)
	})
}

func (s *PaymentTestSuite) TestListPaymentsWithLargeDataset() {
	workflowClientMock := &mocks.Client{}
	service := payment.NewService(s.entClient, workflowClientMock)

	paymentIDsWithRejectedStatus := map[int64]bool{
		103: true,
		104: true,
		105: true,
		110: true,
		111: true,
		120: true,
		135: true,
		147: true,
		152: true,
	}

	for i := 0; i < 60; i++ {
		id := int64(i) + 100
		status := entpayment.StatusPENDING

		if paymentIDsWithRejectedStatus[id] {
			status = entpayment.StatusREJECTED
		}

		s.entClient.Payment.
			Create().
			SetID(id).
			SetStatus(status).
			SetCreateTime(time.Date(1970, 2, 1, 1, 1, i, 0, time.UTC)).
			ExecX(context.TODO())
	}

	s.Run("Query with page size", func() {
		paymentStatuses := []entpayment.Status{entpayment.StatusREJECTED}

		resp, err := service.ListPayments(context.TODO(), &approvalservice.ListPaymentsRequest{
			PageSize:        null.IntFrom(2),
			PaymentStatuses: &paymentStatuses,
		})
		s.NoError(err)
		s.True(resp.HasNext)

		entityPayments := resp.Payments
		s.Len(entityPayments, 2)
		s.Equal(int64(152), entityPayments[0].ID)
		s.Equal(int64(147), entityPayments[1].ID)
	})

	s.Run("Query with before create time and id", func() {
		paymentStatuses := []entpayment.Status{entpayment.StatusREJECTED}

		resp, err := service.ListPayments(context.TODO(), &approvalservice.ListPaymentsRequest{
			BeforeCreateTime:  null.TimeFrom(time.Date(1970, 2, 1, 1, 1, 47, 0, time.UTC)),
			LessThanPaymentID: null.IntFrom(147),
			PageSize:          null.IntFrom(1),

			PaymentStatuses: &paymentStatuses,
		})
		s.NoError(err)
		s.True(resp.HasNext)

		entityPayments := resp.Payments
		s.Len(entityPayments, 1)
		s.Equal(int64(135), entityPayments[0].ID)
	})

	s.Run("Check has next payment", func() {
		paymentStatuses := []entpayment.Status{entpayment.StatusREJECTED}

		resp, err := service.ListPayments(context.TODO(), &approvalservice.ListPaymentsRequest{
			BeforeCreateTime:  null.TimeFrom(time.Date(1970, 2, 1, 1, 1, 5, 0, time.UTC)),
			LessThanPaymentID: null.IntFrom(105),
			PageSize:          null.IntFrom(2),

			PaymentStatuses: &paymentStatuses,
		})
		s.NoError(err)
		s.False(resp.HasNext)

		entityPayments := resp.Payments
		s.Len(entityPayments, 2)
		s.Equal(int64(104), entityPayments[0].ID)
		s.Equal(int64(103), entityPayments[1].ID)
	})
}
