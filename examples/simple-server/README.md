# Simple Server

An example simple NATS server.

## Usage

Start the server:

```bash
go run main.go
```

Subscribe and Publish a message:

```bash
PIN=$(date +"%Y%m%d%H%M%S")

nats sub test-subject >output.txt &
SUB_PID=$!

nats pub test-subject "PIN: $PIN"

grep "PIN: $PIN" output.txt

kill $SUB_PID
```

[User and Password Auth]: https://docs.nats.io/using-nats/developer/connecting/userpass
