# Authenticating with an NKey

An example NATS server using [NKey Auth][].

## Setup

Create a user NKEY for two micro-services.

```bash
$ nk -gen user -pubout > .nkeys/service1.nk
$ nk -gen user -pubout > .nkeys/service2.nk
```

## Usage

Start the server:

```bash
$ go run main.go
```

Query the server with the system user:

```bash
$ nats server info --user system --password systempassword
```

> **Important:** This will only work if a **system account is configured** in
  the NATS server. Without a properly configured system account, system-level
  queries like `nats server info` will not be available, even when using NKEYS
  authentication.

Subscribe and Publish a message:

```bash
$ PIN=$(date +"%Y%m%d%H%M%S")

$ nats sub test-subject --count=1 --nkey ./.nkeys/service1.nk | grep "PIN: $PIN" &
$ nats pub test-subject "PIN: $PIN" --nkey ./.nkeys/service1.nk

$ nats pub test-subject "PIN: $PIN" --nkey ./.nkeys/service2.nk # fail
```

[NKey Auth]: https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/nkey_auth
