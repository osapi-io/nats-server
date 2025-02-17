# Authenticating with a User and Password

An example simple NATS server using [User and Password Auth][].

## Usage

Start the server:

```bash
$ go run main.go
```

Query the server with the system user:

```bash
$ nats server info --user system --password systempassword
```

[User and Password Auth]: https://docs.nats.io/using-nats/developer/connecting/userpass
