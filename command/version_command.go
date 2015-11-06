package command

import (
	"bytes"
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"
)

type VersionCommand struct {
	*Meta

	Revision          string
	Version           string
	VersionPrerelease string
}

func (c *VersionCommand) run(*kingpin.ParseContext) error {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "dom4tools v%s", c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, "-%s", c.VersionPrerelease)

		if c.Revision != "" {
			fmt.Fprintf(&versionString, " (%s)", c.Revision)
		}
	}

	c.Ui.Output(versionString.String())

	return nil
}

func (c *VersionCommand) completion(parseContext *kingpin.ParseContext) error {
	return noCompletion()
}

func ConfigureVersionCommand(app *kingpin.Application, meta *Meta, version string, versionPrerelease string, revision string) (commandName string) {
	commandName = "version"
	c := &VersionCommand{Meta: meta, Version: version, VersionPrerelease: versionPrerelease, Revision: revision}
	cmd := app.Command(commandName, "Version")

	if meta.CompletionOnly {
		cmd.Action(c.completion)
	} else {
		cmd.Action(c.run)
	}

	return commandName
}
