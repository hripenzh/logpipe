# logpipe

A lightweight CLI for tailing and filtering structured logs from multiple sources simultaneously.

---

## Installation

```bash
go install github.com/yourusername/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logpipe.git && cd logpipe && go build -o logpipe .
```

---

## Usage

Tail logs from multiple sources and filter by log level or field:

```bash
# Tail a single log file
logpipe tail --file /var/log/app.log

# Tail multiple sources simultaneously
logpipe tail --file /var/log/app.log --file /var/log/worker.log

# Filter by log level
logpipe tail --file /var/log/app.log --level error

# Filter by a specific field value
logpipe tail --file /var/log/app.log --filter service=api-gateway

# Combine filters
logpipe tail --file /var/log/app.log --level warn --filter env=production
```

Output is colorized and formatted for readability in your terminal.

---

## Flags

| Flag | Description |
|------|-------------|
| `--file` | Path to a log file (repeatable) |
| `--level` | Minimum log level to display (`debug`, `info`, `warn`, `error`) |
| `--filter` | Filter by field value in `key=value` format |
| `--json` | Output raw JSON instead of formatted output |

---

## License

MIT © [yourusername](https://github.com/yourusername)