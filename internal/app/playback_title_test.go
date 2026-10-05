// Tests for the title row of the player panel.
package app

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// panelView sends PlayMsg for each track in turn and returns the plain view.
func panelView(tracks ...provider.Track) string {
	var m tea.Model = sized(New(&resolveProvider{}).WithPlayer(&recordingPlayer{}), 100)
	for _, tr := range tracks {
		m, _ = m.Update(PlayMsg{Track: tr})
	}
	return plain(m.View().Content)
}

// START: TestPanelShowsTitle

func TestPanelShowsTitle(t *testing.T) {
	if view := panelView(provider.Track{ID: "a1", Title: "Alpha"}); !strings.Contains(view, "Alpha") {
		t.Fatalf("view has no title Alpha:\n%s", view)
	}
}

// END: TestPanelShowsTitle

// START: TestPanelTitleFollowsTrack

func TestPanelTitleFollowsTrack(t *testing.T) {
	view := panelView(provider.Track{ID: "a1", Title: "Alpha"}, provider.Track{ID: "a2", Title: "Bravo"})
	if !strings.Contains(view, "Bravo") || strings.Contains(view, "Alpha") {
		t.Fatalf("view should show Bravo and no Alpha:\n%s", view)
	}
}

// END: TestPanelTitleFollowsTrack
