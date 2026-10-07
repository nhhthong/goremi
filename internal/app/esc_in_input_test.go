// Tests that `Esc` in the search input does nothing (task 10.3.1).
package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// START: TestEscInInputKeepsQueryAndFocus

func TestEscInInputKeepsQueryAndFocus(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "daft")
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	got := next.(Model)
	if got.Query() != "daft" || got.Focus() != FocusInput {
		t.Fatalf("query %q, focus %v; want query daft and the focus on the input", got.Query(), got.Focus())
	}
}

// END: TestEscInInputKeepsQueryAndFocus

// START: TestEscInInputDoesNotQuit

func TestEscInInputDoesNotQuit(t *testing.T) {
	if _, cmd := New(fakeProvider{}).Update(tea.KeyPressMsg{Code: tea.KeyEscape}); quits(cmd) {
		t.Fatal("Esc in the input returned a quit command, want none")
	}
}

// END: TestEscInInputDoesNotQuit
