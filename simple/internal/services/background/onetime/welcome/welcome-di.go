package welcome

import (
	"reflect"

	contractsBackgroundtasks "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/backgroundtasks"
	servicesLogger "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/logger"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/rs/zerolog/log"
)

// AddOneTimeWelcomeJobProvider adds service to the DI container
func AddOneTimeWelcomeJobProvider(builder *di.Builder) {
	log.Info().
		Msg("IoC: AddOneTimeWelcomeJobProvider")
	types := di.NewTypeSet()
	types.Add(contractsBackgroundtasks.ReflectTypeIJobsProvider)

	builder.Add(di.Def{
		Scope:            di.App,
		ImplementedTypes: types,
		Type:             reflect.TypeOf(&service{}),
		Build: func(ctn di.Container) (interface{}, error) {
			obj := &service{
				Logger: servicesLogger.GetSingletonLoggerFromContainer(ctn),
			}

			return obj, nil
		},
		Close: func(obj interface{}) error {

			return nil
		},
	})
}
