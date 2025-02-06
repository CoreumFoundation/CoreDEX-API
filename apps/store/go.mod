module github.com/CoreumFoundation/CoreDEX-API/apps/store

go 1.23.3

require (
	github.com/CoreumFoundation/CoreDEX-API/domain v0.0.0-20250204222705-64b06c939bc4
	github.com/CoreumFoundation/CoreDEX-API/utils v0.0.0-20250204222705-64b06c939bc4
	google.golang.org/grpc v1.69.0
	google.golang.org/protobuf v1.36.0
)

replace github.com/CoreumFoundation/CoreDEX-API/domain => ../../domain

replace github.com/CoreumFoundation/CoreDEX-API/utils => ../../utils

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241209162323-e6fa225c2576 // indirect
)
