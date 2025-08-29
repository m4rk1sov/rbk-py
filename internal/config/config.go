package config

import (
	"flag"
	"os"
	"time"
	
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env               string     `json:"env" yaml:"env" env-default:"dev"`
	JWT               JWTConfig  `yaml:"jwt"`
	HTTP              HTTPConfig `yaml:"http"`
	StaticToken       string     `yaml:"static_token" env:"STATIC_TOKEN" env-default:"static_token"`
	TemplateDir       string     `yaml:"template_dir" env:"TEMPLATE_DIR" env-default:"./templates"`
	ServiceContextURL string     `yaml:"service_context_url" env-default:"/document-generator"`
	PDFConverterURL   string     `yaml:"pdf_converter_url" env-default:"http://gotenberg:3000"`
}

type JWTConfig struct {
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
	Secret   string        `yaml:"secret" env:"SECRET_KEY" env-default:"jwt_secret"`
}

type HTTPConfig struct {
	Port         int           `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read-timeout" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write-timeout" env-default:"15s"`
	IdleTimeout  time.Duration `yaml:"idle-timeout" env-default:"2m"`
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
