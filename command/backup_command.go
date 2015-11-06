package command

import (
	"fmt"

	"github.com/promisedlandt/dom4tools/game"

	"gopkg.in/alecthomas/kingpin.v2"
)

type BackupCommand struct {
	*Meta

	Game       *game.Game
	GameName   string
	TurnNumber int
	Force      bool
}

func (c *BackupCommand) run(*kingpin.ParseContext) error {
	if c.Game == nil {
		game, err := c.Meta.RunContext.GameInstallation.AvailableGames.FindGameByName(c.GameName)
		if err != nil {
			return err
		}

		c.Game = &game
	}

	c.Ui.Output(fmt.Sprintf("Backing up game %v, turn number %v", c.Game.Name, c.TurnNumber))
	err := c.Game.Backup(c.TurnNumber, c.Force)
	if err != nil {
		return err
	}

	return nil
}

func (c *BackupCommand) completion(parseContext *kingpin.ParseContext) error {
	return completionWithGames(c.Meta, parseContext)
}

func ConfigureBackupCommand(app *kingpin.Application, meta *Meta) (commandName string) {
	commandName = "backup"
	c := &BackupCommand{Meta: meta}
	cmd := app.Command(commandName, "Backs up the current trn and 2h files of game_name as turn turn_number, allowing them to be restored at a later date.")

	if meta.CompletionOnly {
		cmd.Action(c.completion)
	} else {
		cmd.Action(c.run)
		cmd.Arg("game_name", "Name of the game to backup").Required().StringVar(&c.GameName)
		cmd.Arg("turn_number", "Back up which turn number?").Required().IntVar(&c.TurnNumber)
		cmd.Flag("force", "overwrite existing backup").Short('f').BoolVar(&c.Force)
	}

	return commandName
}
