// The header of the empty screen: the mascot with the badge beside it or under it.
package ui

import (
	"strings"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// START: Badge

// Badge is the three lines of the info badge; the version comes from the build.
func Badge(version string) []string {
	return []string{"Goremi v" + version, "Created by nhhthong", "Current provider: YouTube"}
}

// END: Badge

// START: Header

// WideMin is the width, in columns, from which the layout is wide: list and panel side by side, badge beside the mascot.
const WideMin = 80

// HeaderRows is the number of rows of the header at a width: the mascot rows, and the badge lines under them below WideMin.
func HeaderRows(width int) int {
	if width >= WideMin {
		return MascotRows
	}
	return MascotRows + len(Badge(""))
}

// headerGap is the number of columns between the mascot and the badge.
const headerGap = 2

// Header is the mascot, painted from the theme, with the badge: from WideMin columns the badge sits to the right, centred on the mascot rows, two columns apart; below it goes under the mascot.
func Header(t theme.Theme, width int, version string, notes []Note) string {
	badge := Badge(version)
	if width < WideMin {
		return PaintMascot(t, notes) + "\n" + strings.Join(badge, "\n")
	}
	top := strings.Repeat("\n", (MascotRows-len(badge))/2)
	return lipgloss.JoinHorizontal(lipgloss.Top, PaintMascot(t, notes), strings.Repeat(" ", headerGap), top+strings.Join(badge, "\n"))
}

// END: Header
