package command

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureResubmitCommand(app *kingpin.Application, meta *Meta) (commandName string) {
	commandName = "resubmit"
	c := &SubmitCommand{Meta: meta, Resubmit: true}
	cmd := app.Command(commandName, "Resubmits a game.")

	if meta.CompletionOnly {
		cmd.Action(c.completion)
	} else {
		cmd.Arg("game_name", "Name of the game to resubmit").Required().StringVar(&c.GameName)
		cmd.Arg("turn_number", "Resubmit which turn? Needed for backup").IntVar(&c.TurnNumber)
		cmd.Action(c.run)
	}

	return commandName
}
