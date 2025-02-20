module example.com/server

go 1.24

toolchain go1.24.0

replace github.com/osapi-io/nats-server => ../../../nats-server/

replace github.com/osapi-io/nats-client => ../../../nats-client/

require (
	github.com/lmittmann/tint v1.0.7
	github.com/nats-io/nats-server/v2 v2.10.25
	github.com/nats-io/nats.go v1.39.0
	github.com/osapi-io/nats-client v0.0.0-00010101000000-000000000000
	github.com/osapi-io/nats-server v0.0.0-00010101000000-000000000000
)

require (
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	github.com/nats-io/jwt/v2 v2.7.3 // indirect
	github.com/nats-io/nkeys v0.4.10 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/time v0.9.0 // indirect
)

tool github.com/nats-io/nkeys/nk
