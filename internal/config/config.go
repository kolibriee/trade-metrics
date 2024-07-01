package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	ClickHouse ClickHouse
	Server     Server `mapstructure:"server"`
}

type Server struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

type ClickHouse struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func New(path string, fileName string) (*Config, error) {
	var cfg *Config
	godotenv.Load()
	viper.SetConfigName(fileName)
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New("failed to read config file: " + err.Error())
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.New("failed to unmarshal config: " + err.Error())
	}
	if err := envconfig.Process("DB", &cfg.ClickHouse); err != nil {
		return nil, errors.New("failed to process env variables: " + err.Error())
	}
	fmt.Println(os.Getenv("DB_HOST"))
	return cfg, nil
}
