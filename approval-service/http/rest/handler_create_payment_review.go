package rest

import (
	"errors"
	"net/http"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/api"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
	"gopkg.in/guregu/null.v4"
)

func (s *Server) CreatePaymentReview(w http.ResponseWriter, r *http.Request, paymentID int64) {
	ctx := r.Context()

	var request api.CreateReviewRequest
	if err := decode(r, &request); err != nil {
		respondErr(w, err.Error(), http.StatusBadRequest, s.logger)
		return
	}

	if err := s.paymentService.CreatePaymentReview(
		ctx,
		&approvalservice.CreatePaymentReviewRequest{
			PaymentID:   paymentID,
			ReviewEvent: review.Event(request.Event),
			ReviewerID:  request.ReviewerId,
			Comment:     null.StringFromPtr(request.Comment),
		},
	); err != nil {
		if errors.Is(err, approvalservice.ErrNotFound) {
			respondErr(w, err.Error(), http.StatusNotFound, s.logger)
			return
		}

		respondErr(w, err.Error(), http.StatusInternalServerError, s.logger)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
