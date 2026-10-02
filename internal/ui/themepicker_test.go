// Tests for the theme picker list.
package ui

import (
	"strings"
	"testing"

	"goremi/internal/ui/theme"
)

// START: TestPickerListsThemes

func TestPickerListsThemes(t *testing.T) {
	lines := strings.Split(NewThemePicker(theme.All(), "dark").View(), "\n")
	last := -1
	for _, name := range []string{"light", "dark", "cyberpunk"} {
		found := -1
		for i, l := range lines {
			if strings.Contains(l, name) {
				found = i
			}
		}
		if found <= last {
			t.Fatalf("%q is on line %d, want a line after %d; view %q", name, found, last, lines)
		}
		last = found
	}
}

// END: TestPickerListsThemes

// START: TestPickerReadsGivenList

func TestPickerReadsGivenList(t *testing.T) {
	view := NewThemePicker([]theme.Entry{{Name: "alpha"}, {Name: "beta"}}, "alpha").View()
	for _, name := range []string{"alpha", "beta"} {
		if !strings.Contains(view, name) {
			t.Errorf("view %q lacks %q", view, name)
		}
	}
	if strings.Contains(view, "light") {
		t.Errorf("view %q lists a theme that was not given", view)
	}
}

// END: TestPickerReadsGivenList
