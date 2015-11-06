package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestsavedGamesPath(t *testing.T) {
	gameInstallation := GameInstallation{BasePath: "/home/nl/dominions4/"}

	assert.Equal(t, "/home/nl/dominions4/savedgames", savedGamesPath(gameInstallation))
}
