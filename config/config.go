package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		os.Exit(1)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Unable to decode config into struct: %v", err)
		os.Exit(1)
	}

	AppConfig = &cfg

	if jwtEnv := os.Getenv("JWT_SECRET_KEY"); jwtEnv != "" {
		AppConfig.JWT.Secret = jwtEnv
	}
}
