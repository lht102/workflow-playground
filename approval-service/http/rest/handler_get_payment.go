package rest

import (
	"errors"
	"net/http"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/entutil"
)

func (s *Server) GetPaymentByID(w http.ResponseWriter, r *http.Request, paymentID int64) {
	ctx := r.Context()

	entityPayment, err := s.paymentService.GetPayment(ctx, &approvalservice.GetPaymentRequest{
		PaymentID: paymentID,
	})
	if err != nil {
		if errors.Is(err, approvalservice.ErrNotFound) {
			respondErr(w, err.Error(), http.StatusNotFound, s.logger)
			return
		}

		respondErr(w, err.Error(), http.StatusInternalServerError, s.logger)

		return
	}

	respond(w, entutil.PaymentToAPIPayment(entityPayment), http.StatusOK, s.logger)
}
