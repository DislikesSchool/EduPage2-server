package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	Schools struct {
		Whitelist   []string `yaml:"whitelist"`
		IsBlacklist bool     `yaml:"blacklist" mapstructure:"blacklist"`
	} `yaml:"schools"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	Redis struct {
		Enabled  bool   `yaml:"enabled"`
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
		TTL      struct {
			Timeline  int `yaml:"timeline"`
			Timetable int `yaml:"timetable"`
			Results   int `yaml:"results"`
			DBI       int `yaml:"dbi"`
		} `yaml:"ttl"`
	} `yaml:"redis"`
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

	if hostEnv := os.Getenv("HOST"); hostEnv != "" {
		AppConfig.Server.Host = hostEnv
	}

	if portEnv := os.Getenv("PORT"); portEnv != "" {
		AppConfig.Server.Port = portEnv
	}
}
