package app

import (
	"context"

	"github.com/chapsuk/mserv"
	"github.com/chapsuk/worker"
	"github.com/im-kulikov/helium"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type app struct {
	dig.In

	Logger  *zap.Logger
	Servers mserv.Server
	Workers *worker.Group
}

func newApp(a app) helium.App {
	return &a
}

func (a *app) Run(ctx context.Context) error {
	a.Logger.Info("start servers")
	a.Servers.Start()

	a.Logger.Info("start workers")
	a.Workers.Run()

	a.Logger.Info("start application")
	<-ctx.Done()

	a.Logger.Info("stop servers")
	a.Servers.Stop()

	a.Logger.Info("stop workers")
	a.Workers.Stop()

	a.Logger.Info("stop application")
	return nil
}
