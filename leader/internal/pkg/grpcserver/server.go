package grpcserver

import (
	"fmt"
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/kysre/TurtleMQ/leader/pkg/errors"
	"github.com/kysre/TurtleMQ/leader/pkg/leader"
	"github.com/kysre/TurtleMQ/leader/pkg/queue"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	listener net.Listener
	server   *grpc.Server
}

func New(logger *logrus.Logger, listenPort int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		return nil, errors.Wrap(err, "failed to listen")
	}

	logEntry := logger.WithFields(map[string]interface{}{
		"app": "queue",
	})

	interceptors := []grpc.UnaryServerInterceptor{
		grpclogrus.UnaryServerInterceptor(logEntry),
		errors.UnaryServerInterceptor,
		grpcprometheus.UnaryServerInterceptor,
		grpcrecovery.UnaryServerInterceptor(),
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))

	return &Server{
		listener: listener,
		server:   grpcServer,
	}, nil
}

func (s *Server) RegisterQueueServer(queueServicer queue.QueueServer) {
	queue.RegisterQueueServer(s.server, queueServicer)
}

func (s *Server) RegisterLeaderServer(leaderServicer leader.LeaderServer) {
	leader.RegisterLeaderServer(s.server, leaderServicer)
}

func (s *Server) Serve() error {
	if err := s.server.Serve(s.listener); err != nil {
		return errors.Wrap(err, "failed to serve")
	}
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
