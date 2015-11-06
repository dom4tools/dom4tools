package command

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type CdCommand struct {
	*Meta

	GameName string
}

// Print out the name of a directory, useful for scripting.
// This will either be the Dominions 4 data directory, or the directory of the named game
func (c *CdCommand) run(*kingpin.ParseContext) error {
	if c.GameName != "" {
		game, err := c.Meta.RunContext.GameInstallation.AvailableGames.FindGameByName(c.GameName)
		if err != nil {
			return err
		}

		c.Ui.Output(game.Directory)
	} else {
		c.Ui.Output(c.Meta.RunContext.BasePath)
	}

	return nil
}

func (c *CdCommand) completion(parseContext *kingpin.ParseContext) error {
	return completionWithGames(c.Meta, parseContext)
}

func ConfigureCdCommand(app *kingpin.Application, meta *Meta) (commandName string) {
	commandName = "cd"
	c := &CdCommand{Meta: meta}
	cmd := app.Command(commandName, "Print name of Dominions 4 data directory OR specified game directory (for use in scripting).")

	if meta.CompletionOnly {
		cmd.Action(c.completion)
	} else {
		cmd.Action(c.run)
		cmd.Arg("game_name", "Name of the game").StringVar(&c.GameName)
	}

	return commandName
}
