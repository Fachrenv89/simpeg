package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment         string        `mapstructure:"ENVIRONMENT"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress   string        `mapstructure:"PORT"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables.
// func LoadConfig(path string) (config Config, err error) {
// 	viper.AddConfigPath(path)
// 	viper.SetConfigName("app")
// 	viper.SetConfigType("env")
// 	viper.AutomaticEnv()
// 	viper.ReadInConfig()
// 	if err != nil {
// 		config.DBDriver = viper.GetString("DB_DRIVER")
// 		config.Environment = viper.GetString("ENVIRONMENT")
// 		config.DBSource = viper.GetString("DB_SOURCE")
// 		config.HTTPServerAddress = viper.GetString("PORT")
// 		config.TokenSymmetricKey = viper.GetString("TOKEN_SYMMETRIC_KEY")
// 		config.AccessTokenDuration = viper.GetDuration("ACCESS_TOKEN_DURATION")
// 		return
// 	}

// 	err = viper.Unmarshal(&config)
// 	return
// }


func LoadConfig(path string) (config Config, err error) {
	// viper.AddConfigPath(path)
	// viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	// if err != nil {
		config.DBDriver = viper.GetString("DB_DRIVER")
		config.Environment = viper.GetString("ENVIRONMENT")
		config.DBSource = viper.GetString("DB_SOURCE")
		config.HTTPServerAddress = viper.GetString("PORT")
		config.TokenSymmetricKey = viper.GetString("TOKEN_SYMMETRIC_KEY")
		config.AccessTokenDuration = viper.GetDuration("ACCESS_TOKEN_DURATION")
		// return
	// }

	// err = viper.Unmarshal(&config)
	return
}
