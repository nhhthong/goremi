// Content of the artwork pane (right side of the search view) for the selected track.
package ui

import "goremi/internal/provider"

// START: ArtworkPane

// ArtworkPane returns the pane content for the selected track; empty when there is nothing to show.
// Rendering a real image arrives with the artwork rows, so for now every case is empty.
func ArtworkPane(tracks []provider.Track, selected int) string {
	return ""
}

// END: ArtworkPane
