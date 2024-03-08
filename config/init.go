package config

import (
	"bytes"
	_ "embed"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

//go:embed config.yml
var defaultConfiguration []byte

type Auth struct {
	HashSalt   string        `mapstructure:"hash_salt"`
	SigningKey string        `mapstructure:"signing_key"`
	TokenTTL   time.Duration `mapstructure:"token_ttl"`
}

type Postgres struct {
	Host     string
	User     string
	Password string
	DB       string
	Port     string
}

type Config struct {
	Port     string
	Auth     *Auth
	Postgres *Postgres
}

func Init() (*Config, error) {
	viper.SetEnvPrefix("CONNECT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfiguration)); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	if err := initSecrets(config); err != nil {
		return nil, err
	}

	return config, nil
}

func initSecrets(config *Config) error {
	hashSaltFile := viper.GetString("auth.hash_salt_file")
	if hashSaltFile != "" {
		value, err := os.ReadFile(hashSaltFile)
		if err != nil {
			return err
		}
		config.Auth.HashSalt = string(value)
	}

	signingKeyFile := viper.GetString("auth.signing_key_file")
	if signingKeyFile != "" {
		value, err := os.ReadFile(signingKeyFile)
		if err != nil {
			return err
		}
		config.Auth.SigningKey = string(value)
	}

	var dbFile = viper.GetString("postgres.db_file")
	if dbFile != "" {
		value, err := os.ReadFile(dbFile)
		if err != nil {
			return err
		}
		config.Postgres.DB = string(value)
	}

	passwordFile := viper.GetString("postgres.password_file")
	if passwordFile != "" {
		value, err := os.ReadFile(passwordFile)
		if err != nil {
			return err
		}
		config.Postgres.Password = string(value)
	}

	userFile := viper.GetString("postgres.user_file")
	if userFile != "" {
		value, err := os.ReadFile(userFile)
		if err != nil {
			return err
		}
		config.Postgres.User = string(value)
	}

	return nil
}
