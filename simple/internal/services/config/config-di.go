package config

import (
	"github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal"
	servicesGrpcDotNetGoConfig "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/config"
	di "github.com/fluffy-bunny/sarulabsdi"
)

// GetConfigFromContainer from the Container
func GetConfigFromContainer(ctn di.Container) *internal.Config {
	obj := servicesGrpcDotNetGoConfig.GetConfigFromContainer(ctn).(*internal.Config)
	return obj
}
