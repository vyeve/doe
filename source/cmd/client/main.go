package main

import (
	"context"

	server "doe/source/client"
	"doe/source/data"
	"doe/source/logger"

	"go.uber.org/fx"
)

func main() {
	ctx := context.Background()
	var (
		log logger.Logger
		srv server.Server
	)
	app := fx.New(
		logger.Module,
		data.Module,
		server.Module,
		fx.Populate(&log),
		fx.Populate(&srv),
	)
	defer app.Stop(ctx) // nolint: errcheck

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	if err := srv.Serve(); err != nil {
		log.Error(err)
	}
}
