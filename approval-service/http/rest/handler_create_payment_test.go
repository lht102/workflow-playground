package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/api"
	"github.com/lht102/workflow-playground/approval-service/ent"
	"github.com/stretchr/testify/mock"
)

func (s *RESTTestSuite) TestCreatePaymentHandler() {
	s.Run("Success", func() {
		s.paymentServiceMock.
			EXPECT().
			CreatePayment(mock.Anything, &approvalservice.CreatePaymentRequest{}).
			Return(
				&ent.Payment{
					ID: 1,
				},
				nil,
			).
			Once()

		inputRequest := api.CreatePaymentRequest{}
		inputB, err := json.Marshal(inputRequest)
		s.NoError(err)

		request := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(inputB))
		request.Header.Add("Content-Type", "application/json")
		response, b := testHTTPHandler(s.server, request)
		s.Equal(http.StatusOK, response.StatusCode)

		var actualResponse api.Payment
		err = json.Unmarshal(b, &actualResponse)
		s.NoError(err)
	})

	s.Run("Internal server error", func() {
		s.paymentServiceMock.
			EXPECT().
			CreatePayment(mock.Anything, &approvalservice.CreatePaymentRequest{}).
			Return(
				nil,
				errors.New("some error"),
			).
			Once()

		inputRequest := api.CreatePaymentRequest{}
		inputB, err := json.Marshal(inputRequest)
		s.NoError(err)

		request := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(inputB))
		request.Header.Add("Content-Type", "application/json")
		response, b := testHTTPHandler(s.server, request)
		s.Equal(http.StatusInternalServerError, response.StatusCode)

		var actualResponse api.ErrorResponse
		err = json.Unmarshal(b, &actualResponse)
		s.NoError(err)
	})
}
