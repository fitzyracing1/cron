# Getting started

## Install

```
# keep imports and fix transitive users too (go.mod):
replace github.com/robfig/cron/v3 => github.com/fitzyracing1/cron/v3 v3.0.2
# or switch imports:
go get github.com/fitzyracing1/cron/v3@v3.0.2
```



## Usage

Same API as robfig/cron v3: `cron.New()`, `AddFunc`, `AddJob`, `Start`, `Stop`, `ParseStandard`, plus parser options. Package docs are in `doc.go`.

## Develop and test

```
go vet ./...
go test ./... -race
go test -run='^$' -fuzz=FuzzParse -fuzztime=30s
```
