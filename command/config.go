package command

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func NewDefaultConfig() *viper.Viper {
	config := viper.New()
	config.Set("getstyle", "folder")
	config.Set("submitstyle", "smtp")
	config.Set("smtpsettings", Smtpsettings{
		From:     "your@email.com",
		Port:     "587",
		Server:   "smtp.gmail.com",
		Username: "your.login@email.com",
		Password: "",
	})

	return config
}

type ConfigStruct struct {
	Submitstyle  string       `json:"submitstyle,omitempty"`
	Getstyle     string       `json:"getstyle,omitempty"`
	Smtpsettings Smtpsettings `json:"smtpsettings,omitempty"`
	Imapsettings Imapsettings `json:"imapsettings,omitempty"`
}

type Smtpsettings struct {
	From     string `json:"from,omitempty"`
	Server   string `json:"server,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Imapsettings struct {
	Server   string `json:"server,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

var DefaultConfigStruct ConfigStruct

func LoadConfigFrom(configPath string) (ConfigStruct, error) {
	config := ConfigStruct{}
	viper := viper.New()

	basename := filepath.Base(configPath)
	viper.SetConfigName(strings.TrimSuffix(basename, filepath.Ext(basename)))
	viper.AddConfigPath(filepath.Dir(configPath))

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func SaveConfigTo(viperConfig viper.Viper, path string) error {
	err := viperConfig.Unmarshal(&DefaultConfigStruct)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(DefaultConfigStruct, "", "    ")
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	f.WriteString(string(b))

	return nil
}
