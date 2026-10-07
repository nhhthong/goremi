// Tests for the theme colours the view uses.
package app

import (
	"strings"
	"testing"

	"goremi/internal/ui/theme"
)

// wantCodes fails when text lacks one of the colour codes.
func wantCodes(t *testing.T, text string, codes ...string) {
	t.Helper()
	for _, c := range codes {
		if !strings.Contains(text, c) {
			t.Errorf("%q lacks the colour code %s", text, c)
		}
	}
}

// START: TestSearchLabelAccentDefault

func TestSearchLabelAccentDefault(t *testing.T) {
	wantCodes(t, viewLines(New(fakeProvider{}))[1], "38;2;45;212;191")
}

// END: TestSearchLabelAccentDefault

// START: TestSearchLabelAccentDracula

func TestSearchLabelAccentDracula(t *testing.T) {
	wantCodes(t, viewLines(New(fakeProvider{}).WithTheme(theme.Dracula()))[1], "38;2;189;147;249")
}

// END: TestSearchLabelAccentDracula
