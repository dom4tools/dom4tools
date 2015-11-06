package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwohFileBackupFilename(t *testing.T) {
	twohfile := TwohFile{Filename: "early_agartha.2h"}
	twohNumber := 4

	backupFilename, _ := twohfile.BackupFilename(twohNumber)

	assert.Equal(t, "early_agartha-4.2h", backupFilename)
}

func TestTwohFileBackupFilenameWithoutFilename(t *testing.T) {
	twohfile := TwohFile{}
	twohNumber := 4

	_, err := twohfile.BackupFilename(twohNumber)

	assert.Error(t, err)
}

func TestTwohFileBackupFilepath(t *testing.T) {
	twohfile := TwohFile{Fullpath: "~/dominions4/savedgames/testgame/early_agartha.2h"}
	twohNumber := 5

	backupFilepath, _ := twohfile.BackupFilepath(twohNumber)

	assert.Equal(t, "~/dominions4/savedgames/testgame/early_agartha-5.2h", backupFilepath)
}

func TestTwohFileBackupFilepathWithoutFilepath(t *testing.T) {
	twohfile := TwohFile{}
	twohNumber := 5

	_, err := twohfile.BackupFilepath(twohNumber)

	assert.Error(t, err)
}
