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
	if got := NewThemePicker(theme.All(), "default").Update(pickDown).Selected(); got != "catppuccin" {
		t.Fatalf("Selected() = %q, want catppuccin", got)
	}
}

// END: TestPickerDown

// START: TestPickerUp

func TestPickerUp(t *testing.T) {
	if got := NewThemePicker(theme.All(), "catppuccin").Update(pickUp).Selected(); got != "default" {
		t.Fatalf("Selected() = %q, want default", got)
	}
}

// END: TestPickerUp

// START: TestPickerStopsAtEnds

func TestPickerStopsAtEnds(t *testing.T) {
	if got := NewThemePicker(theme.All(), "default").Update(pickUp).Selected(); got != "default" {
		t.Errorf("↑ on the first line: Selected() = %q, want default", got)
	}
	if got := NewThemePicker(theme.All(), "tokyonight").Update(pickDown).Selected(); got != "tokyonight" {
		t.Errorf("↓ on the last line: Selected() = %q, want tokyonight", got)
	}
}

// END: TestPickerStopsAtEnds

// START: TestPickerColoursSelectedTheme

func TestPickerColoursSelectedTheme(t *testing.T) {
	lines := strings.Split(NewThemePicker(theme.All(), "dracula").View(), "\n")
	for _, c := range []string{"48;2;68;71;90", "38;2;189;147;249"} {
		if !strings.Contains(lines[2], c) {
			t.Errorf("selected line %q lacks %s", lines[2], c)
		}
	}
	if lines[0] != "default" {
		t.Errorf("unselected line %q carries a colour code, want the plain name", lines[0])
	}
}

// END: TestPickerColoursSelectedTheme

// START: TestPickerColoursFollowSelection

func TestPickerColoursFollowSelection(t *testing.T) {
	view := NewThemePicker(theme.All(), "default").Update(pickDown).View()
	if !strings.Contains(view, "48;2;49;50;68") || strings.Contains(view, "48;2;19;78;74") {
		t.Fatalf("view %q: want the catppuccin background and not the default one", view)
	}
}

// END: TestPickerColoursFollowSelection
