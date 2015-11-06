package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameValidGameName(t *testing.T) {
	_, valid1 := ValidGameName("testgame")
	assert.True(t, valid1)

	_, valid2 := ValidGameName("test_game")
	assert.True(t, valid2)

	_, valid3 := ValidGameName("test game")
	assert.False(t, valid3)

	_, valid4 := ValidGameName("newlords")
	assert.False(t, valid4)
}

func TestReplayName(t *testing.T) {
	game, _ := NewGame("testgame", "/home/test/dominions4/savegaes/testgame")

	replayName := game.ReplayName(24)

	assert.Equal(t, "testgame24", replayName)
}
