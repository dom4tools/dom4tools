package main

import (
	"log"
	"os"
	"runtime"

	"github.com/promisedlandt/dom4tools/command"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/alecthomas/kingpin.v2"
)

var CurrentRunContext command.RunContext

var Ui cli.Ui

const (
	ErrorPrefix  = "!!!: "
	AskPrefix    = "???: "
	OutputPrefix = "d4t: "
)

func main() {
	var useBasicUi bool
	var completionMode bool // for use in autocompletion scripts, don't execute commands
	colorizeUi := true
	args := os.Args

	if len(args) > 1 {
		if args[1] == "cd" {
			useBasicUi = true
			colorizeUi = false
		}
	}

	for i, arg := range args {
		if arg == "--bash-completion" {
			useBasicUi = true
			colorizeUi = false
			completionMode = true

			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	for i, arg := range args {
		if arg == "--simple-output" {
			useBasicUi = true
			colorizeUi = false

			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	if useBasicUi {
		Ui = &cli.BasicUi{Writer: os.Stdout}
	} else {
		Ui = &cli.PrefixedUi{
			AskPrefix:    AskPrefix,
			OutputPrefix: OutputPrefix,
			InfoPrefix:   OutputPrefix,
			ErrorPrefix:  ErrorPrefix,
			Ui:           &cli.BasicUi{Writer: os.Stdout},
		}
	}

	// TODO: check if exists
	CurrentRunContext.BasePath = defaultDominions4BasePath()
	CurrentRunContext.BaseConfigurationPath = defaultDominions4BaseConfigurationPath()
	CurrentRunContext.DownloadsDirectory = defaultDownloadsDirectory()

	meta := command.Meta{
		Ui:             Ui,
		RunContext:     &CurrentRunContext,
		Color:          colorizeUi,
		CompletionOnly: completionMode,
	}

	args, err := meta.Process(args[1:])

	if err != nil {
		meta.Ui.Error(err.Error())
		os.Exit(1)
	}

	exitStatus := 0
	var commandNames []string

	app := kingpin.New("d4t", "Manage your Dominions 4 games from the command line.")
	commandNames = append(commandNames, command.ConfigureCdCommand(app, &meta))
	commandNames = append(commandNames, command.ConfigureListCommand(app, &meta))
	commandNames = append(commandNames, command.ConfigureCreateCommand(app, &meta))
	commandNames = append(commandNames, command.ConfigureBackupCommand(app, &meta))
	commandNames = append(commandNames, command.ConfigureRestoreCommand(app, &meta))
	commandNames = append(commandNames, command.ConfigureReplayCommand(app, &meta))
	commandNames = append(commandNames, command.ConfigureSubmitCommand(app, &meta))
	commandNames = append(commandNames, command.ConfigureResubmitCommand(app, &meta))
	commandNames = append(commandNames, command.ConfigureGetCommand(app, &meta))
	commandNames = append(commandNames, command.ConfigureVersionCommand(app, &meta, Version, VersionPrerelease, GitCommit))

	// Show the names of the subcommands but execute no commands
	if completionMode && len(args) == 0 {
		for _, commandName := range commandNames {
			meta.Ui.Output(commandName)
		}
		os.Exit(0)
	}

	_, err = app.Parse(args)

	if err != nil {
		if len(err.Error()) > 0 {
			meta.Ui.Error(err.Error())
		}
		exitStatus = 1
	}

	os.Exit(exitStatus)
}

// Path to the directory where downloaded files are stored by default
func defaultDownloadsDirectory() (downloadsDirectory string) {
	if runtime.GOOS == "windows" {
		downloadsDirectory = os.ExpandEnv("${USERPROFILE}\\Downloads")
	} else {
		base, err := homedir.Expand("~/Downloads")
		if err != nil {
			log.Fatal(err)
		}
		downloadsDirectory = base
	}

	return
}

// Path to the default dom4tools configuration file
func defaultDominions4BaseConfigurationPath() (configurationPath string) {
	if runtime.GOOS == "windows" {
		configurationPath = os.ExpandEnv("${LOCALAPPDATA}\\dom4tools\\config.json")
	} else {
		base, err := homedir.Expand("~/.dom4tools/config.json")
		if err != nil {
			log.Fatal(err)
		}
		configurationPath = base
	}

	return
}

// Path to the default Dominions 4 data directory
func defaultDominions4BasePath() (basePath string) {
	if runtime.GOOS == "windows" {
		basePath = os.ExpandEnv("${APPDATA}\\Dominions4")
	} else {
		base, err := homedir.Expand("~/dominions4")
		if err != nil {
			log.Fatal(err)
		}
		basePath = base
	}

	return
}
