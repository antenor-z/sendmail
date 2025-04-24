package config

import (
	"os"
	"path/filepath"
	"sendmail/util"

	"github.com/BurntSushi/toml"
)

type Mailgun struct {
	ApiKey string `toml:"apiKey"`
	Sender string `toml:"sender"`
	Domain string `toml:"domain"`
}

type Config struct {
	Mailgun Mailgun `toml:"mailgun"`
}

func getConfigPath() string {
	configDir := filepath.Join(os.Getenv("HOME"), ".config")
	return filepath.Join(configDir, "a4sendmail", "config.toml")
}

func GetConfig() Config {
	var config Config
	dat, err := os.ReadFile(getConfigPath())
	util.Check(err)
	_, err = toml.Decode(string(dat), &config)
	util.Check(err)
	return config
}
