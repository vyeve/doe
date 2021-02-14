package endpoint

import (
	"doe/source/logger"

	"go.uber.org/fx"
)

const (
	CommunicationURLEnvKey = "COMMUNICATION_URL"
	HostURLEnvKey          = "SERVER_URL"
)

type (
	ServerParams struct {
		fx.In

		Logger logger.Logger
	}

	RouterParams struct {
		fx.In

		Logger logger.Logger
		RPC    ServiceInterface
	}
)
