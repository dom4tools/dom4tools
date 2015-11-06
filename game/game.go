package game

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/promisedlandt/dom4tools/utility"
)

type GameCollection []Game

// These game names are used by Dominions 4, and we want nothing to do with them
var ReservedGameNames = []string{"newlords"}

// Game represents a Dominions 4 game
type Game struct {
	Name      string
	Directory string

	TwohFile             TwohFile
	TrnFile              TrnFile
	TwohBackups          map[int]TwohFile
	TrnBackups           map[int]TrnFile
	SortedTwohBackupKeys []int
	SortedTrnBackupKeys  []int
}

func NewGame(name string, basedir string) (*Game, error) {
	game := Game{Name: name, Directory: basedir, TwohBackups: make(map[int]TwohFile), TrnBackups: make(map[int]TrnFile)}

	current2hFile, err2hFile := game.Current2hFile()
	current2hFilepath, err2hPath := game.Current2hFilepath()

	if err2hFile == nil && err2hPath == nil {
		game.TwohFile = TwohFile{Filename: current2hFile, Fullpath: current2hFilepath}
	}

	currentTrnFile, errTrnFile := game.CurrentTrnFile()
	currentTrnFilepath, errTrnPath := game.CurrentTrnFilepath()

	if errTrnFile == nil && errTrnPath == nil {
		game.TrnFile = TrnFile{Filename: currentTrnFile, Fullpath: currentTrnFilepath}
	}

	files, err := ioutil.ReadDir(game.Directory)
	if err != nil {
		return &game, err
	}

	twohRegexp := regexp.MustCompile(`-(\d+)\.2h\z`)
	trnRegexp := regexp.MustCompile(`-(\d+)\.trn\z`)

	for _, f := range files {
		if matchData := twohRegexp.FindStringSubmatch(f.Name()); matchData != nil {
			turnNumber, err := strconv.Atoi(matchData[1])

			if err == nil {
				game.TwohBackups[turnNumber] = TwohFile{Filename: f.Name(), Fullpath: path.Join(basedir, f.Name())}
			}
		} else if matchData := trnRegexp.FindStringSubmatch(f.Name()); matchData != nil {
			turnNumber, err := strconv.Atoi(matchData[1])

			if err == nil {
				game.TrnBackups[turnNumber] = TrnFile{Filename: f.Name(), Fullpath: path.Join(basedir, f.Name())}
			}
		}
	}

	for key := range game.TwohBackups {
		game.SortedTwohBackupKeys = append(game.SortedTwohBackupKeys, key)
	}
	sort.Ints(game.SortedTwohBackupKeys)

	for key := range game.TrnBackups {
		game.SortedTrnBackupKeys = append(game.SortedTrnBackupKeys, key)
	}
	sort.Ints(game.SortedTrnBackupKeys)

	return &game, nil
}

// Return the number of the current game turn, that is the turn number of the highest 2h backup + 1
func (game *Game) CurrentTurnNumber() (currentTurnNumber int) {
	if len(game.SortedTwohBackupKeys) > 0 {
		currentTurnNumber = game.SortedTwohBackupKeys[len(game.SortedTwohBackupKeys)-1] + 1
	} else {
		currentTurnNumber = 1
	}

	return currentTurnNumber
}

// Backup the current trn and 2h files for this game
func (game *Game) Backup(turnNumber int, force bool) error {
	current2hPath := game.TwohFile.Fullpath

	target2hPath, err := game.TwohFile.BackupFilepath(turnNumber)
	if err != nil {
		return err
	}

	currentTrnPath := game.TrnFile.Fullpath

	targetTrnPath, err := game.TrnFile.BackupFilepath(turnNumber)
	if err != nil {
		return err
	}

	if !force && (utility.FileExists(target2hPath) || utility.FileExists(targetTrnPath)) {
		return errors.New(fmt.Sprintf("Backup for turn %v already exists in %v, not forcing", turnNumber, game.Directory))
	}

	err = utility.Cp(current2hPath, target2hPath)
	if err != nil {
		return err
	}

	err = utility.Cp(currentTrnPath, targetTrnPath)
	if err != nil {
		return err
	}

	return nil
}

