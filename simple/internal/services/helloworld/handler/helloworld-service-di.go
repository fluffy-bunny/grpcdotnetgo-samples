package handler

import (
	"reflect"

	pb "github.com/fluffy-bunny/grpcdotnetgo-samples/contracts/simple/helloworld"
	servicesConfig "github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal/services/config"
	"github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contextaccessor "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/contextaccessor"
	servicesLogger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/rs/zerolog/log"
)

// AddGreeterService adds service to the DI container
func AddGreeterService(builder *di.Builder) {
	log.Info().
		Msg("IoC: AddGreeterService")
	types := di.NewTypeSet()
	types.Add(pb.TypeIGreeterService)

	builder.Add(di.Def{
		Scope:            di.Request,
		ImplementedTypes: types,
		Type:             reflect.TypeOf(&Service{}),
		Build: func(ctn di.Container) (interface{}, error) {
			return &Service{
				config:          servicesConfig.GetConfigFromContainer(ctn),
				ContextAccessor: contextaccessor.GetIContextAccessorFromContainer(ctn),
				ClaimsPrincipal: claimsprincipal.GetIClaimsPrincipalFromContainer(ctn),
				Logger:          servicesLogger.GetILoggerFromContainer(ctn),
			}, nil
		},
	})
}

// AddGreeter2Service adds service to the DI container
func AddGreeter2Service(builder *di.Builder) {
	log.Info().
		Msg("IoC: AddGreeter2Service")
	types := di.NewTypeSet()
	types.Add(pb.TypeIGreeter2Service)

	builder.Add(di.Def{
		Scope:            di.Request,
		ImplementedTypes: types,
		Type:             reflect.TypeOf(&Service2{}),
		Build: func(ctn di.Container) (interface{}, error) {
			return &Service2{
				config:          servicesConfig.GetConfigFromContainer(ctn),
				ContextAccessor: contextaccessor.GetIContextAccessorFromContainer(ctn),
				ClaimsPrincipal: claimsprincipal.GetIClaimsPrincipalFromContainer(ctn),
				Logger:          servicesLogger.GetILoggerFromContainer(ctn),
			}, nil
		},
	})
}
