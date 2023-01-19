package payment

import (
	"context"
	"fmt"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/ent"
	"github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
	"github.com/lht102/workflow-playground/approval-service/entutil"
	"gopkg.in/guregu/null.v4"
)

type GetPaymentStatusActivityRequest struct {
	PaymentID int64
}

type ApprovePaymentActivityRequest struct {
	PaymentID  int64
	ReviewerID string
	Comment    null.String
}

type RejectPaymentActivityRequest struct {
	PaymentID  int64
	ReviewerID string
	Comment    null.String
}

type DeactivatePaymentActivityRequest struct {
	PaymentID int64
}

type Activities struct {
	entClient *ent.Client
}

func NewActivities(entClient *ent.Client) *Activities {
	return &Activities{
		entClient: entClient,
	}
}

func (a *Activities) GetPaymentStatusActivity(
	ctx context.Context,
	request *GetPaymentStatusActivityRequest,
) (payment.Status, error) {
	paymentStatus, err := getPaymentStatusByID(ctx, a.entClient.Payment, request.PaymentID)
	if err != nil {
		return "", err
	}

	return paymentStatus, nil
}

func (a *Activities) ApprovePaymentActivity(
	ctx context.Context,
	request *ApprovePaymentActivityRequest,
) error {
	if err := entutil.WithTx(ctx, a.entClient, func(tx *ent.Tx) error {
		status, err := getPaymentStatusByID(ctx, tx.Payment, request.PaymentID)
		if err != nil {
			return err
		}

		if status != payment.StatusPENDING {
			return errNonPendingStatus
		}

		if err := upsertPaymentReviewExec(
			ctx,
			tx.Review,
			&approvalservice.CreatePaymentReviewRequest{
				PaymentID:   request.PaymentID,
				ReviewEvent: review.EventAPPROVE,
				ReviewerID:  request.ReviewerID,
				Comment:     request.Comment,
			},
		); err != nil {
			return err
		}

		numOfApprovals, err := getNumberOfApprovalsByPaymentID(ctx, tx.Review, request.PaymentID)
		if err != nil {
			return err
		}

		if numOfApprovals < minimumNumOfApprovalsToCompletePayment {
			return nil
		}

		approved := payment.StatusAPPROVED
		if err := updatePayment(
			ctx,
			tx.Payment,
			&approvalservice.UpdatePaymentRequest{
				PaymentID:     request.PaymentID,
				PaymentStatus: &approved,
			},
		); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("with tx: %w", err)
	}

	return nil
}

func (a *Activities) RejectPaymentActivity(
	ctx context.Context,
	request *RejectPaymentActivityRequest,
) error {
	if err := entutil.WithTx(ctx, a.entClient, func(tx *ent.Tx) error {
		status, err := getPaymentStatusByID(ctx, tx.Payment, request.PaymentID)
		if err != nil {
			return err
		}

		if status != payment.StatusPENDING {
			return errNonPendingStatus
		}

		if err := upsertPaymentReviewExec(
			ctx,
			tx.Review,
			&approvalservice.CreatePaymentReviewRequest{
				PaymentID:   request.PaymentID,
				ReviewEvent: review.EventREJECT,
				ReviewerID:  request.ReviewerID,
				Comment:     request.Comment,
			},
		); err != nil {
			return err
		}

		rejected := payment.StatusREJECTED
		if err := updatePayment(
			ctx,
			tx.Payment,
			&approvalservice.UpdatePaymentRequest{
				PaymentID:     request.PaymentID,
				PaymentStatus: &rejected,
			},
		); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("with tx: %w", err)
	}

	return nil
}

func (a *Activities) DeactivatePaymentActivity(
	ctx context.Context,
	request *DeactivatePaymentActivityRequest,
) error {
	if err := entutil.WithTx(ctx, a.entClient, func(tx *ent.Tx) error {
		status, err := getPaymentStatusByID(ctx, tx.Payment, request.PaymentID)
		if err != nil {
			return err
		}

		if status != payment.StatusPENDING {
			return errNonPendingStatus
		}

		rejected := payment.StatusREJECTED
		if err := updatePayment(
			ctx,
			tx.Payment,
			&approvalservice.UpdatePaymentRequest{
				PaymentID:     request.PaymentID,
				PaymentStatus: &rejected,
				Remark:        null.StringFrom("This is no longer active"),
			},
		); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("with tx: %w", err)
	}

	return nil
}
