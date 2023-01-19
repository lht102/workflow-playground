package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/api"
	"github.com/lht102/workflow-playground/approval-service/ent"
	"github.com/stretchr/testify/mock"
)

func (s *RESTTestSuite) TestGetPaymentHandler() {
	s.Run("Success", func() {
		s.paymentServiceMock.
			EXPECT().
			GetPayment(mock.Anything, &approvalservice.GetPaymentRequest{
				PaymentID: 1,
			}).
			Return(
				&ent.Payment{
					ID: 1,
				},
				nil,
			).
			Once()

		request := httptest.NewRequest(http.MethodGet, "/payments/1", nil)
		response, b := testHTTPHandler(s.server, request)
		s.Equal(http.StatusOK, response.StatusCode)

		var actualResponse api.ListPaymentsResponse
		err := json.Unmarshal(b, &actualResponse)
		s.NoError(err)
	})

	s.Run("Not found", func() {
		s.paymentServiceMock.
			EXPECT().
			GetPayment(mock.Anything, &approvalservice.GetPaymentRequest{
				PaymentID: 2,
			}).
			Return(
				nil,
				approvalservice.ErrNotFound,
			).
			Once()

		request := httptest.NewRequest(http.MethodGet, "/payments/2", nil)
		response, b := testHTTPHandler(s.server, request)
		s.Equal(http.StatusNotFound, response.StatusCode)

		var actualResponse api.ErrorResponse
		err := json.Unmarshal(b, &actualResponse)
		s.NoError(err)
	})
}
