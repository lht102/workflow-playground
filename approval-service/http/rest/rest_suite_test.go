package rest_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lht102/workflow-playground/approval-service/http/rest"
	"github.com/lht102/workflow-playground/approval-service/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type RESTTestSuite struct {
	suite.Suite

	server             *rest.Server
	paymentServiceMock *mocks.PaymentService
}

func TestRest(t *testing.T) {
	paymentServiceMock := mocks.NewPaymentService(t)

	server := rest.NewServer(paymentServiceMock, 8081, zap.NewExample())

	s := &RESTTestSuite{
		server:             server,
		paymentServiceMock: paymentServiceMock,
	}

	suite.Run(t, s)
}

func testHTTPHandler(handler http.Handler, r *http.Request) (*http.Response, []byte) {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	resp := w.Result()
	b, _ := io.ReadAll(resp.Body)

	return resp, b
}
