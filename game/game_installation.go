package game

import (
	"io/ioutil"
	"log"
	"path"
)

// GameInstallation represents a Dominions 4 installation
type GameInstallation struct {
	BasePath       string
	SavedGamesPath string
	AvailableGames GameCollection
}

func NewGameInstallation(basePath string) *GameInstallation {
	gameInstallation := GameInstallation{BasePath: basePath}
	gameInstallation.Update()

	return &gameInstallation
}

func (gameInstallation *GameInstallation) Update() {
	gameInstallation.SavedGamesPath = savedGamesPath(*gameInstallation)
	gameInstallation.AvailableGames = availableGames(*gameInstallation)
}

func savedGamesPath(gameInstallation GameInstallation) string {
	return path.Join(gameInstallation.BasePath, "savedgames")
}

// Read all available games for the given installation from the installation directory
func availableGames(gameInstallation GameInstallation) []Game {
	var games []Game

	files, err := ioutil.ReadDir(gameInstallation.SavedGamesPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if _, valid := ValidGameName(f.Name()); valid {
			game, err := NewGame(f.Name(), path.Join(gameInstallation.SavedGamesPath, f.Name()))
			if err == nil {
				games = append(games, *game)
			}
		}
	}

	return games
}
