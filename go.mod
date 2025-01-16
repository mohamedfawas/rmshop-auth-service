module github.com/mohamedfawas/rmshop-auth-service

go 1.23.4

require (
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/lib/pq v1.10.9
	github.com/mohamedfawas/rmshop-proto v0.0.0-20250115221258-ce447a396e83
	golang.org/x/crypto v0.32.0
	google.golang.org/grpc v1.69.4
)

require (
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)

// replace github.com/mohamedfawas/rmshop-proto => ../rmshop-proto
