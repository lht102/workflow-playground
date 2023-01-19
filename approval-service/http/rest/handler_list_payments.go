package rest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/api"
	entpayment "github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/entutil"
	"gopkg.in/guregu/null.v4"
)

type PageTokenParams struct {
	LessThanPaymentID null.Int
	BeforeCreateTime  null.Time
}

func (s *Server) ListPayments(w http.ResponseWriter, r *http.Request, params api.ListPaymentsParams) {
	ctx := r.Context()

	pageTokenParams, err := unmarshalPageTokenParams(null.StringFromPtr(params.PageToken).String)
	if err != nil {
		respondErr(w, err.Error(), http.StatusBadRequest, s.logger)
		return
	}

	var entityPaymentStatuses *[]entpayment.Status

	if params.Statuses != nil {
		entityStatuses := make([]entpayment.Status, 0, len(*params.Statuses))
		for _, s := range *params.Statuses {
			entityStatuses = append(entityStatuses, entpayment.Status(s))
		}

		entityPaymentStatuses = &entityStatuses
	}

	var beforeCreateTime null.Time

	if pageTokenParams.BeforeCreateTime.Valid {
		beforeCreateTime = pageTokenParams.BeforeCreateTime
	} else if params.BeforeCreateTimestamp != nil {
		beforeCreateTime = null.TimeFrom(time.Unix(*params.BeforeCreateTimestamp, 0).UTC())
	}

	var pageSize null.Int
	if params.PageSize != nil {
		pageSize = null.IntFrom(int64(*params.PageSize))
	}

	resp, err := s.paymentService.ListPayments(ctx, &approvalservice.ListPaymentsRequest{
		BeforeCreateTime:  beforeCreateTime,
		LessThanPaymentID: pageTokenParams.LessThanPaymentID,
		PageSize:          pageSize,
		PaymentStatuses:   entityPaymentStatuses,
	})
	if err != nil {
		respondErr(w, err.Error(), http.StatusInternalServerError, s.logger)
		return
	}

	var nextPageToken null.String

	if len(resp.Payments) != 0 && resp.HasNext {
		lastEntityPayment := resp.Payments[len(resp.Payments)-1]
		token, err := marshalPageTokenParams(&PageTokenParams{
			LessThanPaymentID: null.IntFrom(lastEntityPayment.ID),
			BeforeCreateTime:  null.TimeFrom(lastEntityPayment.CreateTime),
		})

		if err != nil {
			respondErr(w, err.Error(), http.StatusBadRequest, s.logger)
			return
		}

		nextPageToken = null.StringFrom(token)
	}

	respond(
		w,
		api.ListPaymentsResponse{
			Payments:      entutil.PaymentsToAPIPayments(resp.Payments),
			NextPageToken: nextPageToken.Ptr(),
		},
		http.StatusOK,
		s.logger,
	)
}

func unmarshalPageTokenParams(pageToken string) (*PageTokenParams, error) {
	if pageToken == "" {
		return &PageTokenParams{}, nil
	}

	b, err := base64.RawURLEncoding.DecodeString(pageToken)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var pageTokenParams PageTokenParams
	if err := json.Unmarshal(b, &pageTokenParams); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return &pageTokenParams, nil
}

func marshalPageTokenParams(pageTokenParams *PageTokenParams) (string, error) {
	b, err := json.Marshal(pageTokenParams)
	if err != nil {
		return "", fmt.Errorf("json unmarshal: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
