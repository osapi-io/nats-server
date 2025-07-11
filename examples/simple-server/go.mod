module example.com/server

go 1.24.0

replace github.com/osapi-io/nats-server => ../../../nats-server/

replace github.com/osapi-io/nats-client => ../../../nats-client/

require (
	github.com/nats-io/nats-server/v2 v2.11.4
	github.com/osapi-io/nats-server v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/go-tpm v0.9.5 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	github.com/nats-io/jwt/v2 v2.7.4 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/time v0.11.0 // indirect
)
