package game

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/promisedlandt/dom4tools/utility"
)

type ImapConfig struct {
	Port     string
	Server   string
	Username string
	Password string
}

// Get the turn by using a builtin mailer
func (game *Game) GetTurnByMailBuiltin(mailConfig ImapConfig) error {
	return nil
}

// Get the turn by checking the download directory
func (game *Game) GetTurnFromFolder(folder string) error {
	downloadFilepath := filepath.Join(folder, game.TrnFile.Filename)

	if !utility.FileExists(downloadFilepath) {
		return errors.New(fmt.Sprintf("%v does not exist", downloadFilepath))
	}

	err := utility.Mv(downloadFilepath, game.TrnFile.Fullpath)
	if err != nil {
		return err
	}

	return nil
}
