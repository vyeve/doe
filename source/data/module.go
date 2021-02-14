package data

import "go.uber.org/fx"

var Module = fx.Provide(
	NewSQLTx,
	NewRepository,
)
