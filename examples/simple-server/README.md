# Simple server

An example of a simple NATS server.

## Usage

Start the server:

```bash
$ go run main.go
```

Subscribe to a subject and publish a message to it:

```bash
$ PIN=$(date +"%Y%m%d%H%M%S")

$ nats sub test-subject --count=1 | grep "PIN: $PIN" &
$ nats pub test-subject "PIN: $PIN"
```
