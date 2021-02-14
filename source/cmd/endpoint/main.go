package main

import (
	"context"

	"doe/source/endpoint"
	"doe/source/logger"

	"go.uber.org/fx"
)

func main() {
	ctx := context.Background()
	var (
		log logger.Logger
		srv endpoint.Router
	)
	app := fx.New(
		logger.Module,
		endpoint.Module,
		fx.Populate(&log),
		fx.Populate(&srv),
	)
	defer app.Stop(ctx) // nolint: errcheck
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error(err)
	}
}
