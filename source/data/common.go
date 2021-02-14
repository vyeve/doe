package data

import (
	"doe/source/logger"

	"go.uber.org/fx"
)

type Params struct {
	fx.In

	LifeCycle fx.Lifecycle `optional:"true"`
	Logger    logger.Logger
	DBConn    SQLDb
}
