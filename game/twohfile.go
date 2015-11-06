package game

import (
	"errors"
	"path/filepath"
)

type TwohFile struct {
	Filename string
	Fullpath string
}

// Returns the backup file name for a 2h file for the given turn number for this game
func (twohfile *TwohFile) BackupFilename(turnNumber int) (string, error) {
	if twohfile.Filename == "" {
		return "", errors.New("No 2h file to back up found")
	}

	return Backup2hFilename(twohfile.Filename, turnNumber), nil
}

// Returns the backup file path for a 2h file for the given turn number for this game
func (twohfile *TwohFile) BackupFilepath(turnNumber int) (string, error) {
	if twohfile.Fullpath == "" {
		return "", errors.New("No turn file to back up found")
	}

	dir, file := filepath.Split(twohfile.Fullpath)

	backupFilename := Backup2hFilename(file, turnNumber)

	return filepath.Join(dir, backupFilename), nil
}

func (twohfile *TwohFile) BackupBasename() (string, error) {
	if twohfile.Fullpath == "" {
		return "", errors.New("No 2h file to get basename for")
	}

	return Backup2hBasename(twohfile.Filename), nil
}
