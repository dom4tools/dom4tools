package game

import (
	"regexp"
	"strconv"
	"strings"
)

// Backup2hFilename returns the backup file name for a 2h file for the given turn number.
// Example: early_agartha.2h -> early_agartha-12.2h
func Backup2hFilename(twohFile string, turnNumber int) string {
	return BackupFilename(twohFile, turnNumber, "2h")
}

// BackupTrnFilename returns the backup file name for a trn file for the given turn number.
// Example: early_agartha.trn -> early_agartha-12.trn
func BackupTrnFilename(trnFile string, turnNumber int) string {
	return BackupFilename(trnFile, turnNumber, "trn")
}

// BackupTrnBasename returns the name of a trn file without any backup information.
// Example: early_agartha-12.trn -> early_agartha.trn
func Backup2hBasename(twohFile string) string {
	return BackupBasename(twohFile, "2h")
}

// BackupTrnBasename returns the name of a trn file without any backup information.
// Example: early_agartha-12.trn -> early_agartha.trn
func BackupTrnBasename(trnFile string) string {
	return BackupBasename(trnFile, "trn")
}

// BackupFilename returns the backup file name for a file with the given extension for the given turn number.
func BackupFilename(inputFilename string, turnNumber int, fileExtension string) string {
	fileExtensionStart := strings.LastIndex(inputFilename, "."+fileExtension)
	newName := inputFilename[0:fileExtensionStart] + "-" + strconv.Itoa(turnNumber) + inputFilename[fileExtensionStart:]

	return newName
}

// BackupBasename returns the name of a file without any backup information.
func BackupBasename(inputFilename string, fileExtension string) string {
	filenameRegexp := regexp.MustCompile(`(.*)-\d+(\.` + fileExtension + `)`)

	if matchData := filenameRegexp.FindStringSubmatch(inputFilename); matchData != nil {
		return matchData[1] + matchData[2]
	} else {
		return inputFilename
	}
}
