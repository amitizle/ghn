package config

import (
	"github.com/amitizle/ghn/pkg/notifiers"
	"github.com/spf13/viper"
)

// Config is a struct that holds the configuration that is
// used for the application. It includes logger and notifiers configuration.
type Config struct {
	Notifiers []*notifierInstace `yaml:"notifiers"`
	Log       *LogConfig         `yaml:"log"`
	Github    *GithubConfig      `yaml:"github"`
}

type notifierInstace struct {
	Notifier notifiers.Notifier

	Type   string                 `yaml:"type"`
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

// LogConfig contains the configuration for the logger, e.g. level
type LogConfig struct {
	Level string `yaml:"level"`
}

// GithubConfig represents the configuration for Github.
// In particular, it includes the oauth token
type GithubConfig struct {
	Token string `yaml:"string"`
}

func init() {
	viper.SetDefault("log.level", "debug")
}

// New returns a new Config struct
func New() *Config {
	return &Config{}
}
