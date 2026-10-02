// Tests for the full screen: the views ask for the alternate screen.
package app

import (
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: TestViewAltScreen

func TestViewAltScreen(t *testing.T) {
	if !New(fakeProvider{}).View().AltScreen {
		t.Fatal("the app view does not ask for the alternate screen")
	}
}

// END: TestViewAltScreen

// START: TestViewAltScreenAfterSearch

func TestViewAltScreenAfterSearch(t *testing.T) {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}}}}}
	next, _ := New(p).Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	if !search(next.(Model), "daft").View().AltScreen {
		t.Fatal("the app view after a search does not ask for the alternate screen")
	}
}

// END: TestViewAltScreenAfterSearch

// START: TestThemeViewAltScreen

func TestThemeViewAltScreen(t *testing.T) {
	if !NewThemeModel(filepath.Join(t.TempDir(), "config.toml")).View().AltScreen {
		t.Fatal("the theme picker view does not ask for the alternate screen")
	}
}

// END: TestThemeViewAltScreen
