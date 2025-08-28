package config

import (
	"flag"
	"os"
	"time"
	
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string     `json:"env" yaml:"env" env-default:"dev"`
	JWT  JWTConfig  `yaml:"jwt"`
	HTTP HTTPConfig `yaml:"http"`
}

type JWTConfig struct {
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type HTTPConfig struct {
	Port    int           `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"1h"`
}

func Load() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("Config path is empty")
	}
	
	return loadPath(configPath)
}

func fetchConfigPath() string {
	var res string
	
	flag.StringVar(&res, "config", "", "path to config file")
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	flag.Parse()
	
	return res
}

func loadPath(configPath string) *Config {
	// check for file existence
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config path is empty: " + configPath)
	}
	
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}
	
	return &cfg
}
