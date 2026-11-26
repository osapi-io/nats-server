# Authentication with an NKey

An example NATS server that uses [NKey Auth].

## Setup

Make a user NKey for each of two services. Each file has the seed on the first
line and the public key on the second. The example reads the public keys from
the files. Git ignores `.nkeys/`, because a seed is a private key.

```bash
$ mkdir -p .nkeys
$ go tool nk -gen user -pubout > .nkeys/service1.nk
$ go tool nk -gen user -pubout > .nkeys/service2.nk
```

## Usage

Start the server:

```bash
$ go run main.go
```

> **Important:** Configure a system account in the NATS server before you do
> this query. Without a system account, system queries such as
> `nats server info` do not work. This is also true with NKey authentication.

Query the server with the system user:

```bash
$ nats server info --user system --password systempassword
```

Subscribe to a subject. Then publish a message to it:

```bash
$ PIN=$(date +"%Y%m%d%H%M%S")

$ nats sub test-subject --count=1 --nkey ./.nkeys/service1.nk | grep "PIN: $PIN" &
$ nats pub test-subject "PIN: $PIN" --nkey ./.nkeys/service1.nk

$ nats pub test-subject "PIN: $PIN" --nkey ./.nkeys/service2.nk # fail
```

[nkey auth]: https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/nkey_auth
