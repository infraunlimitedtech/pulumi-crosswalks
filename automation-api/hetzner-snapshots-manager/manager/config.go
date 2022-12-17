package manager

import (
	"github.com/spf13/viper"
)

type Config struct {
	Debug         bool
	APIServerPort int `mapstructure:"api-server-port"`
	Stack         Stack
}

type Stack struct {
	Name string
	Path string
}

func GetConfig() (*Config, error) {
	var config *Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
