package command

import "gopkg.in/alecthomas/kingpin.v2"

type ListCommand struct {
	*Meta
}

// Lists all games we can find for the current installations
func (c *ListCommand) run(*kingpin.ParseContext) error {
	for _, game := range c.Meta.RunContext.GameInstallation.AvailableGames {
		c.Ui.Output(game.Name)
	}

	return nil
}

func (c *ListCommand) completion(parseContext *kingpin.ParseContext) error {
	return noCompletion()
}

func ConfigureListCommand(app *kingpin.Application, meta *Meta) (commandName string) {
	commandName = "list"
	c := &ListCommand{Meta: meta}
	cmd := app.Command(commandName, "List all games.")

	if meta.CompletionOnly {
		cmd.Action(c.completion)
	} else {
		cmd.Action(c.run)
	}

	return commandName
}
