package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lht102/workflow-playground/approval-service/api"
	"go.uber.org/zap"
)

func respond(w http.ResponseWriter, data any, statusCode int, logger *zap.Logger) {
	if data == nil {
		w.WriteHeader(statusCode)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Failed to response with json:", zap.Error(err))
	}
}

func respondErr(w http.ResponseWriter, errMsg string, statusCode int, logger *zap.Logger) {
	respond(
		w,
		api.ErrorResponse{
			Code:    statusCode,
			Message: errMsg,
		},
		statusCode,
		logger,
	)
}

func decode(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	return nil
}
