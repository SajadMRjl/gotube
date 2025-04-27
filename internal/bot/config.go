package bot

import (
	"github.com/spf13/viper"
)

type Config struct {
	Telegram struct {
		Token   string
		Debug   bool
		Timeout int
	}
	Server struct {
		Port       string
		WebhookURL string `mapstructure:"webhook_url"`
	}
	Logging struct {
		Level       string
		Development bool
	}
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
