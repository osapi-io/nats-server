# Simple Server

An example simple NATS server.

## Usage

Start the server:

```bash
$ go run main.go
```

Subscribe and Publish a message:

```bash
$ PIN=$(date +"%Y%m%d%H%M%S")

$ nats sub test-subject --count=1 | grep "PIN: $PIN" &
$ nats pub test-subject "PIN: $PIN"
```

[User and Password Auth]: https://docs.nats.io/using-nats/developer/connecting/userpass
