# Simple Server

An example simple NATS server.

## Usage

Start the server:

```bash
$ go run main.go
```

Subscribe and Publish a message:

```bash
$ nats sub test-subject --count=1 | grep "PIN: $PIN" &
$ SUB_PID=$!

$ nats pub test-subject "PIN: $PIN"

$ wait $SUB_PID
```

[User and Password Auth]: https://docs.nats.io/using-nats/developer/connecting/userpass
