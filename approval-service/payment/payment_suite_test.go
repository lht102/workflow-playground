package payment_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lht102/workflow-playground/approval-service/ent"
	"github.com/lht102/workflow-playground/approval-service/ent/enttest"
	"github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type PaymentTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env       *testsuite.TestWorkflowEnvironment
	entClient *ent.Client
}

func (s *PaymentTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	s.entClient = enttest.
		Open(s.T(), "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1").
		Debug()
}

func (s *PaymentTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
	s.entClient.Close()
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentTestSuite))
}

func getReviewByPaymentIDReviewerID(
	entClient *ent.Client,
	paymentID int64,
	reviewerID string,
) (*ent.Review, error) {
	return entClient.Review.
		Query().
		Where(
			review.And(
				review.PaymentIDEQ(paymentID),
				review.ReviewerIDEQ(reviewerID),
			),
		).
		Only(context.TODO())
}

func getPaymentStatusByID(
	entClient *ent.Client,
	paymentID int64,
) (payment.Status, error) {
	s, err := entClient.Payment.
		Query().
		Where(payment.IDEQ(paymentID)).
		Select(payment.FieldStatus).
		String(context.TODO())
	if err != nil {
		return "", err
	}

	return payment.Status(s), nil
}

func getPaymentByRequestID(
	entClient *ent.Client,
	requestID uuid.UUID,
) (*ent.Payment, error) {
	entityPayment, err := entClient.Payment.
		Query().
		Where(payment.RequestIDEQ(requestID)).
		Only(context.TODO())
	if err != nil {
		return nil, err
	}

	return entityPayment, nil
}
