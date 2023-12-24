package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/1kovalevskiy/sso/config"
	error_ "github.com/1kovalevskiy/sso/internal/error"
	authgrpc "github.com/1kovalevskiy/sso/internal/grpc"
	"github.com/1kovalevskiy/sso/internal/interceptor"
	"github.com/1kovalevskiy/sso/internal/usecase"
	repo "github.com/1kovalevskiy/sso/internal/usecase/repo_sqlite"
	"github.com/1kovalevskiy/sso/pkg/grpcserver"
	"github.com/1kovalevskiy/sso/pkg/logger"
	sqlite_ "github.com/1kovalevskiy/sso/pkg/sqlite"
)

func Run(cfg *config.Config) {
	const op = "internal - app - Run"
	l := logger.New("local")

	sqlite, err := sqlite_.New(cfg.SQL.URL, cfg.SQL.Timeout)
	if err != nil {
		l.Error(op+" - sql.New", error_.Err(err))
		return
	}
	defer sqlite.Close()

	interceptor := interceptor.NewInterceptor(l)

	server := grpcserver.New(l, cfg.GRPC.Port, interceptor)

	authUseCase := usecase.New(l, repo.New(sqlite))

	server.Register(authgrpc.New(authUseCase))

	server.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info(op+" - signal: " + s.String())
	case err = <-server.Notify():
		l.Error(op+" - grpcServer.Notify:", error_.Err(err))
	}

	server.Shutdown()
}
