package plugin

import (
	"github.com/fluffy-bunny/grpcdotnetgo-samples/simple/internal/startup"
	coreContracts "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/core"
	pluginContracts "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/plugin"
	grpcdotnetgo_plugin "github.com/fluffy-bunny/grpcdotnetgo/pkg/plugin"
)

func init() {
	grpcdotnetgo_plugin.AddPlugin(NewPlugin())
}

type pluginService struct {
}

// NewPlugin ...
func NewPlugin() pluginContracts.IGRPCDotNetGoPlugin {
	return &pluginService{}
}

// startup gets name of plugin
func (p *pluginService) GetName() string {
	return "example"
}

// GetStartup gets startup object
func (p *pluginService) GetStartup() coreContracts.IStartup {
	return startup.NewStartup()
}
