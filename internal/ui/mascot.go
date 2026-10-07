// The mascot of the header: the line-art gopher with a guitar, in braille characters, painted from the theme.
package ui

import (
	"strings"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// START: Mascot

// mascotRows is the line-art gopher in braille dots (2x4 per cell), made once from an image that is no longer in the repository, and embedded; no file is read at run time. The notes of the image are removed: the program generates its own (notes.go).
var mascotRows = []string{
	"⠀⠀⠀⠀⢠⠞⣉⡙⣆⡤⠔⠒⠋⠋⠙⠉⠒⠦⣴⢉⡉⢦⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠈⢦⣀⠟⠁⠀⢀⣀⡀⠀⠀⠀⠀⣀⣈⡻⣤⠞⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⡏⠀⠀⣴⣿⣿⣿⣧⠀⠀⣼⣿⣿⣿⣿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⢸⠀⠀⠀⣿⣿⡏⠘⣻⢀⣤⡿⣿⣅⣙⠏⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠸⡄⠀⠀⠈⠛⠛⠒⠁⢿⣶⡾⠎⠉⠁⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⢳⡀⠀⠀⠀⠀⠀⠀⠘⠛⠃⠀⠀⠀⣰⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⣸⠁⠀⠤⣄⣀⣀⢀⡴⠶⢦⡀⠀⣀⣹⣤⣤⣾⢿⠛⣹⠄⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⢘⢦⡀⠀⣸⣶⡌⣡⠶⢆⣶⣿⣿⣿⡿⢿⣿⡏⠻⠏⠏⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⣾⠀⠩⡏⠙⣿⡅⠳⣤⡜⠉⢸⠃⠀⣧⠤⠚⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⢿⠀⠀⠱⡀⠈⠁⣀⠔⠒⠚⠁⠀⠀⢸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠈⣳⣄⡀⠉⠉⠉⣁⣀⣀⡀⠀⢀⣠⣏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠤⠠⠤⠐⠐⣒⢰⣡⣀⡏⣉⣉⢉⠉⠉⢉⡉⣉⠉⣆⣤⡧⢐⣒⠐⠒⠀⠤⠤⠠⠄⠀",
}

// MascotRows and MascotWidth are the size of the mascot in rows and columns.
const (
	MascotRows  = 12
	MascotWidth = 32
)

// Mascot is the art as plain text, 12 rows of 32 braille characters.
func Mascot() string { return strings.Join(mascotRows, "\n") }

// END: Mascot

// START: PaintMascot

// PaintMascot draws the notes over the art and colours it column by column, from t.LogoFrom on the first column to t.LogoTo on the last; a note is drawn in t.Accent and blank cells stay plain.
func PaintMascot(t theme.Theme, notes []Note) string {
	colours := lipgloss.Blend1D(MascotWidth, lipgloss.Color(t.LogoFrom), lipgloss.Color(t.LogoTo))
	grid := make([][]rune, len(mascotRows))
	for i, row := range mascotRows {
		grid[i] = []rune(row)
	}
	isNote := map[[2]int]bool{}
	for _, n := range notes {
		for i, c := range noteCells {
			grid[n.Row+c[0]][n.Col+c[1]] = noteRunes[i]
			isNote[[2]int{n.Row + c[0], n.Col + c[1]}] = true
		}
	}
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent))
	rows := make([]string, len(grid))
	for y, row := range grid {
		var b strings.Builder
		for x, r := range row {
			switch {
			case isNote[[2]int{y, x}]:
				b.WriteString(accent.Render(string(r)))
			case r == 0x2800: // the blank braille cell
				b.WriteRune(r)
			default:
				b.WriteString(lipgloss.NewStyle().Foreground(colours[x]).Render(string(r)))
			}
		}
		rows[y] = b.String()
	}
	return strings.Join(rows, "\n")
}

// END: PaintMascot
