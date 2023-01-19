package entutil

import (
	"github.com/lht102/workflow-playground/approval-service/api"
	"github.com/lht102/workflow-playground/approval-service/ent"
)

func PaymentToAPIPayment(payment *ent.Payment) api.Payment {
	reviews := make([]api.Review, 0, len(payment.Edges.Reviews))

	for _, entityReview := range payment.Edges.Reviews {
		review := api.Review{
			Comment:    entityReview.Comment,
			CreateTime: entityReview.CreateTime,
			Event:      api.ReviewEvent(entityReview.Event),
			Id:         entityReview.ID,
			ReviewerId: entityReview.ReviewerID,
			UpdateTime: entityReview.UpdateTime,
		}
		reviews = append(reviews, review)
	}

	return api.Payment{
		CreateTime: payment.CreateTime,
		Id:         payment.ID,
		Remark:     payment.Remark,
		RequestId:  payment.RequestID,
		Reviews:    reviews,
		Status:     api.PaymentStatus(payment.Status),
		UpdateTime: payment.UpdateTime,
	}
}

func PaymentsToAPIPayments(payments ent.Payments) []api.Payment {
	apiPayments := make([]api.Payment, 0, len(payments))

	for _, payment := range payments {
		apiPayments = append(apiPayments, PaymentToAPIPayment(payment))
	}

	return apiPayments
}
