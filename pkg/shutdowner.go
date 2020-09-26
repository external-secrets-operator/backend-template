package pkg

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func WaitForShutdown(listeners ...func() error) error {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	logrus.WithField("signal", sig).Info("shutdown signal received")

	wg, _ := errgroup.WithContext(context.TODO())
	for _, listener := range listeners {
		wg.Go(listener)
	}
	return wg.Wait()
}
