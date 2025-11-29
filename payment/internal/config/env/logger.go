package env

import (
	"strings"

	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level                 string `env:"LOGGER_LEVEL,required"`
	AsJson                bool   `env:"LOGGER_AS_JSON,required"`
	Outputs               string `env:"LOG_OUTPUTS" envDefault:"stdout"`
	OTELCollectorEndpoint string `env:"OTEL_COLLECTOR_ENDPOINT" envDefault:""`
	ServiceName           string `env:"SERVICE_NAME,required"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{raw: raw}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *loggerConfig) AsJson() bool {
	return cfg.raw.AsJson
}

func (cfg *loggerConfig) AsJSON() bool {
	return cfg.raw.AsJson
}

func (cfg *loggerConfig) Outputs() []string {
	parts := strings.Split(cfg.raw.Outputs, ",")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}

		out = append(out, trimmed)
	}

	return out
}

func (cfg *loggerConfig) OTELCollectorEndpoint() string {
	return cfg.raw.OTELCollectorEndpoint
}

func (cfg *loggerConfig) ServiceName() string {
	return cfg.raw.ServiceName
}
