package rest_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	approvalservice "github.com/lht102/workflow-playground/approval-service"
	"github.com/lht102/workflow-playground/approval-service/api"
	"github.com/lht102/workflow-playground/approval-service/ent"
	"github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"
)

func (s *RESTTestSuite) TestListPaymentsHandler() {
	s.Run("Without filtering", func() {
		s.paymentServiceMock.
			EXPECT().
			ListPayments(
				mock.Anything,
				&approvalservice.ListPaymentsRequest{
					PageSize: null.IntFrom(50),
				},
			).
			Return(
				&approvalservice.ListPaymentsResponse{
					Payments: ent.Payments{
						&ent.Payment{
							ID: 1,
						},
					},
					HasNext: false,
				},
				nil,
			).
			Once()

		request := httptest.NewRequest(http.MethodGet, "/payments", nil)
		response, b := testHTTPHandler(s.server, request)
		s.Equal(http.StatusOK, response.StatusCode)

		var actualResponse api.ListPaymentsResponse
		err := json.Unmarshal(b, &actualResponse)
		s.NoError(err)
		s.Nil(actualResponse.NextPageToken)
	})

	s.Run("Filter by approved payment", func() {
		statuses := []payment.Status{
			payment.StatusAPPROVED,
		}

		s.paymentServiceMock.
			EXPECT().
			ListPayments(
				mock.Anything,
				&approvalservice.ListPaymentsRequest{
					PageSize:        null.IntFrom(50),
					PaymentStatuses: &statuses,
				},
			).
			Return(
				&approvalservice.ListPaymentsResponse{
					Payments: ent.Payments{
						&ent.Payment{
							ID: 1,
						},
					},
					HasNext: false,
				},
				nil,
			).
			Once()

		request := httptest.NewRequest(http.MethodGet, "/payments?statuses=APPROVED", nil)
		response, b := testHTTPHandler(s.server, request)
		s.Equal(http.StatusOK, response.StatusCode)

		var actualResponse api.ListPaymentsResponse
		err := json.Unmarshal(b, &actualResponse)
		s.NoError(err)
		s.Nil(actualResponse.NextPageToken)
	})

	s.Run("Query by before create timestamp", func() {
		beforeCreateTime := time.Date(2023, 1, 1, 1, 0, 0, 0, time.UTC)

		s.paymentServiceMock.
			EXPECT().
			ListPayments(
				mock.Anything,
				&approvalservice.ListPaymentsRequest{
					PageSize:         null.IntFrom(50),
					BeforeCreateTime: null.TimeFrom(beforeCreateTime),
				},
			).
			Return(
				&approvalservice.ListPaymentsResponse{
					Payments: ent.Payments{
						&ent.Payment{
							ID: 1,
						},
					},
					HasNext: false,
				},
				nil,
			).
			Once()

		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/payments?before_create_timestamp=%d", beforeCreateTime.Unix()),
			nil,
		)
		response, b := testHTTPHandler(s.server, request)
		s.Equal(http.StatusOK, response.StatusCode)

		var actualResponse api.ListPaymentsResponse
		err := json.Unmarshal(b, &actualResponse)
		s.NoError(err)
		s.Nil(actualResponse.NextPageToken)
	})

	s.Run("Check page token", func() {
		var nextPageToken string

		s.Run("Have next page token", func() {
			s.paymentServiceMock.
				EXPECT().
				ListPayments(
					mock.Anything,
					&approvalservice.ListPaymentsRequest{
						PageSize: null.IntFrom(1),
					},
				).
				Return(
					&approvalservice.ListPaymentsResponse{
						Payments: ent.Payments{
							&ent.Payment{
								ID:         100,
								CreateTime: time.Date(1970, 12, 1, 1, 0, 0, 0, time.UTC),
							},
						},
						HasNext: true,
					},
					nil,
				).
				Once()

			request := httptest.NewRequest(http.MethodGet, "/payments?page_size=1", nil)
			response, b := testHTTPHandler(s.server, request)
			s.Equal(http.StatusOK, response.StatusCode)

			var actualResponse api.ListPaymentsResponse
			err := json.Unmarshal(b, &actualResponse)
			s.NoError(err)
			s.NotNil(actualResponse.NextPageToken)
			nextPageToken = *actualResponse.NextPageToken
		})

		s.Run("Use next page token", func() {
			s.paymentServiceMock.
				EXPECT().
				ListPayments(
					mock.Anything,
					&approvalservice.ListPaymentsRequest{
						BeforeCreateTime:  null.TimeFrom(time.Date(1970, 12, 1, 1, 0, 0, 0, time.UTC)),
						LessThanPaymentID: null.IntFrom(100),
						PageSize:          null.IntFrom(1),
					},
				).
				Return(
					&approvalservice.ListPaymentsResponse{
						Payments: ent.Payments{
							&ent.Payment{
								ID: 1,
							},
						},
						HasNext: false,
					},
					nil,
				).
				Once()

			request := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/payments?page_token=%s&page_size=1", nextPageToken),
				nil,
			)
			response, b := testHTTPHandler(s.server, request)
			s.Equal(http.StatusOK, response.StatusCode)

			var actualResponse api.ListPaymentsResponse
			err := json.Unmarshal(b, &actualResponse)
			s.NoError(err)
			s.Nil(actualResponse.NextPageToken)
		})
	})
}
