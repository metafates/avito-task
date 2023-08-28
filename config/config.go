package config

import (
	"errors"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Port string   `koanf:"port"`
	DB   DBConfig `koanf:"db"`
}

type DBConfig struct {
	PostgresURI string `koanf:"postgres"`
}

func LoadEnvs(filenames ...string) error {
	return godotenv.Load(filenames...)
}

func Load() (config Config, err error) {
	k := koanf.New(".")

	// Default values
	err = k.Load(confmap.Provider(map[string]any{
		"port": "1234",
	}, "."), nil)

	if err != nil {
		return
	}

	const envPrefix = "SERVER_"
	err = k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.ReplaceAll(
			strings.ToLower(strings.TrimPrefix(s, envPrefix)),
			"_",
			".",
		)
	}), nil)

	if err != nil {
		return
	}

	err = k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{
		Tag:       "koanf",
		FlatPaths: false,
	})

	if err != nil {
		return
	}

	if config.DB.PostgresURI == "" {
		err = errors.New("postgres uri is empty")
		return
	}

	return
}
