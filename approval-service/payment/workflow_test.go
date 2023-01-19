package payment_test

import (
	"context"
	"time"

	entpayment "github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/payment"
	"github.com/stretchr/testify/mock"
)

func (s *PaymentTestSuite) TestApprovePayment() {
	var a *payment.Activities

	s.env.
		OnActivity(a.GetPaymentStatusActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.GetPaymentStatusActivityRequest,
			) (entpayment.Status, error) {
				return entpayment.StatusPENDING, nil
			},
		).
		Once()

	s.env.
		OnActivity(a.ApprovePaymentActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.ApprovePaymentActivityRequest,
			) error {
				return nil
			},
		).
		Once()

	s.env.RegisterDelayedCallback(func() {
		signal := payment.ApprovePaymentSignal{
			ReviewerID: "alice@abc.com",
		}

		s.env.SignalWorkflow(payment.SignalChannelApprovePaymentChannel, signal)
	}, time.Millisecond)

	s.env.
		OnActivity(a.GetPaymentStatusActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.GetPaymentStatusActivityRequest,
			) (entpayment.Status, error) {
				return entpayment.StatusAPPROVED, nil
			},
		).
		After(2 * time.Millisecond).
		Once()

	s.env.ExecuteWorkflow(
		payment.ApprovalWorkflow,
		&payment.ApprovalWorkflowRequest{
			PaymentID: 1,
		},
	)

	s.True(s.env.IsWorkflowCompleted())
}

func (s *PaymentTestSuite) TestRejectPayment() {
	var a *payment.Activities

	s.env.
		OnActivity(a.GetPaymentStatusActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.GetPaymentStatusActivityRequest,
			) (entpayment.Status, error) {
				return entpayment.StatusPENDING, nil
			},
		).
		Once()

	s.env.
		OnActivity(a.RejectPaymentActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.RejectPaymentActivityRequest,
			) error {
				return nil
			},
		).
		Once()

	s.env.
		OnActivity(a.GetPaymentStatusActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.GetPaymentStatusActivityRequest,
			) (entpayment.Status, error) {
				return entpayment.StatusREJECTED, nil
			},
		).
		After(2 * time.Millisecond).
		Once()

	s.env.RegisterDelayedCallback(func() {
		signal := payment.RejectPaymentSignal{
			ReviewerID: "alice@abc.com",
		}

		s.env.SignalWorkflow(payment.SignalChannelRejectPaymentChannel, signal)
	}, time.Millisecond)

	s.env.ExecuteWorkflow(
		payment.ApprovalWorkflow,
		&payment.ApprovalWorkflowRequest{
			PaymentID: 1,
		},
	)

	s.True(s.env.IsWorkflowCompleted())
}

func (s *PaymentTestSuite) TestDeactivatePayment() {
	var a *payment.Activities

	s.env.
		OnActivity(a.GetPaymentStatusActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.GetPaymentStatusActivityRequest,
			) (entpayment.Status, error) {
				return entpayment.StatusPENDING, nil
			},
		).
		Once()

	s.env.
		OnActivity(a.DeactivatePaymentActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.DeactivatePaymentActivityRequest,
			) error {
				return nil
			},
		).
		Once()

	s.env.
		OnActivity(a.GetPaymentStatusActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.GetPaymentStatusActivityRequest,
			) (entpayment.Status, error) {
				return entpayment.StatusREJECTED, nil
			},
		).
		Once()

	s.env.ExecuteWorkflow(
		payment.ApprovalWorkflow,
		&payment.ApprovalWorkflowRequest{
			PaymentID: 1,
		},
	)

	s.True(s.env.IsWorkflowCompleted())
}

func (s *PaymentTestSuite) TestApprovedPayment() {
	var a *payment.Activities

	s.env.
		OnActivity(a.GetPaymentStatusActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.GetPaymentStatusActivityRequest,
			) (entpayment.Status, error) {
				return entpayment.StatusAPPROVED, nil
			},
		).
		Once()

	s.env.RegisterDelayedCallback(func() {
		signal := payment.RejectPaymentSignal{
			ReviewerID: "alice@abc.com",
		}

		s.env.SignalWorkflow(payment.SignalChannelRejectPaymentChannel, signal)
	}, time.Millisecond)

	s.env.ExecuteWorkflow(
		payment.ApprovalWorkflow,
		&payment.ApprovalWorkflowRequest{
			PaymentID: 1,
		},
	)

	s.True(s.env.IsWorkflowCompleted())
}

func (s *PaymentTestSuite) TestRejectedPayment() {
	var a *payment.Activities

	s.env.
		OnActivity(a.GetPaymentStatusActivity, mock.Anything, mock.Anything).
		Return(
			func(
				_ context.Context,
				request *payment.GetPaymentStatusActivityRequest,
			) (entpayment.Status, error) {
				return entpayment.StatusREJECTED, nil
			},
		).
		Once()

	s.env.RegisterDelayedCallback(func() {
		signal := payment.ApprovePaymentSignal{
			ReviewerID: "alice@abc.com",
		}

		s.env.SignalWorkflow(payment.SignalChannelApprovePaymentChannel, signal)
	}, time.Millisecond)

	s.env.ExecuteWorkflow(
		payment.ApprovalWorkflow,
		&payment.ApprovalWorkflowRequest{
			PaymentID: 1,
		},
	)

	s.True(s.env.IsWorkflowCompleted())
}
