package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Adopten123/go-todoapp-1/docs"
	core_logger "github.com/Adopten123/go-todoapp-1/internal/core/logger"
	core_http_middleware "github.com/Adopten123/go-todoapp-1/internal/core/transport/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux *http.ServeMux
	cfg Config
	log *core_logger.Logger

	middlewares []core_http_middleware.Middleware
}

func NewHTTPServer(
	cfg Config,
	log *core_logger.Logger,
	middlewares ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		cfg:         cfg,
		log:         log,
		middlewares: middlewares,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()),
		)
	}
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"GET /swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)

	s.mux.HandleFunc(
		"GET /swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middlewares...)

	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)

		s.log.Warn("starting HTTP server", zap.String("addr", s.cfg.Addr))

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
		s.log.Warn("shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		s.log.Info("HTTP server stopped")
	}

	return nil
}
