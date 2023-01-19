package payment_test

import (
	"context"

	entpayment "github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
	"github.com/lht102/workflow-playground/approval-service/payment"
)

func (s *PaymentTestSuite) TestGetPaymentStatusActivity() {
	paymentID := int64(1)
	s.entClient.Payment.
		Create().
		SetID(paymentID).
		SetStatus(entpayment.StatusREJECTED).
		ExecX(context.TODO())

	a := payment.NewActivities(s.entClient)
	paymentStatus, err := a.GetPaymentStatusActivity(
		context.TODO(),
		&payment.GetPaymentStatusActivityRequest{
			PaymentID: paymentID,
		},
	)
	s.NoError(err)
	s.Equal(entpayment.StatusREJECTED, paymentStatus)
}

func (s *PaymentTestSuite) TestApprovePaymentActivity() {
	paymentID := int64(1)
	s.entClient.Payment.
		Create().
		SetID(paymentID).
		SetStatus(entpayment.StatusPENDING).
		ExecX(context.TODO())

	a := payment.NewActivities(s.entClient)

	err := a.ApprovePaymentActivity(
		context.TODO(),
		&payment.ApprovePaymentActivityRequest{
			PaymentID:  paymentID,
			ReviewerID: "alice@abc.com",
		},
	)
	s.NoError(err)

	review1, err := getReviewByPaymentIDReviewerID(s.entClient, paymentID, "alice@abc.com")
	s.NoError(err)
	s.Equal(review.EventAPPROVE, review1.Event)

	err = a.ApprovePaymentActivity(
		context.TODO(),
		&payment.ApprovePaymentActivityRequest{
			PaymentID:  paymentID,
			ReviewerID: "bob@abc.com",
		},
	)
	s.NoError(err)

	paymentStatus, err := getPaymentStatusByID(s.entClient, paymentID)
	s.NoError(err)
	s.Equal(entpayment.StatusAPPROVED, paymentStatus)
}

func (s *PaymentTestSuite) TestRejectPaymentActivity() {
	paymentID := int64(1)
	s.entClient.Payment.
		Create().
		SetID(paymentID).
		SetStatus(entpayment.StatusPENDING).
		ExecX(context.TODO())

	a := payment.NewActivities(s.entClient)

	err := a.RejectPaymentActivity(
		context.TODO(),
		&payment.RejectPaymentActivityRequest{
			PaymentID:  paymentID,
			ReviewerID: "alice@abc.com",
		},
	)
	s.NoError(err)

	review1, err := getReviewByPaymentIDReviewerID(s.entClient, paymentID, "alice@abc.com")
	s.NoError(err)
	s.Equal(review.EventREJECT, review1.Event)

	paymentStatus, err := getPaymentStatusByID(s.entClient, paymentID)
	s.NoError(err)
	s.Equal(entpayment.StatusREJECTED, paymentStatus)
}

func (s *PaymentTestSuite) TestDeactivatePaymentActivity() {
	paymentID := int64(1)
	s.entClient.Payment.
		Create().
		SetID(paymentID).
		SetStatus(entpayment.StatusPENDING).
		ExecX(context.TODO())

	a := payment.NewActivities(s.entClient)

	err := a.DeactivatePaymentActivity(
		context.TODO(),
		&payment.DeactivatePaymentActivityRequest{
			PaymentID: paymentID,
		},
	)
	s.NoError(err)

	entityPayment, err := s.entClient.Payment.Get(context.TODO(), paymentID)
	s.NoError(err)
	s.Equal(entpayment.StatusREJECTED, entityPayment.Status)
	s.NotNil(entityPayment.Remark)
	s.NotEmpty(*entityPayment.Remark)
}
