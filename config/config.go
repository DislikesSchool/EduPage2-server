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
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
		TTL      struct {
			Timeline  int `yaml:"timeline"`
			Timetable int `yaml:"timetable"`
			Results   int `yaml:"results"`
			DBI       int `yaml:"dbi"`
		} `yaml:"ttl"`
	} `yaml:"redis"`
	Database struct {
		Enabled bool   `yaml:"enabled"`
		Driver  string `yaml:"driver"`
		DSN     string `yaml:"dsn"`
	} `yaml:"database"`
	Encryption struct {
		Enabled bool   `yaml:"enabled"`
		Key     string `yaml:"key"`
	} `yaml:"encryption"`
	Meilisearch struct {
		Enabled  bool   `yaml:"enabled"`
		Host     string `yaml:"host"`
		APIKey   string `yaml:"api_key"`
		Messages struct {
			IndexName  string `yaml:"index_name"`
			PrimaryKey string `yaml:"primary_key"`
		} `yaml:"messages"`
	} `yaml:"meilisearch"`
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

	if encKey := os.Getenv("ENCRYPTION_KEY"); encKey != "" {
		AppConfig.Encryption.Key = encKey
		AppConfig.Encryption.Enabled = true
	}
}
