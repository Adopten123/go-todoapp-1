package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/Adopten123/go-todoapp-1/internal/core/logger"
	"github.com/Adopten123/go-todoapp-1/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Adopten123/go-todoapp-1/internal/core/transport/http/middleware"
	core_http_server "github.com/Adopten123/go-todoapp-1/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/Adopten123/go-todoapp-1/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Adopten123/go-todoapp-1/internal/features/tasks/service"
	tasks_transport_http "github.com/Adopten123/go-todoapp-1/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Adopten123/go-todoapp-1/internal/features/users/repository/postgres"
	users_service "github.com/Adopten123/go-todoapp-1/internal/features/users/service"
	users_transport_http "github.com/Adopten123/go-todoapp-1/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	log, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Debug("application time zone", zap.Any("zone", timeZone))

	log.Debug("initializing postgres connection pool")

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		log.Fatal("failed to init postgres connection pool", zap.Error(err))
	}

	log.Debug("initializing feature", zap.String("feature", "users"))
	usersRepo := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepo)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	log.Debug("initializing feature", zap.String("feature", "tasks"))

	tasksRepo := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepo)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	log.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		log,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(log),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouters(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRouters(tasksTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}
