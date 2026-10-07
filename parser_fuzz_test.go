//go:build go1.18
// +build go1.18

package cron

import (
	"testing"
	"time"
)

// FuzzParse checks that no input makes the parser (or Next on the result)
// panic. `go test` runs the seed corpus; `go test -fuzz=FuzzParse` explores.
func FuzzParse(f *testing.F) {
	seeds := []string{
		"* * * * *", "0 0 0 1 1 ?", "*/5 1-3,7 ? jan-mar mon-fri", "1-5/2 * * * *",
		"@every 1h", "@daily", "TZ=UTC 5 * * * *", "CRON_TZ=Asia/Tokyo @daily",
		"TZ=UTC ", "TZ=a b", "-", "*/", "1-", "a-b/c",
	}
	seeds = append(seeds, tzWithoutSchedule...)
	for _, s := range seeds {
		f.Add(s)
	}
	parsers := []Parser{
		standardParser,
		secondParser,
		NewParser(SecondOptional | Minute | Hour | Dom | Month | Dow | Descriptor),
		NewParser(Dom | Month | DowOptional),
	}
	from := time.Date(2026, 10, 7, 0, 0, 0, 0, time.UTC)
	f.Fuzz(func(t *testing.T, spec string) {
		for _, p := range parsers {
			if sched, err := p.Parse(spec); err == nil {
				sched.Next(from)
			}
		}
	})
}