// Restore backed up trn and 2h file for this game
func (game *Game) Restore(turnNumber int) error {
	current2hPath, err := game.Current2hFilepath()
	if err != nil {
		return err
	}

	currentTrnPath, err := game.CurrentTrnFilepath()
	if err != nil {
		return err
	}

	// It's fine if only the 2h or the trn file have backups we'll just restore the one that exists.
	// But if neither file exists, we error out.
	backup2hPath, err := game.TwohFile.BackupFilepath(turnNumber)
	backup2hExists := err == nil

	backupTrnPath, err := game.TrnFile.BackupFilepath(turnNumber)
	backupTrnExists := err == nil

	if !(backup2hExists || backupTrnExists) {
		return errors.New(fmt.Sprintf("Neither trn nor 2h backups exist for turn %v in %v", turnNumber, game.Directory))
	}

	if backup2hExists {
		err = utility.Cp(backup2hPath, current2hPath)
		if err != nil {
			return err
		}
	}

	if backupTrnExists {
		err = utility.Cp(backupTrnPath, currentTrnPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create the directory for a game
func (game *Game) Create(savedGamesPath string) error {
	newGameDirectory := path.Join(savedGamesPath, game.Name)

	return os.MkdirAll(newGameDirectory, 0755)
}

// Delete the directory for a game
func (game *Game) Delete() error {
	return os.RemoveAll(game.Directory)
}

// Returns the full path for this games current 2h file
func (game *Game) Current2hFilepath() (string, error) {
	fileName, err := game.Current2hFile()
	if err != nil {
		return "", err
	}

	return path.Join(game.Directory, fileName), nil
}

// Returns the full path for this games current trn file
func (game *Game) CurrentTrnFilepath() (string, error) {
	fileName, err := game.CurrentTrnFile()
	if err != nil {
		return "", err
	}

	return path.Join(game.Directory, fileName), nil
}

// Find the current 2h file for a game
func (game *Game) Current2hFile() (string, error) {
	var possibleMatches []string

	files, err := ioutil.ReadDir(game.Directory)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if Valid2hFileName(f.Name()) {
			possibleMatches = append(possibleMatches, f.Name())
		}
	}

	switch {
	case len(possibleMatches) == 0:
		return "", errors.New(fmt.Sprintf("Could not find a 2h file for %s", game.Name))
	case len(possibleMatches) > 1:
		// Take the shortest filename. Ugly, but we have no way to check for valid nation names yet
		sort.Sort(utility.ByLength(possibleMatches))
	}

	return possibleMatches[0], nil
}

// Find the current trn file for a game
func (game *Game) CurrentTrnFile() (string, error) {
	var possibleMatches []string

	files, err := ioutil.ReadDir(game.Directory)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if ValidTrnFileName(f.Name()) {
			possibleMatches = append(possibleMatches, f.Name())
		}
	}

	switch {
	case len(possibleMatches) == 0:
		return "", errors.New(fmt.Sprintf("Could not find a trn file for %s", game.Name))
	case len(possibleMatches) > 1:
		// Take the shortest filename. Ugly, but we have no way to check for valid nation names yet
		sort.Sort(utility.ByLength(possibleMatches))
	}

	return possibleMatches[0], nil
}

// The name of the replay for the given turn number for this game.
// Example: PretendersOfReddit13
func (game *Game) ReplayName(turnNumber int) string {
	return game.Name + strconv.Itoa(turnNumber)
}

// Is the given filename a valid 2h file name?
func Valid2hFileName(filename string) bool {
	return strings.HasSuffix(filename, ".2h")
}

// Is the given filename a valid trn file name?
func ValidTrnFileName(filename string) bool { return strings.HasSuffix(filename, ".trn") }

// Find the game with the given name from the game collection and return it.
// Case insensitive.
func (games GameCollection) FindGameByName(name string) (Game, error) {
	for _, game := range games {
		if strings.ToLower(game.Name) == strings.ToLower(name) {
			return game, nil
		}
	}

	return Game{}, errors.New(fmt.Sprintf("Could not find a game called %s", name))
}

// Checks whether a given game name is valid.
// Game names are invalid when thet are reserved names (used by Dominions 4 itself),
// or when they contain spaces.
func ValidGameName(name string) (messages []string, valid bool) {
	valid = true

	reservedGameNameIndex := -1

	for index, reservedGameName := range ReservedGameNames {
		if strings.ToLower(name) == strings.ToLower(reservedGameName) {
			reservedGameNameIndex = index
			break
		}
	}

	if reservedGameNameIndex > -1 {
		valid = false
		messages = append(messages, fmt.Sprintf("\"%s\" is a reserved game name", ReservedGameNames[reservedGameNameIndex]))
	}

	if strings.Contains(name, " ") {
		valid = false
		messages = append(messages, fmt.Sprint("Game names must not contain spaces"))
	}

	return
}
