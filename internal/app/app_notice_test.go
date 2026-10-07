// Tests for the message line and the screen colours: plain text keeps the terminal's colours.
package app

import (
	"testing"

	"goremi/internal/provider"
	"goremi/internal/ui/theme"
)

// failedSearchLine returns the raw message line under Search: after a network failure, drawn with t.
func failedSearchLine(t *testing.T, th theme.Theme) string {
	t.Helper()
	p := &scriptedProvider{steps: []step{{err: provider.ErrNetwork}}}
	return rowLine(t, search(New(p).WithTheme(th), "daft"), networkText)
}

// START: TestViewNoTerminalColoursDefault

func TestViewNoTerminalColoursDefault(t *testing.T) {
	if v := New(fakeProvider{}).View(); v.BackgroundColor != nil || v.ForegroundColor != nil {
		t.Fatalf("background %v, foreground %v; want neither set", v.BackgroundColor, v.ForegroundColor)
	}
}

// END: TestViewNoTerminalColoursDefault

// START: TestViewNoTerminalColoursDracula

func TestViewNoTerminalColoursDracula(t *testing.T) {
	if v := New(fakeProvider{}).WithTheme(theme.Dracula()).View(); v.BackgroundColor != nil || v.ForegroundColor != nil {
		t.Fatalf("background %v, foreground %v; want neither set", v.BackgroundColor, v.ForegroundColor)
	}
}

// END: TestViewNoTerminalColoursDracula
