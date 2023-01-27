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
	Firebase  FirebaseConfig  `yaml:"firebase"`
}

type AppConfig struct {
	Name    string `env-required:"true" env:"APP_NAME" yaml:"name"`
	Version string `env-required:"true" env:"APP_VERSION" yaml:"version"`
	Debug   bool   `env:"APP_DEBUG" yaml:"debug" env-default:"false"`
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

type FirebaseConfig struct {
	Type                    string `env:"FIREBASE_TYPE" yaml:"type" json:"type"`
	ProjectId               string `env:"FIREBASE_PROJECT_ID" yaml:"project_id" json:"project_id"`
	PrivateKeyId            string `env:"FIREBASE_PRIVATE_KEY_ID" yaml:"private_key_id" json:"private_key_id"`
	PrivateKey              string `env:"FIREBASE_PRIVATE_KEY" yaml:"private_key" json:"private_key"`
	ClientEmail             string `env:"FIREBASE_CLIENT_EMAIL" yaml:"client_email" json:"client_email"`
	ClientId                string `env:"FIREBASE_CLIENT_ID" yaml:"client_id" json:"client_id"`
	AuthURI                 string `env:"FIREBASE_AUTH_URI" yaml:"auth_uri" json:"auth_uri"`
	TokenURI                string `env:"FIREBASE_TOKEN_URI" yaml:"token_uri" json:"token_uri"`
	AuthProviderX509CertURL string `env:"FIREBASE_PROVIDER_X509_CERT_URL" yaml:"auth_provider_x509_cert_url" json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `env:"FIREBASE_CLIENT_X509_CERT_URL" yaml:"client_x509_cert_url" json:"client_x509_cert_url"`
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
