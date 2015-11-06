package command

import (
	"errors"
	"fmt"

	"github.com/promisedlandt/dom4tools/game"

	"gopkg.in/alecthomas/kingpin.v2"
)

type GetCommand struct {
	*Meta

	GameName string
	Game     *game.Game
	//ImapConfig ImapConfig
}

// Gets the given game
func (c *GetCommand) run(*kingpin.ParseContext) error {
	game, err := c.Meta.RunContext.GameInstallation.AvailableGames.FindGameByName(c.GameName)
	if err != nil {
		return err
	}

	c.Game = &game

	switch c.Meta.Config.Getstyle {
	case "folder":
		if len(c.Meta.RunContext.DownloadsDirectory) == 0 {
			return errors.New("No download folder set in config")
		}

		err = game.GetTurnFromFolder(c.Meta.RunContext.DownloadsDirectory)
		if err != nil {
			return err
		}

		c.Ui.Output(fmt.Sprintf("Got turn for %v", c.Game.Name))

	//case "imap":
	//if len(c.Meta.Config.Imapsettings.Port) == 0 {
	//return errors.New("no port set in imapsettings")
	//}

	//if len(c.Meta.Config.Imapsettings.Server) == 0 {
	//return errors.New("no server set in imapsettings")
	//}

	//if len(c.Meta.Config.Imapsettings.Username) == 0 {
	//return errors.New("no username set in imapsettings")
	//}

	//if len(c.Meta.Config.Imapsettings.Password) == 0 {
	//return errors.New("no password set in imapsettings")
	//}

	//c.ImapConfig = ImapConfig{Port: c.Meta.Config.Imapsettings.Port, Server: c.Meta.Config.Imapsettings.Server, Username: c.Meta.Config.Imapsettings.Username, Password: c.Meta.Config.Imapsettings.Password}

	//log.Printf("Getting turn for %s", game.Name)

	//err = game.GetTurnByMailBuiltin(c.ImapConfig)
	//if err != nil {
	//return err
	//}
	default:
		return errors.New("No getstyle set in config")
	}

	return nil
}

func (c *GetCommand) completion(parseContext *kingpin.ParseContext) error {
	return completionWithGames(c.Meta, parseContext)
}

func ConfigureGetCommand(app *kingpin.Application, meta *Meta) (commandName string) {
	commandName = "get"
	c := &GetCommand{Meta: meta}
	cmd := app.Command(commandName, "Gets turn for a game.")

	if meta.CompletionOnly {
		cmd.Action(c.completion)
	} else {
		cmd.Arg("game_name", "Name of the game you want to get the turn for").Required().StringVar(&c.GameName)
		cmd.Action(c.run)
	}

	return commandName
}
