module github.com/radiophysiker/microservices-homework/iam

go 1.25.3

replace github.com/radiophysiker/microservices-homework/shared => ../shared

replace github.com/radiophysiker/microservices-homework/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/gomodule/redigo v1.9.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/joho/godotenv v1.5.1
	github.com/radiophysiker/microservices-homework/platform v0.0.0-20251112151515-a870437b7b54
	github.com/radiophysiker/microservices-homework/shared v0.0.0-20251112151515-a870437b7b54
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.45.0
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pressly/goose/v3 v3.26.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
)
