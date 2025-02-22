# No Authentication

An example NATS server without Auth.

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

$ nats sub test-subject --count=1 | grep "PIN: $PIN" &
$ nats pub test-subject "PIN: $PIN"
```
