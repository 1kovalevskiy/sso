package grpcserver

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	notify chan error
	log    *slog.Logger
	port   int
}

func New(log *slog.Logger, port int, interceptors ...grpc.ServerOption) *Server {
	gRPCServer := grpc.NewServer(interceptors...)

	return &Server{
		server: gRPCServer,
		notify: make(chan error, 1),
		log:    log,
		port:   port,
	}
}

func (s *Server) Register(registrator func(*grpc.Server)) {
	registrator(s.server)
}

func (s *Server) Start() {
	const op = "grpcserver.Start"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.notify <- fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		s.log.Info("grpc server started", slog.String("addr", l.Addr().String()))
		s.notify <- s.server.Serve(l)
		close(s.notify)
	}()

}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() {
	s.server.GracefulStop()
}
