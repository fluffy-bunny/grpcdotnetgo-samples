package transient

import (
	"reflect"

	servicesConfig "github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal/services/config"
	exampleServices "github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal/services/contracts"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/rs/zerolog/log"
)

var (
	rtService  = reflect.TypeOf(&Service{}).Elem()
	rtService2 = reflect.TypeOf(&Service2{}).Elem()
)

// GetTransientServiceFromContainer from the Container
func GetTransientServiceFromContainer(ctn di.Container) *Service {
	return ctn.GetByType(rtService).(*Service)
}

// AddTransientService adds service to the DI container
func AddTransientService(builder *di.Builder) {
	log.Info().
		Msg("IoC: AddTransientService")

	implementedTypes := di.NewTypeSet()
	implementedTypes.Add(exampleServices.ReflectTypeISomething)

	builder.Add(di.Def{
		Scope:            di.App,
		ImplementedTypes: implementedTypes,
		Type:             reflect.TypeOf(&Service{}),
		Unshared:         true,
		Build: func(ctn di.Container) (interface{}, error) {
			service := &Service{
				config: servicesConfig.GetConfigFromContainer(ctn),
			}
			return service, nil
		},
	})
}

// GetTransientService2FromContainer from the Container
func GetTransientService2FromContainer(ctn di.Container) *Service2 {
	return ctn.GetByType(rtService2).(*Service2)
}

// AddTransientService2 to the Container
func AddTransientService2(builder *di.Builder) {
	log.Info().
		Msg("IoC: AddTransientService2")

	implementedTypes := di.NewTypeSet()
	implementedTypes.Add(exampleServices.ReflectTypeISomething)

	builder.Add(di.Def{
		Scope:            di.App,
		ImplementedTypes: implementedTypes,
		Type:             reflect.TypeOf(&Service2{}),
		Unshared:         true,
		Build: func(ctn di.Container) (interface{}, error) {
			service := &Service2{
				config: servicesConfig.GetConfigFromContainer(ctn),
			}
			return service, nil
		},
	})
}
