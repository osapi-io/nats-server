# No authentication

An example NATS server with no authentication.

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

$ nats sub test-subject --count=1 | grep "PIN: $PIN" &
$ nats pub test-subject "PIN: $PIN"
```
