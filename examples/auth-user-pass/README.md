# Authenticating with a User and Password

An example NATS server using [User and Password Auth][].

## Usage

Start the server:

```bash
$ go run main.go
```

Query the server with the system user:

```bash
$ nats server info --user system --password systempassword
```

Subscribe and Publish a message:

```bash
$ PIN=$(date +"%Y%m%d%H%M%S")

$ nats sub test-subject --count=1 --user myuser --password mypassword | grep "PIN: $PIN" &
$ nats pub test-subject "PIN: $PIN" --user myuser --password mypassword
```

[User and Password Auth]: https://docs.nats.io/using-nats/developer/connecting/userpass
