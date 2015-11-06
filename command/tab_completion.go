package command

import "gopkg.in/alecthomas/kingpin.v2"

type TabCompletable interface {
	completion() error
}

// Tab completion with all games we can find
func completionWithGames(meta *Meta, parseContext *kingpin.ParseContext) error {
	listCommand := &ListCommand{Meta: meta}
	return listCommand.run(parseContext)
}

// Empty tab completion
func noCompletion() error {
	return nil
}
