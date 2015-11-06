package game

import (
	"errors"
	"path/filepath"
)

type TrnFile struct {
	Filename string
	Fullpath string
}

// Returns the backup file name for a trn file for the given turn number for this game.
func (trnfile *TrnFile) BackupFilename(turnNumber int) (string, error) {
	if trnfile.Filename == "" {
		return "", errors.New("No turn file to back up found")
	}

	return BackupTrnFilename(trnfile.Filename, turnNumber), nil
}

// Returns the backup file path for a Trn file for the given turn number for this game.
func (trnfile *TrnFile) BackupFilepath(turnNumber int) (string, error) {
	if trnfile.Fullpath == "" {
		return "", errors.New("No turn file to back up found")
	}

	dir, file := filepath.Split(trnfile.Fullpath)

	backupFilename := BackupTrnFilename(file, turnNumber)

	return filepath.Join(dir, backupFilename), nil
}

func (trnfile *TrnFile) BackupBasename() (string, error) {
	if trnfile.Fullpath == "" {
		return "", errors.New("No turn file to get basename for")
	}

	return BackupTrnBasename(trnfile.Filename), nil
}
