package rest

import (
	"net/http"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/lht102/workflow-playground/approval-service/api"
	"go.uber.org/zap"
)

func (s *Server) routes() {
	r := s.router

	swagger, err := api.GetSwagger()
	if err != nil {
		s.logger.Error("Failed to setup API spec", zap.Error(err))
		return
	}

	swagger.Servers = nil

	r.Use(middleware.OapiRequestValidatorWithOptions(
		swagger,
		&middleware.Options{
			ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
				respondErr(w, message, statusCode, s.logger)
			},
		},
	))
	api.HandlerWithOptions(s, api.ChiServerOptions{
		BaseRouter: r,
	})
}
