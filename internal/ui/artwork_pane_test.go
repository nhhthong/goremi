// Tests for the artwork pane content that follows the selected track.
package ui

import (
	"testing"

	"goremi/internal/provider"
)

// START: TestArtworkPaneNoArtwork

func TestArtworkPaneNoArtwork(t *testing.T) {
	if got := ArtworkPane([]provider.Track{{Title: "Get Lucky"}}, 0); got != "" {
		t.Fatalf("want an empty pane, got %q", got)
	}
}

// END: TestArtworkPaneNoArtwork

// START: TestArtworkPaneSelectedWithoutArtwork

func TestArtworkPaneSelectedWithoutArtwork(t *testing.T) {
	tracks := []provider.Track{{Title: "Get Lucky"}, {Title: "One More Time", Artwork: "http://x/a.jpg"}}
	if got := ArtworkPane(tracks, 0); got != "" {
		t.Fatalf("want an empty pane, got %q", got)
	}
}

// END: TestArtworkPaneSelectedWithoutArtwork

// START: TestArtworkPaneNoSelection

func TestArtworkPaneNoSelection(t *testing.T) {
	tracks := []provider.Track{{Title: "Get Lucky", Artwork: "http://x/a.jpg"}}
	for _, selected := range []int{len(tracks), -1} {
		if got := ArtworkPane(tracks, selected); got != "" {
			t.Errorf("selected=%d: want an empty pane, got %q", selected, got)
		}
	}
}

// END: TestArtworkPaneNoSelection
