// Tests for the clock row and the controls row of the player panel.
package ui

import (
	"regexp"
	"testing"
	"time"

	"goremi/internal/ui/theme"
)

// ansiCode matches the colour and style codes a rendered line carries.
var ansiCode = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// checkClock fails when ClockText(elapsed, total) is not want.
func checkClock(t *testing.T, elapsed, total int, want string) {
	t.Helper()
	if got := ClockText(time.Duration(elapsed)*time.Second, time.Duration(total)*time.Second); got != want {
		t.Fatalf("ClockText(%ds, %ds) = %q, want %q", elapsed, total, got, want)
	}
}

// START: TestClockMinutes

func TestClockMinutes(t *testing.T) { checkClock(t, 151, 248, "2:31 / 4:08") }

// END: TestClockMinutes

// START: TestClockSeconds

func TestClockSeconds(t *testing.T) { checkClock(t, 5, 30, "0:05 / 0:30") }

// END: TestClockSeconds

// START: TestClockHours

func TestClockHours(t *testing.T) { checkClock(t, 3723, 3800, "1:02:03 / 1:03:20") }

// END: TestClockHours

// START: TestClockHourBoundary

func TestClockHourBoundary(t *testing.T) { checkClock(t, 3599, 3600, "59:59 / 1:00:00") }

// END: TestClockHourBoundary

// START: TestControlsLineGlyphs

func TestControlsLineGlyphs(t *testing.T) {
	got := ansiCode.ReplaceAllString(ControlsLine(theme.Default(), false), "")
	if want := "◀◀ ◀ ❚❚ ▶ ▶▶"; got != want {
		t.Fatalf("ControlsLine = %q, want %q", got, want)
	}
}

// END: TestControlsLineGlyphs
