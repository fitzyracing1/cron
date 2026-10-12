# github.com/fitzyracing1/cron/v3

`github.com/fitzyracing1/cron/v3` is a drop-in fork of [robfig/cron](https://github.com/robfig/cron) and part of [360 Bench](https://github.com/fitzyracing1/360-bench). It is based on upstream `github.com/robfig/cron/v3` v3.0.1, and the API is unchanged. All credit for the original library goes to its authors and contributors.

## What the fork fixes (summary of [README.md](../README.md))

- No more `panic: slice bounds out of range [:-1]` when a schedule is only a timezone prefix (`TZ=UTC`, `CRON_TZ=...`). It now returns `missing schedule after timezone: <spec>`.

The full details and links to the upstream issues are in [README.md](../README.md).

Pages: [Getting started](getting-started.md) · [Architecture](architecture.md) · [Contributing](contributing.md)

## CI

`.github/workflows/test.yml` runs `go vet ./...`, `go test ./... -race` and a 30-second `FuzzParse` fuzz run. A run is **lit** when it succeeds and **dark** otherwise.

## License

MIT. The original notice is kept: Copyright (C) 2012 Rob Figueiredo.
Modifications: Copyright (c) 2026 Joshua Almeida.
