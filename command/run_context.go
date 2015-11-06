package command

import (
	"github.com/jacobstr/confer"
	"github.com/promisedlandt/dom4tools/game"
)

type RunContext struct {
	BasePath              string
	BaseConfigurationPath string
	DownloadsDirectory    string
	Config                *confer.Config

	GameInstallation game.GameInstallation
}

func (runContext *RunContext) Finalize() error {
	gameInstallation := game.NewGameInstallation(runContext.BasePath)

	runContext.GameInstallation = *gameInstallation

	config := confer.NewConfig()
	config.ReadPaths(runContext.BaseConfigurationPath)
	runContext.Config = config

	return nil
}
