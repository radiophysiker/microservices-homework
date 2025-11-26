package env

import (
	"github.com/caarlos0/env/v11"
)

type migrationsEnvConfig struct {
	Directory string `env:"MIGRATION_DIRECTORY" envDefault:"./migrations"`
}

type migrationsConfig struct {
	raw migrationsEnvConfig
}

func NewMigrationsConfig() (*migrationsConfig, error) {
	var raw migrationsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &migrationsConfig{raw: raw}, nil
}

func (cfg *migrationsConfig) Directory() string {
	return cfg.raw.Directory
}
