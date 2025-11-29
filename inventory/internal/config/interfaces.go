package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	Outputs() []string
	OTELCollectorEndpoint() string
	ServiceName() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type IAMGRPCConfig interface {
	IAMAddress() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
