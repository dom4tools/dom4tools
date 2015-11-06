package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
	"github.com/promisedlandt/dom4tools/utility"
)

type Meta struct {
	Ui             cli.Ui
	RunContext     *RunContext
	Color          bool
	CompletionOnly bool
	Config         ConfigStruct

	oldUi cli.Ui
	color bool
}

func (m *Meta) Process(args []string) ([]string, error) {
	context := m.RunContext

	// Colorization only works on posix shells
	if runtime.GOOS == "windows" {
		m.Color = false
	}

	m.color = m.Color

	// Set the UI
	m.oldUi = m.Ui
	m.Ui = &cli.ConcurrentUi{
		Ui: &ColorizeUi{
			Colorize:    m.Colorize(),
			OutputColor: "[green]",
			InfoColor:   "[green]",
			ErrorColor:  "[red]",
			WarnColor:   "[yellow]",
			Ui:          m.oldUi,
		},
	}

	configPath := m.RunContext.BaseConfigurationPath

	// Create default configuration file if it doesn't already exist
	if !utility.FileExists(configPath) {
		m.Ui.Info(fmt.Sprintf("No configuration found at %v, creating default config", m.RunContext.BaseConfigurationPath))
		err := m.CreateDefaultConfig()
		if err != nil {
			return args, err
		}
	}

	// Since we store sensitive information in the config file, including email passwords, make sure it's only readable by owner
	configFileInfo, err := os.Stat(m.RunContext.BaseConfigurationPath)
	if err != nil {
		return args, err
	}

	// Don't bother with security on windows
	if runtime.GOOS != "windows" && configFileInfo.Mode() != 0600 {
		return args, errors.New(fmt.Sprintf("Permissions for %v were %v and not -rw-------. Please update (e.g. chmod 0600 %v) as sensitive information might be stored in the config.", m.RunContext.BaseConfigurationPath, configFileInfo.Mode(), m.RunContext.BaseConfigurationPath))
	}

	// Finally, load config from file (even if we just wrote it)
	config, err := LoadConfigFrom(configPath)
	if err != nil {
		return args, err
	}

	m.Config = config

	err = context.Finalize()
	if err != nil {
		return args, err
	}

	return args, nil
}

// Create default config at the configuration path of the run context
func (m *Meta) CreateDefaultConfig() error {
	configDir := filepath.Dir(m.RunContext.BaseConfigurationPath)
	if !utility.FileExists(configDir) {
		m.Ui.Info(fmt.Sprintf("Creating %v", configDir))
		err := os.Mkdir(configDir, 0700)
		if err != nil {
			return err
		}
	}

	m.Ui.Info(fmt.Sprintf("Creating %v", m.RunContext.BaseConfigurationPath))
	file, err := os.Create(m.RunContext.BaseConfigurationPath)
	if err != nil {
		return err
	}

	// Chmod doesn't work on Windows, and there appears to be no equivalent
	if runtime.GOOS != "windows" {
		err = file.Chmod(0600)
		if err != nil {
			return err
		}
	}

	viperConfig := NewDefaultConfig()
	err = SaveConfigTo(*viperConfig, m.RunContext.BaseConfigurationPath)
	if err != nil {
		return err
	}

	return nil
}

func (m *Meta) Colorize() *colorstring.Colorize {
	return &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: !m.color,
		Reset:   true,
	}
}
