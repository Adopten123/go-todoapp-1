package core_http_response

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_logger "github.com/Adopten123/go-todoapp-1/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	w   http.ResponseWriter
}

func NewHTTPResponseHandler(
	log *core_logger.Logger,
	w http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		w:   w,
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", msg)

	h.log.Error(msg, zap.Error(err))
	h.w.WriteHeader(statusCode)

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	if err := json.NewEncoder(h.w).Encode(response); err != nil {
		h.log.Error("Write HTTP response", zap.Error(err))
	}
}
