// Layout of the search view: track list on the left, artwork pane on the right.
package ui

import "charm.land/lipgloss/v2"

// START: JoinPanes

// JoinPanes places the list on the left and the pane on the right, top-aligned, two spaces apart.
func JoinPanes(list, pane string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, list, "  ", pane)
}

// END: JoinPanes

// START: PlayerPanel

// PlayerPanel is the right-hand panel; until a track has played it shows the Goremi art.
func PlayerPanel() string { return Logo }

// END: PlayerPanel
