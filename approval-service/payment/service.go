package payment

import (
	"context"
	"fmt"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/ent"
	"github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
	"go.temporal.io/sdk/client"
)

type Service struct {
	entClient      *ent.Client
	workflowClient client.Client
}

var _ approvalservice.PaymentService = (*Service)(nil)

func NewService(
	entClient *ent.Client,
	workflowClient client.Client,
) *Service {
	return &Service{
		entClient:      entClient,
		workflowClient: workflowClient,
	}
}

func (s *Service) CreatePayment(
	ctx context.Context,
	request *approvalservice.CreatePaymentRequest,
) (*ent.Payment, error) {
	entityPayment, err := createPayment(ctx, s.entClient.Payment, request)
	if err != nil {
		return nil, err
	}

	// TODO: fix the workflow implementation? if application crash before this, no workflow is being created
	if _, err := s.workflowClient.ExecuteWorkflow(
		ctx,
		getStartApprovalWorkflowOptions(entityPayment.ID),
		ApprovalWorkflow,
		&ApprovalWorkflowRequest{
			PaymentID: entityPayment.ID,
		},
	); err != nil {
		return nil, fmt.Errorf("execute payment approval workflow: %w", err)
	}

	return entityPayment, nil
}

func (s *Service) CreatePaymentReview(
	ctx context.Context,
	request *approvalservice.CreatePaymentReviewRequest,
) error {
	paymentStatus, err := getPaymentStatusByID(ctx, s.entClient.Payment, request.PaymentID)
	if err != nil {
		return err
	}

	if paymentStatus != payment.StatusPENDING {
		return errNonPendingStatus
	}

	// use SignalWithStartWorkflow to prevent the workflow isn't being created before
	switch request.ReviewEvent {
	case review.EventAPPROVE:
		if _, err := s.workflowClient.SignalWithStartWorkflow(
			ctx,
			getApprovalWorkflowID(request.PaymentID),
			SignalChannelApprovePaymentChannel,
			&ApprovePaymentSignal{
				ReviewerID: request.ReviewerID,
				Comment:    request.Comment,
			},
			getStartApprovalWorkflowOptions(request.PaymentID),
			ApprovalWorkflow,
			&ApprovalWorkflowRequest{
				PaymentID: request.PaymentID,
			},
		); err != nil {
			return fmt.Errorf("signal approve payment: %w", err)
		}
	case review.EventREJECT:
		if _, err := s.workflowClient.SignalWithStartWorkflow(
			ctx,
			getApprovalWorkflowID(request.PaymentID),
			SignalChannelRejectPaymentChannel,
			&RejectPaymentSignal{
				ReviewerID: request.ReviewerID,
				Comment:    request.Comment,
			},
			getStartApprovalWorkflowOptions(request.PaymentID),
			ApprovalWorkflow,
			&ApprovalWorkflowRequest{
				PaymentID: request.PaymentID,
			},
		); err != nil {
			return fmt.Errorf("signal reject payment: %w", err)
		}
	default:
		return errUnknownReviewEvent
	}

	return nil
}

func (s *Service) GetPayment(
	ctx context.Context,
	request *approvalservice.GetPaymentRequest,
) (*ent.Payment, error) {
	entityPayment, err := getPaymentByID(ctx, s.entClient.Payment, request.PaymentID)
	if err != nil {
		return nil, err
	}

	return entityPayment, nil
}

func (s *Service) ListPayments(
	ctx context.Context,
	request *approvalservice.ListPaymentsRequest,
) (*approvalservice.ListPaymentsResponse, error) {
	payments, err := listPayments(ctx, s.entClient.Payment, request)
	if err != nil {
		return nil, err
	}

	return payments, nil
}
