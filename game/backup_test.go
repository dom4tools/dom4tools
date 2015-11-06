package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackup2hFilename(t *testing.T) {
	filename := "early_agartha.2h"
	turnNumber := 17

	backup2hFilename := Backup2hFilename(filename, turnNumber)

	assert.Equal(t, "early_agartha-17.2h", backup2hFilename)
}

func TestBackupTrnFilename(t *testing.T) {
	filename := "early_agartha.trn"
	turnNumber := 5

	backupTrnFilename := BackupTrnFilename(filename, turnNumber)

	assert.Equal(t, "early_agartha-5.trn", backupTrnFilename)
}

func TestBackup2hBasename(t *testing.T) {
	filename := "early_agartha-12.2h"

	backupFilename := Backup2hBasename(filename)
	assert.Equal(t, "early_agartha.2h", backupFilename)
}

func TestBackupTrnBasename(t *testing.T) {
	filename := "early_agartha-12.trn"

	backupFilename := BackupTrnBasename(filename)
	assert.Equal(t, "early_agartha.trn", backupFilename)
}
