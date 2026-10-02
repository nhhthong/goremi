// Tests for the theme colours of the list rows.
package app

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// rowLine returns the view line whose plain text is exactly want.
func rowLine(t *testing.T, m Model, want string) string {
	t.Helper()
	for _, l := range viewLines(m) {
		if plain(l) == want {
			return l
		}
	}
	t.Fatalf("no view line %q in %q", want, viewLines(m))
	return ""
}

// START: TestSelectedRowColours

func TestSelectedRowColours(t *testing.T) {
	wantCodes(t, rowLine(t, abModel(79), "▶ A"), "48;2;49;50;68", "38;2;137;180;250")
}

// END: TestSelectedRowColours

// START: TestSelectedRowBold

func TestSelectedRowBold(t *testing.T) {
	m := abModel(79)
	if sel := rowLine(t, m, "▶ A"); !strings.Contains(sel, "\x1b[1;") {
		t.Errorf("selected line %q is not bold", sel)
	}
	if other := rowLine(t, m, "  B"); strings.Contains(other, "\x1b[1;") {
		t.Errorf("unselected line %q is bold", other)
	}
}

// END: TestSelectedRowBold

// START: TestLoadMoreSelectedColours

func TestLoadMoreSelectedColours(t *testing.T) {
	var m tea.Model = abModel(79)
	m, _ = m.Update(down)
	m, _ = m.Update(down)
	wantCodes(t, rowLine(t, m.(Model), "▶ load more..."), "48;2;49;50;68")
}

// END: TestLoadMoreSelectedColours

// START: TestUnselectedRowForeground

func TestUnselectedRowForeground(t *testing.T) {
	line := rowLine(t, abModel(79), "  B")
	wantCodes(t, line, "38;2;205;214;244")
	if strings.Contains(line, "48;2;") {
		t.Errorf("unselected line %q has a background", line)
	}
}

// END: TestUnselectedRowForeground

// START: TestLoadMoreUnselectedMuted

func TestLoadMoreUnselectedMuted(t *testing.T) {
	wantCodes(t, rowLine(t, abModel(79), "  load more..."), "38;2;108;112;134")
}

// END: TestLoadMoreUnselectedMuted
