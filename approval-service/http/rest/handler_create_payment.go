package rest

import (
	"net/http"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/api"
	"github.com/lht102/workflow-playground/approval-service/entutil"
)

func (s *Server) CreatePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request api.CreatePaymentRequest
	if err := decode(r, &request); err != nil {
		respondErr(w, err.Error(), http.StatusBadRequest, s.logger)
		return
	}

	entityPayment, err := s.paymentService.CreatePayment(
		ctx,
		&approvalservice.CreatePaymentRequest{
			RequestID: request.RequestId,
		},
	)
	if err != nil {
		respondErr(w, err.Error(), http.StatusInternalServerError, s.logger)
		return
	}

	respond(w, entutil.PaymentToAPIPayment(entityPayment), http.StatusOK, s.logger)
}
