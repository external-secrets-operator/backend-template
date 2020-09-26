package internal

import (
	"backend-template/generated/api"
	"backend-template/pkg"
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"strings"
)

type App interface {
	Start() error
	Stop() error
}

type app struct {
	config         Config
	actuatorServer pkg.ActuatorServer
	grpcServer     pkg.GrpcServer
}

var appName = strings.ToUpper("backend_template")

func NewApp() (App, error) {
	app := app{}
	if err := envconfig.Process(appName, &app.config); err != nil {
		return nil, fmt.Errorf("failed to populate config; %w", err)
	}

	logrus.WithField("config", app.config).Info("starting application")

	app.actuatorServer = pkg.NewActuatorServer(app.config.HttpPort)

	svc := NewBackendService()
	app.grpcServer = pkg.NewGrpcServer(app.config.GrpcPort, func(s *grpc.Server) {
		api.RegisterBackendService(s, &svc)
	})
	return &app, nil
}

func (a *app) Start() error {
	if err := a.actuatorServer.Start(); err != nil {
		return fmt.Errorf("failed to start actuator server; %w", err)
	}
	if err := a.grpcServer.Start(); err != nil {
		return fmt.Errorf("failed to grpc server; %w", err)
	}
	return nil
}

func (a *app) Stop() error {
	wg, _ := errgroup.WithContext(context.TODO())
	wg.Go(a.grpcServer.Stop)
	wg.Go(a.actuatorServer.Stop)
	return wg.Wait()
}

type Config struct {
	GrpcPort int32 `split_words:"true" default:"8080"`
	HttpPort int32 `split_words:"true" default:"8181"`
}
