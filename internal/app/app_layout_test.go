// Helpers for the layout tests: the model sizes and the lines of the view. The layout cases themselves are in v011_layout_test.go.
package app

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: helpers

// art0 is the first row of the old GOREMI art, the fragment that tells whether the art is on screen.
func art0() string { return "▄▀▀▀ ▄▀▀▄ █▀▀▄" }

// sized sends a WindowSizeMsg of the given width to the model.
func sized(m Model, width int) Model {
	next, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 40})
	return next.(Model)
}

// abModel returns a model of the given width whose list holds the tracks A and B.
func abModel(width int) Model {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}, {Title: "B"}}}}}
	return search(sized(New(p), width), "daft")
}

// isSearch tells whether a view line is the Search: line.
func isSearch(l string) bool { return strings.HasPrefix(plain(l), "❯") }

// END: helpers
