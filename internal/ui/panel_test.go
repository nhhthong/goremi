// Tests for the player panel shown before any track has played.
package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// START: TestPanelIsMultiLineArt

func TestPanelIsMultiLineArt(t *testing.T) {
	p := PlayerPanel(theme.Default())
	if p == "" || len(strings.Split(p, "\n")) < 3 {
		t.Fatalf("PlayerPanel(theme.Default()) = %q, want at least 3 lines", p)
	}
}

// END: TestPanelIsMultiLineArt

// START: TestPanelFitsWidth

func TestPanelFitsWidth(t *testing.T) {
	for i, l := range strings.Split(PlayerPanel(theme.Default()), "\n") {
		if w := lipgloss.Width(l); w > 40 {
			t.Errorf("line %d is %d columns wide, want at most 40", i, w)
		}
	}
}

// END: TestPanelFitsWidth
