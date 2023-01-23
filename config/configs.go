package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App       AppConfig       `yaml:"app"`
	Http      HttpConfig      `yaml:"http"`
	Pg        PgConfig        `yaml:"postgres"`
	Cache     CacheConfig     `yaml:"cache"`
	Redis     RedisConfig     `yaml:"redis"`
	Memcached MemcachedConfig `yaml:"memcached"`
}

type AppConfig struct {
	Name    string `env-required:"true" env:"APP_NAME" yaml:"name"`
	Version string `env-required:"true" env:"APP_VERSION" yaml:"version"`
}

type HttpConfig struct {
	Host string `env-required:"true" env:"HTTP_HOST" yaml:"host"`
	Port string `env-required:"true" env:"HTTP_PORT" yaml:"port"`
}

type PgConfig struct {
	Host   string `env-required:"true" env:"PG_HOST" yaml:"host"`
	Port   string `env-required:"true" env:"PG_PORT" yaml:"port"`
	User   string `env-required:"true" env:"PG_USER" yaml:"user"`
	Pass   string `env-required:"true" env:"PG_PASS" yaml:"pass"`
	DbName string `env-required:"true" env:"PG_DBNAME" yaml:"dbname"`
}

type CacheConfig struct {
	Driver string `env:"CACHE_DRIVER" yaml:"driver"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" yaml:"host"`
	Port     string `env:"REDIS_PORT" yaml:"port"`
	Password string `env:"REDIS_PASSWORD" yaml:"password"`
}

type MemcachedConfig struct {
	Hosts []string `env-separator:"|" env:"MEMCACHED_HOSTS" yaml:"hosts"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == "prod" {
		err = cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	} else {
		err = cleanenv.ReadConfig(dir+"/config.local.yaml", cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
