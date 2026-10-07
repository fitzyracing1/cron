package cron

import (
	"strings"
	"testing"
	"time"
)

// A timezone prefix with no schedule after it used to panic with
// "slice bounds out of range [:-1]" instead of returning an error.
// robfig/cron#554 ("TZ=0") and #470 ("TZ=TZ=TZ=").
var tzWithoutSchedule = []string{
	"TZ=0",
	"TZ=UTC",
	"TZ=",
	"TZ=TZ=TZ=",
	"CRON_TZ=",
	"CRON_TZ=America/New_York",
	"TZ=UTC\t*\t*\t*\t*\t*", // tabs only; the prefix has always needed a space
}

func parseNoPanic(t *testing.T, p Parser, spec string) (sched Schedule, err error) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Parse(%q) panicked: %v", spec, r)
		}
	}()
	return p.Parse(spec)
}

func TestParseTimezoneWithoutScheduleReturnsError(t *testing.T) {
	parsers := map[string]Parser{
		"standard": standardParser,
		"second":   secondParser,
		"optional-second": NewParser(
			SecondOptional | Minute | Hour | Dom | Month | Dow | Descriptor),
	}
	for name, p := range parsers {
		for _, spec := range tzWithoutSchedule {
			sched, err := parseNoPanic(t, p, spec)
			if err == nil {
				t.Errorf("%s: Parse(%q) = %v, want an error", name, spec, sched)
				continue
			}
			if !strings.Contains(err.Error(), "missing schedule after timezone") {
				t.Errorf("%s: Parse(%q) error = %q", name, spec, err)
			}
			if sched != nil {
				t.Errorf("%s: Parse(%q) returned a schedule with the error", name, spec)
			}
		}
	}
}

func TestParseStandardTZ0(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ParseStandard(\"TZ=0\") panicked: %v", r)
		}
	}()
	if _, err := ParseStandard("TZ=0"); err == nil {
		t.Fatal(`ParseStandard("TZ=0") returned no error`)
	}
}

func TestAddFuncTimezoneWithoutScheduleReturnsError(t *testing.T) {
	c := New()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AddFunc panicked: %v", r)
		}
	}()
	if _, err := c.AddFunc("CRON_TZ=Asia/Tokyo", func() {}); err == nil {
		t.Fatal("AddFunc accepted a timezone with no schedule")
	}
	if n := len(c.Entries()); n != 0 {
		t.Fatalf("got %d entries, want 0", n)
	}
}

// The neighbouring cases that already returned errors or parsed fine must
// behave exactly as before.
func TestParseTimezoneBehaviourUnchanged(t *testing.T) {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Skip("no tzdata:", err)
	}
	ok := []struct {
		p    Parser
		spec string
		want Schedule
	}{
		{standardParser, "TZ=UTC 5 * * * *", every5min(time.UTC)},
		{standardParser, "CRON_TZ=Asia/Tokyo 5 * * * *", every5min(tokyo)},
		{secondParser, "TZ=Asia/Tokyo @midnight", midnight(tokyo)},
		{secondParser, "CRON_TZ=UTC  0 5 * * * *", every5min(time.UTC)},
	}
	for _, c := range ok {
		got, err := parseNoPanic(t, c.p, c.spec)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.spec, err)
			continue
		}
		if !equalSchedule(got, c.want) {
			t.Errorf("Parse(%q) = %v, want %v", c.spec, got, c.want)
		}
	}

	errs := []struct{ spec, err string }{
		{"TZ=UTC ", "expected exactly 5 fields, found 0"},
		{"TZ=Nowhere/Atlantis 5 * * * *", "provided bad location Nowhere/Atlantis"},
		{"TZ=UTC\t* * * * *", "provided bad location"},
	}
	for _, c := range errs {
		_, err := parseNoPanic(t, standardParser, c.spec)
		if err == nil || !strings.Contains(err.Error(), c.err) {
			t.Errorf("Parse(%q) error = %v, want it to contain %q", c.spec, err, c.err)
		}
	}
}

func equalSchedule(a, b Schedule) bool {
	sa, ok1 := a.(*SpecSchedule)
	sb, ok2 := b.(*SpecSchedule)
	if !ok1 || !ok2 {
		return false
	}
	return sa.Second == sb.Second && sa.Minute == sb.Minute && sa.Hour == sb.Hour &&
		sa.Dom == sb.Dom && sa.Month == sb.Month && sa.Dow == sb.Dow &&
		sa.Location.String() == sb.Location.String()
}
