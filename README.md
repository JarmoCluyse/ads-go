# ads-go

A Go client library for communicating with Beckhoff devices using the ADS protocol.

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

## Features

- Connect to Beckhoff devices over ADS protocol
- Read and write ADS
- Read and write system state

---

## Installation

```bash
go get github.com/jarmocluyse/ads-go@latest
```

---

## Usage

Usage examples and sample code can be found in [`cmd/main.go`](./cmd/main.go).

---

## Logging

By default, the client uses a **silent logger** (no logs).  
To enable logging, pass your own `*slog.Logger` to `NewClient`.

Example:

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
client := ads.NewClient(settings, logger)
```

---

## Contributing

Contributions, issues, and feature requests are welcome!  
Feel free to check [issues page](https://github.com/jarmocluyse/beckhoffads/issues).

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
