package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/Adopten123/go-todoapp-1/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux *http.ServeMux
	cfg Config
	log *core_logger.Logger
}

func NewHTTPServer(
	cfg Config,
	log *core_logger.Logger,
) *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(),
		cfg: cfg,
		log: log,
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    h.cfg.Addr,
		Handler: h.mux,
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)

		h.log.Warn("starting HTTP server", zap.String("addr", h.cfg.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), h.cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		h.log.Info("HTTP server stopped")
	}

	return nil
}
