package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Environment         string        `mapstructure:"ENVIRONMENT"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	MigrationURL        string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return

		config.Environment = viper.GetString("ENVIRONMENT")
		config.DBDriver = viper.GetString("DB_DRIVER")
		config.DBSource = viper.GetString("DB_SOURCE")
		config.MigrationURL = viper.GetString("MIGRATION_URL")
		config.HTTPServerAddress = viper.GetString("HTTP_SERVER_ADDRESS")
		config.TokenSymmetricKey = viper.GetString("TOKEN_SYMMETRIC_KEY")
		config.AccessTokenDuration = viper.GetDuration("ACCESS_TOKEN_DURATION")
	}

	err = viper.Unmarshal(&config)
	return
}
