// Results list rendering: one line per track, a selection marker, and a "load more..." line.
package ui

import (
	"strings"

	"goremi/internal/provider"
)

// START: RenderResults

// RenderResults renders one line per track; the selected line starts with "▶ ".
func RenderResults(tracks []provider.Track, selected int) string {
	lines := make([]string, len(tracks))
	for i, t := range tracks {
		prefix := "  "
		if i == selected {
			prefix = "▶ "
		}
		lines[i] = prefix + t.Title
	}
	return strings.Join(lines, "\n")
}

// END: RenderResults

// START: RenderResultsMore

// RenderResultsMore renders the tracks followed by a selectable "load more..." line.
func RenderResultsMore(tracks []provider.Track, selected int) string {
	if len(tracks) == 0 {
		return ""
	}
	prefix := "  "
	if selected == len(tracks) {
		prefix = "▶ "
	}
	return RenderResults(tracks, selected) + "\n" + prefix + "load more..."
}

// END: RenderResultsMore
