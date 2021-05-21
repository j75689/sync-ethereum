//+build wireinject

//The build tag makes sure the stub is not built in the final build.

package http

import (
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/delivery/http"
	"sync-ethereum/internal/wireset"

	"github.com/google/wire"
)

func Initialize(configPath string) (Application, error) {
	wire.Build(
		newApplication,
		config.NewConfig,
		wireset.InitLogger,
		http.NewHttpServer,
	)
	return Application{}, nil
}
