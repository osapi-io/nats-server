# Authenticating with a user and password

An example NATS server that uses [User and Password Auth].

## Usage

Start the server:

```bash
$ go run main.go
```

Query the server with the system user:

```bash
$ nats server info --user system --password systempassword
```

Subscribe to a subject and publish a message to it:

```bash
$ PIN=$(date +"%Y%m%d%H%M%S")

$ nats sub test-subject --count=1 --user myuser --password mypassword | grep "PIN: $PIN" &
$ nats pub test-subject "PIN: $PIN" --user myuser --password mypassword
```

[user and password auth]: https://docs.nats.io/using-nats/developer/connecting/userpass
