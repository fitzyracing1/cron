# Architecture

| Path | Role |
|---|---|
| `cron.go` | `Cron` scheduler, entries, run loop |
| `parser.go` | spec parsing (the fix lives here) |
| `spec.go, constantdelay.go` | schedule types |
| `chain.go, option.go, logger.go` | job wrappers, options, logging |
| `doc.go` | package documentation |


