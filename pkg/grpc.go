package pkg

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"net"
)

type GrpcServer interface {
	Start() error
	Stop() error
}

type grpcServer struct {
	srv  *grpc.Server
	port int32
}

func NewGrpcServer(port int32, listeners ...func(*grpc.Server)) GrpcServer {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	for _, listener := range listeners {
		listener(srv)
	}

	grpc_health_v1.RegisterHealthServer(srv, healthService{})

	grpc_prometheus.Register(srv)

	return &grpcServer{
		srv:  srv,
		port: port,
	}
}

func (s *grpcServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen %d; %w", s.port, err)
	}
	go func() {
		if err := s.srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return nil
}

func (s *grpcServer) Stop() error {
	s.srv.GracefulStop()
	return nil
}

type healthService struct{}

func (healthService) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (healthService) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watching is not supported")
}
