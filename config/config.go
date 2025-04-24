package config

import (
	"os"
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

func GetConfig() Config {
	var config Config
	dat, err := os.ReadFile("config.toml")
	util.Check(err)
	_, err = toml.Decode(string(dat), &config)
	util.Check(err)
	return config
}
