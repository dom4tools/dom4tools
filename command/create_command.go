package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/promisedlandt/dom4tools/game"

	"gopkg.in/alecthomas/kingpin.v2"
)

type CreateCommand struct {
	*Meta

	Force       bool
	NewGameName string
}

// Create a new game for the installation.
// That means creating a new directory (if allowed)
func (c *CreateCommand) run(*kingpin.ParseContext) error {
	// Is the new game name valid?
	if validationMessages, valid := game.ValidGameName(c.NewGameName); !valid {
		for _, validationMessage := range validationMessages {
			c.Ui.Error(validationMessage)
		}
		return errors.New("")
	}

	// Does another game with the same name exist? (case insensitive)
	gameWithSameNameIndex := -1
	for index, game := range c.Meta.RunContext.GameInstallation.AvailableGames {
		if strings.ToLower(game.Name) == strings.ToLower(c.NewGameName) {
			gameWithSameNameIndex = index
			break
		}
	}

	if gameWithSameNameIndex > -1 {
		existingGame := c.Meta.RunContext.GameInstallation.AvailableGames[gameWithSameNameIndex]

		if c.Force {
			c.Ui.Info(fmt.Sprintf("Overwriting existing game %s at %s", existingGame.Name, existingGame.Directory))

			if err := existingGame.Delete(); err != nil {
				return err
			}
		} else {
			c.Ui.Error(fmt.Sprintf("Game already exists: %s at %s", existingGame.Name, existingGame.Directory))
			c.Ui.Error("if you want to overwrite, call with -f or --force")
			return errors.New("")
		}
	} else {
		c.Ui.Output(fmt.Sprintf("Creating %v", c.NewGameName))
	}

	newGame := game.Game{Name: c.NewGameName}

	if err := newGame.Create(c.Meta.RunContext.GameInstallation.SavedGamesPath); err != nil {
		return err
	}

	c.Meta.RunContext.GameInstallation.Update()

	return nil
}

func (c *CreateCommand) completion(parseContext *kingpin.ParseContext) error {
	return noCompletion()
}

func ConfigureCreateCommand(app *kingpin.Application, meta *Meta) (commandName string) {
	commandName = "create"
	c := &CreateCommand{Meta: meta}
	cmd := app.Command(commandName, "Create a new game.")

	if meta.CompletionOnly {
		cmd.Action(c.completion)
	} else {
		cmd.Arg("game_name", "Name of the game").Required().StringVar(&c.NewGameName)
		cmd.Flag("force", "overwrite existing game").Short('f').BoolVar(&c.Force)
		cmd.Action(c.run)
	}

	return commandName
}
