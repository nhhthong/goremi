// Tests for moving the selection in the theme picker and for its colours.
package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/ui/theme"
)

var (
	pickDown = tea.KeyPressMsg{Code: tea.KeyDown}
	pickUp   = tea.KeyPressMsg{Code: tea.KeyUp}
)

// START: TestPickerDown

func TestPickerDown(t *testing.T) {
	if got := NewThemePicker(theme.All(), "light").Update(pickDown).Selected(); got != "dark" {
		t.Fatalf("Selected() = %q, want dark", got)
	}
}

// END: TestPickerDown

// START: TestPickerUp

func TestPickerUp(t *testing.T) {
	if got := NewThemePicker(theme.All(), "dark").Update(pickUp).Selected(); got != "light" {
		t.Fatalf("Selected() = %q, want light", got)
	}
}

// END: TestPickerUp

// START: TestPickerStopsAtEnds

func TestPickerStopsAtEnds(t *testing.T) {
	if got := NewThemePicker(theme.All(), "light").Update(pickUp).Selected(); got != "light" {
		t.Errorf("↑ on the first line: Selected() = %q, want light", got)
	}
	if got := NewThemePicker(theme.All(), "cyberpunk").Update(pickDown).Selected(); got != "cyberpunk" {
		t.Errorf("↓ on the last line: Selected() = %q, want cyberpunk", got)
	}
}

// END: TestPickerStopsAtEnds

// START: TestPickerColoursSelectedTheme

func TestPickerColoursSelectedTheme(t *testing.T) {
	lines := strings.Split(NewThemePicker(theme.All(), "cyberpunk").View(), "\n")
	for _, c := range []string{"48;2;26;26;58", "38;2;0;240;255"} {
		if !strings.Contains(lines[2], c) {
			t.Errorf("selected line %q lacks %s", lines[2], c)
		}
	}
	if !strings.Contains(lines[0], "38;2;224;224;255") {
		t.Errorf("unselected line %q lacks 38;2;224;224;255", lines[0])
	}
}

// END: TestPickerColoursSelectedTheme

// START: TestPickerColoursFollowSelection

func TestPickerColoursFollowSelection(t *testing.T) {
	view := NewThemePicker(theme.All(), "dark").Update(pickDown).View()
	if !strings.Contains(view, "48;2;26;26;58") || strings.Contains(view, "48;2;49;50;68") {
		t.Fatalf("view %q: want the cyberpunk background and not the dark one", view)
	}
}

// END: TestPickerColoursFollowSelection
