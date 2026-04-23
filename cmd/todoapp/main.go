package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/Adopten123/go-todoapp-1/internal/core/config"
	core_logger "github.com/Adopten123/go-todoapp-1/internal/core/logger"
	"github.com/Adopten123/go-todoapp-1/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Adopten123/go-todoapp-1/internal/core/transport/http/middleware"
	core_http_server "github.com/Adopten123/go-todoapp-1/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/Adopten123/go-todoapp-1/internal/features/statistics/repository/postgres"
	statistics_service "github.com/Adopten123/go-todoapp-1/internal/features/statistics/service"
	statistics_transport_http "github.com/Adopten123/go-todoapp-1/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/Adopten123/go-todoapp-1/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Adopten123/go-todoapp-1/internal/features/tasks/service"
	tasks_transport_http "github.com/Adopten123/go-todoapp-1/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Adopten123/go-todoapp-1/internal/features/users/repository/postgres"
	users_service "github.com/Adopten123/go-todoapp-1/internal/features/users/service"
	users_transport_http "github.com/Adopten123/go-todoapp-1/internal/features/users/transport/http"
	web_fs_repository "github.com/Adopten123/go-todoapp-1/internal/features/web/repository/file_system"
	web_service "github.com/Adopten123/go-todoapp-1/internal/features/web/service"
	web_transport_http "github.com/Adopten123/go-todoapp-1/internal/features/web/transport/http"
	"go.uber.org/zap"

	_ "github.com/Adopten123/go-todoapp-1/docs"
)

// @title 		Go Todo API
// @version 	1.0
// @description Todo Application REST-API scheme
// @host 		127.0.0.1:5050
// @BasePath 	/api/v1
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

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

	log.Debug("application time zone", zap.Any("zone", time.Local))

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

	log.Debug("initializing feature", zap.String("feature", "statistics"))

	statisticsRepo := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepo)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	log.Debug("initializing feature", zap.String("feature", "web"))

	webRepo := web_fs_repository.NewWebRepository()
	webService := web_service.NewWebService(webRepo)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	log.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		log,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(log),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouters(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRouters(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRouters(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)
	httpServer.RegisterRouters(webTransportHTTP.Routes()...)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}
