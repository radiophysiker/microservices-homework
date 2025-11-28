package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	Outputs() []string
	OTELCollectorEndpoint() string
	ServiceName() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	ServiceVersion() string
	Environment() string
}

type PaymentGRPCConfig interface {
	Address() string
}
