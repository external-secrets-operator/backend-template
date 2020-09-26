package pkg

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ActuatorServer interface {
	Start() error
	Stop() error
}

type actuatorServer struct {
	srv *http.Server
}

func NewActuatorServer(port int32) ActuatorServer {
	s := actuatorServer{}
	handler := s.buildRootHandler()
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	return &s
}

func (s *actuatorServer) Start() error {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return nil
}

func (s *actuatorServer) Stop() error {
	return s.srv.Shutdown(context.TODO())
}

func (s *actuatorServer) buildRootHandler() http.Handler {
	rootRouter := mux.NewRouter()
	rootRouter.HandleFunc("/health", s.handleHealthRequest).Methods(http.MethodGet)
	rootRouter.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})).Methods(http.MethodGet)
	rootRouter.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	return rootRouter
}

func (s *actuatorServer) handleHealthRequest(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprint(writer, `{"status": "SERVING"}`)
	if err != nil {
		logrus.WithField("err", err).Error("failed to write data to response")
	}
}
