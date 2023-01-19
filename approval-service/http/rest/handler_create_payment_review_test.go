package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/api"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"
)

func (s *RESTTestSuite) TestCreatePaymentReviewHandler() {
	s.Run("Success", func() {
		s.paymentServiceMock.
			EXPECT().
			CreatePaymentReview(mock.Anything, &approvalservice.CreatePaymentReviewRequest{
				PaymentID:   1,
				ReviewEvent: review.EventAPPROVE,
				ReviewerID:  "bob@abc.com",
				Comment:     null.StringFrom("some text"),
			}).
			Return(
				nil,
			).
			Once()

		inputRequest := api.CreateReviewRequest{
			Event:      api.APPROVE,
			ReviewerId: "bob@abc.com",
			Comment:    null.StringFrom("some text").Ptr(),
		}
		inputB, err := json.Marshal(inputRequest)
		s.NoError(err)

		request := httptest.NewRequest(http.MethodPost, "/payments/1/reviews", bytes.NewBuffer(inputB))
		request.Header.Add("Content-Type", "application/json")
		response, _ := testHTTPHandler(s.server, request)
		s.Equal(http.StatusNoContent, response.StatusCode)
	})

	s.Run("Not found", func() {
		s.paymentServiceMock.
			EXPECT().
			CreatePaymentReview(mock.Anything, &approvalservice.CreatePaymentReviewRequest{
				PaymentID:   2,
				ReviewEvent: review.EventREJECT,
				ReviewerID:  "tom@abc.com",
				Comment:     null.StringFrom("other text"),
			}).
			Return(
				approvalservice.ErrNotFound,
			).
			Once()

		inputRequest := api.CreateReviewRequest{
			Event:      api.REJECT,
			ReviewerId: "tom@abc.com",
			Comment:    null.StringFrom("other text").Ptr(),
		}
		inputB, err := json.Marshal(inputRequest)
		s.NoError(err)

		request := httptest.NewRequest(http.MethodPost, "/payments/2/reviews", bytes.NewBuffer(inputB))
		request.Header.Add("Content-Type", "application/json")
		response, _ := testHTTPHandler(s.server, request)
		s.Equal(http.StatusNotFound, response.StatusCode)
	})
}
