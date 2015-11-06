package command

import (
	"fmt"
	"path/filepath"

	"github.com/promisedlandt/dom4tools/utility"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ReplayCommand struct {
	*Meta

	Force     bool
	Destroy   bool
	GameName  string
	StartTurn int
	TurnCount int
}

func (c *ReplayCommand) run(parseContext *kingpin.ParseContext) error {
	game, err := c.Meta.RunContext.GameInstallation.AvailableGames.FindGameByName(c.GameName)
	if err != nil {
		return err
	}

	var endTurn int
	lastSavedTurn := game.SortedTrnBackupKeys[len(game.SortedTrnBackupKeys)-1]

	if c.TurnCount > 0 {
		endTurn = c.StartTurn + c.TurnCount - 1

		if endTurn > lastSavedTurn {
			endTurn = lastSavedTurn
		}
	} else {
		endTurn = lastSavedTurn
	}

	if c.Destroy {
		c.Ui.Output(fmt.Sprintf("Replaying turns for %v, starting at %v, ending at %v", game.Name, c.StartTurn, endTurn))

		for turn := c.StartTurn; turn <= endTurn; turn++ {

			replayGameName := game.ReplayName(turn)
			replayGame, err := c.Meta.RunContext.GameInstallation.AvailableGames.FindGameByName(replayGameName)

			if err != nil {
				c.Ui.Output(fmt.Sprintf("No game found for turn %v", turn))
			} else {
				c.Ui.Output(fmt.Sprintf("Deleting %v", replayGame.Name))
				replayGame.Delete()
			}
		}
	} else {
		c.Ui.Output(fmt.Sprintf("Deleting replays for %v, starting at %v, ending at %v", game.Name, c.StartTurn, endTurn))

		for turn := c.StartTurn; turn <= endTurn; turn++ {
			trnFile, ok := game.TrnBackups[turn]
			if !ok {
				c.Ui.Error(fmt.Sprintf("No .trn file found for turn %v, skipping", turn))
				continue
			}

			twohFile, ok := game.TwohBackups[turn]
			if !ok {
				c.Ui.Error(fmt.Sprintf("No .2h file found for turn %v, skipping", turn))
				continue
			}

			newGameName := game.ReplayName(turn)
			newGameCmd := CreateCommand{Meta: c.Meta, Force: c.Force, NewGameName: newGameName}
			newGameCmd.run(parseContext)
			newGame, err := c.Meta.RunContext.GameInstallation.AvailableGames.FindGameByName(newGameName)
			if err != nil {
				return err
			}

			backupTrnBasename, err := trnFile.BackupBasename()
			if err != nil {
				return err
			}

			backupTwohBasename, err := twohFile.BackupBasename()
			if err != nil {
				return err
			}

			err = utility.Cp(trnFile.Fullpath, filepath.Join(newGame.Directory, backupTrnBasename))
			if err != nil {
				return err
			}

			err = utility.Cp(twohFile.Fullpath, filepath.Join(newGame.Directory, backupTwohBasename))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ReplayCommand) completion(parseContext *kingpin.ParseContext) error {
	return completionWithGames(c.Meta, parseContext)
}

func ConfigureReplayCommand(app *kingpin.Application, meta *Meta) (commandName string) {
	commandName = "replay"
	c := &ReplayCommand{Meta: meta}
	cmd := app.Command(commandName, "Create new games for every backed up turn for the given game.")

	if meta.CompletionOnly {
		cmd.Action(c.completion)
	} else {
		cmd.Action(c.run)
		cmd.Arg("game_name", "Name of the game to replay").Required().StringVar(&c.GameName)
		cmd.Flag("force", "overwrite existing games").Short('f').BoolVar(&c.Force)
		cmd.Flag("delete", "delete games instead of creating them").Short('d').BoolVar(&c.Destroy)
		cmd.Flag("start-turn", "Start on which turn?").Short('s').Default("1").IntVar(&c.StartTurn)
		cmd.Flag("count", "Replay how many turns?").Short('c').IntVar(&c.TurnCount)
	}

	return commandName
}
