// Layout of the search view: track list on the left, player panel on the right.
package ui

import "charm.land/lipgloss/v2"

// START: JoinPanes

// JoinPanes places the list on the left and the pane on the right, top-aligned, two spaces apart.
func JoinPanes(list, pane string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, list, "  ", pane)
}

// END: JoinPanes
