package config

import (
	"flag"
	"fmt"
	"github.com/aakosarev/kanban-board/back/pkg/constants"
	"github.com/aakosarev/kanban-board/back/pkg/logger"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "kanban microservice config path")
}

type Config struct {
	ServiceName string         `mapstructure:"serviceName"`
	Http        Http           `mapstructure:"http"`
	Cookie      Cookie         `mapstructure:"cookie"`
	Session     Session        `mapstructure:"session"`
	Postgres    Postgres       `mapstructure:"postgres"`
	Redis       Redis          `mapstructure:"redis"`
	Logger      *logger.Config `mapstructure:"logger"`
}

type Http struct {
	Port                string `mapstructure:"port" validate:"required"`
	Development         bool   `mapstructure:"development"`
	BasePath            string `mapstructure:"basePath" validate:"required"`
	AuthPath            string `mapstructure:"authPath" validate:"required"`
	ColumnPath          string `mapstructure:"columnPath" validate:"required"`
	TaskPath            string `mapstructure:"taskPath" validate:"required"`
	BoardPath           string `mapstructure:"boardPath" validate:"required"`
	DebugErrorsResponse bool   `mapstructure:"debugErrorsResponse"`
}

type Session struct {
	Prefix string `mapstructure:"prefix"`
	Name   string `mapstructure:"name"`
	Expire int    `mapstructure:"expire"`
}

type Cookie struct {
	MaxAge   int  `mapstructure:"maxAge"`
	Secure   bool `mapstructure:"secure"`
	HTTPOnly bool `mapstructure:"httpOnly"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type Redis struct {
	RedisAddr      string `mapstructure:"redisAddr"`
	RedisPassword  string `mapstructure:"redisPassword"`
	RedisDB        string `mapstructure:"redisDB"`
	RedisDefaultDB string `mapstructure:"redisDefaultDB"`
	MinIdleConns   int    `mapstructure:"minIdleConns"`
	PoolSize       int    `mapstructure:"poolSize"`
	PoolTimeout    int    `mapstructure:"poolTimeout"`
	Password       string `mapstructure:"password"`
	DB             int    `mapstructure:"db"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/internal/config/config.yaml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, nil
}
