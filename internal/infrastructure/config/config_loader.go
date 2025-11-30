package config

import "github.com/spf13/viper"

type Config struct {
	DatabaseURI     string `mapstructure:"DATABASE_URI"`
	DatabaseName    string `mapstructure:"DATABASE_NAME"`
	ServerPort      string `mapstructure:"SERVER_PORT"`
	JWTSecret       string `mapstructure:"JWT_SECRET"`
	JWTIssuer       string `mapstructure:"JWT_ISSUER"`
	CasbinModelPath string `mapstructure:"CASBIN_MODEL_PATH"`
	Production      bool   `mapstructure:"PRODUCTION"`
	UserCollection       string `mapstructure:"USER_COLLECTION"`
	AssetCollection      string `mapstructure:"ASSET_COLLECTION"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

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
