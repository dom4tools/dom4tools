package command

import (
	"errors"
	"fmt"

	"github.com/promisedlandt/dom4tools/game"

	"gopkg.in/alecthomas/kingpin.v2"
)

type SubmitCommand struct {
	*Meta

	Resubmit   bool
	SkipBackup bool
	GameName   string
	TurnNumber int
	Game       *game.Game
	SmtpConfig SmtpConfig
}

// Submits the given game
func (c *SubmitCommand) run(parseContext *kingpin.ParseContext) error {
	if c.Game == nil {
		game, err := c.Meta.RunContext.GameInstallation.AvailableGames.FindGameByName(c.GameName)
		if err != nil {
			return err
		}

		c.Game = &game
	}

	if c.TurnNumber <= 0 {
		c.TurnNumber = c.Game.CurrentTurnNumber()

		if c.Resubmit {
			c.TurnNumber--
		}
	}

	if c.TurnNumber <= 0 {
		return errors.New(fmt.Sprintf("No turn set to submit, try: d4t submit %v TURN_NUMBER", c.Game.Name))
	}

	if !c.SkipBackup {
		backupCommand := BackupCommand{Meta: c.Meta, Game: c.Game, TurnNumber: c.TurnNumber, Force: c.Resubmit}
		err := backupCommand.run(parseContext)
		if err != nil {
			return err
		}
	}

	switch c.Meta.Config.Submitstyle {
	case "smtp":
		if len(c.Meta.Config.Smtpsettings.From) == 0 {
			return errors.New("no \"from\" set in smtpsettings")
		}

		if len(c.Meta.Config.Smtpsettings.Port) == 0 {
			return errors.New("no port set in smtpsettings")
		}

		if len(c.Meta.Config.Smtpsettings.Server) == 0 {
			return errors.New("no server set in smtpsettings")
		}

		if len(c.Meta.Config.Smtpsettings.Username) == 0 {
			return errors.New("no username set in smtpsettings")
		}

		if len(c.Meta.Config.Smtpsettings.Password) == 0 {
			return errors.New("no password set in smtpsettings")
		}

		c.SmtpConfig = SmtpConfig{To: "turns@llamaserver.net", From: c.Meta.Config.Smtpsettings.From, Port: c.Meta.Config.Smtpsettings.Port, Server: c.Meta.Config.Smtpsettings.Server, Username: c.Meta.Config.Smtpsettings.Username, Password: c.Meta.Config.Smtpsettings.Password, Subject: fmt.Sprintf("%v turn %v", c.Game.Name, c.TurnNumber), Body: "", AttachmentPath: c.Game.TwohFile.Fullpath}

		c.Ui.Output(fmt.Sprintf("Submitting game %s, turn %v", c.Game.Name, c.TurnNumber))

		err := c.SmtpConfig.SubmitTurnBuiltin()
		if err != nil {
			return err
		}
	default:
		return errors.New("No submitstyle set in config")
	}

	return nil
}

func (c *SubmitCommand) completion(parseContext *kingpin.ParseContext) error {
	return completionWithGames(c.Meta, parseContext)
}

func ConfigureSubmitCommand(app *kingpin.Application, meta *Meta) (commandName string) {
	commandName = "submit"
	c := &SubmitCommand{Meta: meta}
	cmd := app.Command(commandName, "Submits a game.")

	if meta.CompletionOnly {
		cmd.Action(c.completion)
	} else {
		cmd.Arg("game_name", "Name of the game to submit").Required().StringVar(&c.GameName)
		cmd.Arg("turn_number", "Submit which turn? Needed for backup").IntVar(&c.TurnNumber)
		cmd.Flag("skip-backup", "don't back up").Short('b').BoolVar(&c.SkipBackup)
		cmd.Action(c.run)
	}

	return commandName
}
