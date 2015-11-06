package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrnFileBackupFilename(t *testing.T) {
	trnfile := TrnFile{Filename: "early_agartha.trn"}
	trnNumber := 4

	backupFilename, _ := trnfile.BackupFilename(trnNumber)

	assert.Equal(t, "early_agartha-4.trn", backupFilename)
}

func TestTrnFileBackupFilenameWithoutFilename(t *testing.T) {
	trnfile := TrnFile{}
	trnNumber := 4

	_, err := trnfile.BackupFilename(trnNumber)

	assert.Error(t, err)
}

func TestTrnFileBackupFilepath(t *testing.T) {
	trnfile := TrnFile{Fullpath: "~/dominions4/savedgames/testgame/early_agartha.trn"}
	trnNumber := 5

	backupFilepath, _ := trnfile.BackupFilepath(trnNumber)

	assert.Equal(t, "~/dominions4/savedgames/testgame/early_agartha-5.trn", backupFilepath)
}

func TestTrnFileBackupFilepathWithoutFilepath(t *testing.T) {
	trnfile := TrnFile{}
	trnNumber := 5

	_, err := trnfile.BackupFilepath(trnNumber)

	assert.Error(t, err)
}
