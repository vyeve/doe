package domainserver

import (
	"doe/source/data"
	"doe/source/logger"

	"go.uber.org/fx"
)

const (
	HostURLEnvKey = "SERVER_URL"
)

const (
	insertLimit = 100
)

type Params struct {
	fx.In

	Logger logger.Logger
	Repo   data.Repository
}
