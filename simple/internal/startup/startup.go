package startup

import (
	"fmt"
	"os"
	"path/filepath"

	pb "github.com/fluffy-bunny/grpcdotnetgo-samples/contracts/simple/helloworld"
	"github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal"
	backgroundCounterService "github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal/services/background/cron/counter"
	backgroundWelcomeService "github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal/services/background/onetime/welcome"
	handlerGreeterService "github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal/services/helloworld/handler"
	singletonService "github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal/services/singleton"
	transientService "github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal/services/transient"
	"github.com/fluffy-bunny/grpcdotnetgo/pkg/auth/oauth2"
	claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	coreContracts "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/core"
	middleware_dicontext "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/dicontext/middleware"
	middleware_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/logger"
	middleware_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/oidc"
	middleware_grpc_recovery "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/recovery"
	grpcDIProtoError "github.com/fluffy-bunny/grpcdotnetgo/pkg/proto/error"
	mockoidcservice "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/test/mockoidcservice"
	di "github.com/fluffy-bunny/sarulabsdi"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	_ "github.com/jnewmano/grpc-json-proxy/codec" // justified
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getConfigPath() string {
	var configPath string
	_, err := os.Stat("../etc/config")
	if !os.IsNotExist(err) {
		configPath, _ = filepath.Abs("../etc/config")
		log.Info().Str("path", configPath).Msg("Configuration Root Folder")
	}
	return configPath
}

// Startup ...
type Startup struct {
	MockOIDCService interface{}
	ConfigOptions   *coreContracts.ConfigOptions
	RootContainer   di.Container
}

// NewStartup ...
func NewStartup() coreContracts.IStartup {
	startup := &Startup{}
	startup.ctor()
	return startup
}

// GetStartupManifest wrapper
func (s *Startup) GetStartupManifest() coreContracts.StartupManifest {
	return coreContracts.StartupManifest{
		Name:    "hello",
		Version: "test.1",
	}
}

// OnPreServerStartup wrapper
func (s *Startup) OnPreServerStartup() error {
	return nil
}

// OnPostServerShutdown Wrapper
func (s *Startup) OnPostServerShutdown() {}

func (s *Startup) ctor() {
	s.ConfigOptions = &coreContracts.ConfigOptions{
		Destination: &internal.Config{},
		RootConfig:  internal.ConfigDefaultYaml,
		ConfigPath:  getConfigPath(),
	}
}

// GetConfigOptions ...
func (s *Startup) GetConfigOptions() *coreContracts.ConfigOptions {
	return s.ConfigOptions
}

// SetRootContainer ...
func (s *Startup) SetRootContainer(container di.Container) {
	s.RootContainer = container
}

// GetPort ...
func (s *Startup) GetPort() int {
	config := s.ConfigOptions.Destination.(*internal.Config)
	return config.Example.GRPCPort
}

// ConfigureServices ...
func (s *Startup) ConfigureServices(builder *di.Builder) {
	// this is how  you get your config before you register your services
	config := s.ConfigOptions.Destination.(*internal.Config)

	var mm = make(map[string]middleware_oidc.EntryPointConfig)

	for k, v := range config.Example.OIDCConfig.EntryPoints {
		mm[k] = v
	}
	for k, v := range mm {
		delete(config.Example.OIDCConfig.EntryPoints, k)
		config.Example.OIDCConfig.EntryPoints[v.FullMethodName] = v
	}
	handlerGreeterService.AddGreeterService(builder)
	handlerGreeterService.AddGreeter2Service(builder)

	singletonService.AddSingletonService(builder)

	transientService.AddTransientService(builder)
	if config.Example.EnableTransient2 {
		transientService.AddTransientService2(builder)
	}

	backgroundCounterService.AddCronCounterJobProvider(builder)
	backgroundWelcomeService.AddOneTimeWelcomeJobProvider(builder)

	mockoidcservice.AddMockOIDCService(builder)

	middleware_oidc.AddOIDCConfigAccessor(builder, config)
	//	backgroundOidcService.AddCronOidcJobProvider(builder)
	//	services_oidc.AddOIDCAuthHandler(builder)
}

// Configure ...
func (s *Startup) Configure(unaryServerInterceptorBuilder coreContracts.IUnaryServerInterceptorBuilder) {
	// this is how  you get your config before you register your services
	config := s.ConfigOptions.Destination.(*internal.Config)

	grpcFuncAuthConfig := oauth2.NewGrpcFuncAuthConfig(config.Example.OIDCConfig.Authority,
		"bearer", 5)
	for _, v := range config.Example.OIDCConfig.EntryPoints {
		methodClaims := oauth2.MethodClaims{
			OR:  []claimsprincipal.Claim{},
			AND: []claimsprincipal.Claim{},
		}

		for _, vv := range v.ClaimsConfig.AND {
			methodClaims.AND = append(methodClaims.AND, claimsprincipal.Claim{
				Type:  vv.Type,
				Value: vv.Value,
			})
		}

		for _, vv := range v.ClaimsConfig.OR {
			methodClaims.OR = append(methodClaims.OR, claimsprincipal.Claim{
				Type:  vv.Type,
				Value: vv.Value,
			})
		}

		grpcFuncAuthConfig.FullMethodNameToClaims[v.FullMethodName] = methodClaims
	}
	oidcContext, err := oauth2.BuildOpenIdConnectContext(grpcFuncAuthConfig)
	if err != nil {
		panic(err)
	}

	//var recoveryFunc middleware_grpc_recovery.RecoveryHandlerFunc
	recoveryOpts := []middleware_grpc_recovery.Option{
		middleware_grpc_recovery.WithRecoveryHandlerUnary(recoveryUnaryFunc),
	}
	unaryServerInterceptorBuilder.Use(grpc_ctxtags.UnaryServerInterceptor())
	unaryServerInterceptorBuilder.Use(middleware_logger.EnsureContextLoggingUnaryServerInterceptor())
	unaryServerInterceptorBuilder.Use(middleware_logger.EnsureCorrelationIDUnaryServerInterceptor())
	unaryServerInterceptorBuilder.Use(middleware_dicontext.UnaryServerInterceptor(s.RootContainer))
	unaryServerInterceptorBuilder.Use(middleware_logger.LoggingUnaryServerInterceptor())

	//	authHandler := middleware_grpc_auth.GetAuthFuncAccessorFromContainer(serviceProvider.GetContainer())
	//	unaryServerInterceptorBuilder.Use(middleware_grpc_auth.UnaryServerInterceptor(authHandler))

	unaryServerInterceptorBuilder.Use(oauth2.OAuth2UnaryServerInterceptor(oidcContext))
	unaryServerInterceptorBuilder.Use(oauth2.FinalAuthVerificationMiddleware(s.RootContainer))

	unaryServerInterceptorBuilder.Use(middleware_grpc_recovery.UnaryServerInterceptor(recoveryOpts...))

	s.MockOIDCService = mockoidcservice.GetMockOIDCServiceFromContainer(s.RootContainer)
}

// RegisterGRPCEndpoints ...
func (s *Startup) RegisterGRPCEndpoints(server *grpc.Server) []interface{} {
	var endpoints []interface{}
	endpoints = append(endpoints, pb.RegisterGreeterServerDI(server))
	endpoints = append(endpoints, pb.RegisterGreeter2ServerDI(server))
	return endpoints
}
func recoveryUnaryFunc(fullMethodName string, p interface{}) (interface{}, error) {
	fmt.Printf("p: %+v\n", p)

	replyFunc := pb.Get_helloworldFullEmptyResponseFromFullMethodName(fullMethodName)
	if replyFunc != nil {
		reply, ok2 := replyFunc().(grpcDIProtoError.IError)
		if ok2 {
			myError := reply.GetError()
			myError.Code = 503
			myError.Message = "Unexpected error2"
			return reply, status.Error(codes.Internal, "Unexpected error2")
		}
	}

	return nil, status.Error(codes.Internal, "Unexpected error1")
}
