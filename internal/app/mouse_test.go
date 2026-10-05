// Tests that the [ui] mouse setting reaches the view.
package app

import (
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// START: TestMouseOnCellMotion

func TestMouseOnCellMotion(t *testing.T) {
	if got := New(fakeProvider{}).WithMouse(true).View().MouseMode; got != tea.MouseModeCellMotion {
		t.Fatalf("MouseMode = %v, want MouseModeCellMotion", got)
	}
}

// END: TestMouseOnCellMotion

// START: TestMouseOffNone

func TestMouseOffNone(t *testing.T) {
	if got := New(fakeProvider{}).WithMouse(false).View().MouseMode; got != tea.MouseModeNone {
		t.Fatalf("MouseMode = %v, want MouseModeNone", got)
	}
}

// END: TestMouseOffNone

// START: TestConfigMouseReachesView

func TestConfigMouseReachesView(t *testing.T) {
	off := NewFromConfig(writeConfig(t, "[ui]\nmouse = false\n"), fakeProvider{}).View().MouseMode
	none := NewFromConfig(filepath.Join(t.TempDir(), "missing.toml"), fakeProvider{}).View().MouseMode
	if off != tea.MouseModeNone || none != tea.MouseModeCellMotion {
		t.Fatalf("MouseMode with mouse = false: %v, with no file: %v, want MouseModeNone and MouseModeCellMotion", off, none)
	}
}

// END: TestConfigMouseReachesView
