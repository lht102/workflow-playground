package payment

import (
	"fmt"
	"time"

	"github.com/lht102/workflow-playground/approval-service/ent"
	"github.com/lht102/workflow-playground/approval-service/ent/payment"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"gopkg.in/guregu/null.v4"
)

const (
	SignalChannelApprovePaymentChannel = "approve_payment_channel"
	SignalChannelRejectPaymentChannel  = "reject_payment_channel"

	taskQueuePayment = "payment_task_queue"
)

type ApprovePaymentSignal struct {
	ReviewerID string
	Comment    null.String
}

type RejectPaymentSignal struct {
	ReviewerID string
	Comment    null.String
}

type ApprovalWorkflowRequest struct {
	PaymentID int64
}

type WorkflowWorker struct {
	worker.Worker
}

func NewWorkflowWorker(
	entClient *ent.Client,
	workflowClient client.Client,
	opts worker.Options,
) *WorkflowWorker {
	internalWorker := worker.New(workflowClient, taskQueuePayment, opts)

	a := NewActivities(entClient)
	internalWorker.RegisterActivity(a.ApprovePaymentActivity)
	internalWorker.RegisterActivity(a.RejectPaymentActivity)
	internalWorker.RegisterActivity(a.GetPaymentStatusActivity)
	internalWorker.RegisterWorkflow(ApprovalWorkflow)

	return &WorkflowWorker{
		Worker: internalWorker,
	}
}

func ApprovalWorkflow(
	ctx workflow.Context,
	request *ApprovalWorkflowRequest,
) error {
	logger := workflow.GetLogger(ctx)

	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    3 * time.Minute,
		NonRetryableErrorTypes: []string{
			errNonPendingStatus.Error(),
			errUnknownReviewEvent.Error(),
		},
	}

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	approvePaymentChannel := workflow.GetSignalChannel(ctx, SignalChannelApprovePaymentChannel)
	rejectPaymentChannel := workflow.GetSignalChannel(ctx, SignalChannelRejectPaymentChannel)

	for {
		var (
			paymentStatus payment.Status
			a             *Activities
		)

		if err := workflow.
			ExecuteActivity(
				ctx,
				a.GetPaymentStatusActivity,
				GetPaymentStatusActivityRequest{
					PaymentID: request.PaymentID,
				},
			).
			Get(
				ctx,
				&paymentStatus,
			); err != nil {
			return fmt.Errorf("execute get payment activity: %w", err)
		}

		if paymentStatus != payment.StatusPENDING {
			break
		}

		selector := workflow.NewSelector(ctx)
		selector.AddReceive(approvePaymentChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal ApprovePaymentSignal
			c.Receive(ctx, &signal)

			var (
				a *Activities
			)
			if err := workflow.
				ExecuteActivity(
					ctx,
					a.ApprovePaymentActivity,
					&ApprovePaymentActivityRequest{
						PaymentID:  request.PaymentID,
						ReviewerID: signal.ReviewerID,
						Comment:    signal.Comment,
					},
				).
				Get(
					ctx,
					nil,
				); err != nil {
				logger.Error("Failed to approve payment")

				return
			}
		})
		selector.AddReceive(rejectPaymentChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal RejectPaymentSignal
			c.Receive(ctx, &signal)

			var (
				a *Activities
			)
			if err := workflow.
				ExecuteActivity(
					ctx,
					a.RejectPaymentActivity,
					&RejectPaymentActivityRequest{
						PaymentID:  request.PaymentID,
						ReviewerID: signal.ReviewerID,
						Comment:    signal.Comment,
					},
				).
				Get(
					ctx,
					nil,
				); err != nil {
				logger.Error("Failed to reject payment", "error", err)

				return
			}
		})
		selector.AddFuture(workflow.NewTimer(ctx, durationToDeactivatePaymentFromLastUpdateTime), func(f workflow.Future) {
			var (
				a *Activities
			)
			if err := workflow.
				ExecuteActivity(
					ctx,
					a.DeactivatePaymentActivity,
					&DeactivatePaymentActivityRequest{
						PaymentID: request.PaymentID,
					},
				).
				Get(
					ctx,
					nil,
				); err != nil {
				logger.Error("Failed to deactivate payment", "error", err)

				return
			}
		})

		selector.Select(ctx)
	}

	return nil
}

func getApprovalWorkflowID(paymentID int64) string {
	return fmt.Sprintf("payment-approval-workflow-%d", paymentID)
}

func getStartApprovalWorkflowOptions(paymentID int64) client.StartWorkflowOptions {
	return client.StartWorkflowOptions{
		ID:        getApprovalWorkflowID(paymentID),
		TaskQueue: taskQueuePayment,
	}
}
