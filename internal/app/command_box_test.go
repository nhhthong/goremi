// Tests of the command box: `/` in the search bar lists /quit and /theme above the bar (ui tasks 3.1.1 to 3.1.7).
package app

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// START: helpers

// slash types text into an empty search bar of a 100x30 model.
func slash(text string) Model {
	m, _ := typed(sizeTo(New(&scriptedProvider{}), 100, 30), text)
	return m
}

// press sends one key to the model and returns the model and the command.
func pressKeyMsg(m Model, k tea.KeyPressMsg) (Model, tea.Cmd) {
	next, cmd := m.Update(k)
	return next.(Model), cmd
}

// END: helpers

// START: TestSlashOpensCommandList

func TestSlashOpensCommandList(t *testing.T) {
	lines := plainLines(slash("/"))
	s := searchIndex(lines)
	if s < 2 || !strings.Contains(lines[s-2], "/quit") || !strings.Contains(lines[s-1], "/theme") {
		t.Fatalf("the two lines above the bar = %q, want /quit then /theme", lines[max(s-2, 0):s])
	}
}

// END: TestSlashOpensCommandList

// START: TestNoListWithoutSlash

func TestNoListWithoutSlash(t *testing.T) {
	if got := plain(slash("a").View().Content); strings.Contains(got, "/quit") || strings.Contains(got, "/theme") {
		t.Fatalf("a list shows with no slash:\n%s", got)
	}
}

// END: TestNoListWithoutSlash

// START: TestListFiltersByPrefix

func TestListFiltersByPrefix(t *testing.T) {
	got := plain(slash("/q").View().Content)
	if !strings.Contains(got, "/quit") || strings.Contains(got, "/theme") {
		t.Fatalf("view after /q:\n%s\nwant /quit only", got)
	}
}

// END: TestListFiltersByPrefix

// START: TestListClosesWithoutMatch

func TestListClosesWithoutMatch(t *testing.T) {
	got := plain(slash("/x").View().Content)
	if strings.Contains(got, "/quit") || strings.Contains(got, "/theme") {
		t.Fatalf("view after /x:\n%s\nwant no list", got)
	}
}

// END: TestListClosesWithoutMatch

// START: TestArrowsMoveHighlight

func TestArrowsMoveHighlight(t *testing.T) {
	m, _ := pressKeyMsg(slash("/"), tea.KeyPressMsg{Code: tea.KeyDown})
	if got := plain(m.View().Content); !strings.Contains(got, "▶ /theme") || strings.Contains(got, "▶ /quit") {
		t.Fatalf("after Down:\n%s\nwant /theme highlighted", got)
	}
	m, _ = pressKeyMsg(m, tea.KeyPressMsg{Code: tea.KeyUp})
	if got := plain(m.View().Content); !strings.Contains(got, "▶ /quit") {
		t.Fatalf("after Up:\n%s\nwant /quit highlighted", got)
	}
}

// END: TestArrowsMoveHighlight

// START: TestTabCompletesHighlighted

func TestTabCompletesHighlighted(t *testing.T) {
	m, _ := pressKeyMsg(slash("/"), tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = pressKeyMsg(m, tea.KeyPressMsg{Code: tea.KeyTab})
	if m.Query() != "/theme" || m.Focus() != FocusInput {
		t.Fatalf("Query() = %q, focus %v; want /theme and the focus still on the bar", m.Query(), m.Focus())
	}
}

// END: TestTabCompletesHighlighted

// START: TestEnterOnQuitQuits

func TestEnterOnQuitQuits(t *testing.T) {
	if _, cmd := pressKeyMsg(slash("/q"), tea.KeyPressMsg{Code: tea.KeyEnter}); !quits(cmd) {
		t.Fatal("Enter on /quit returned no quit command")
	}
}

// END: TestEnterOnQuitQuits

// START: TestEnterOnThemeSendsOpenTheme

func TestEnterOnThemeSendsOpenTheme(t *testing.T) {
	m, cmd := pressKeyMsg(slash("/t"), tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil || cmd() != (OpenThemeMsg{}) || m.Query() != "" {
		t.Fatalf("command %v, Query() = %q; want OpenThemeMsg and an empty bar", cmd, m.Query())
	}
}

// END: TestEnterOnThemeSendsOpenTheme

// START: TestEscClosesListKeepsText

func TestEscClosesListKeepsText(t *testing.T) {
	m, _ := pressKeyMsg(slash("/q"), tea.KeyPressMsg{Code: tea.KeyEscape})
	if got := plain(m.View().Content); strings.Contains(got, "▶ /quit") || m.Query() != "/q" {
		t.Fatalf("Query() = %q, view:\n%s\nwant the text kept and the list closed", m.Query(), got)
	}
}

// END: TestEscClosesListKeepsText

// START: TestUnknownCommandMessage

func TestUnknownCommandMessage(t *testing.T) {
	p := &scriptedProvider{}
	m, _ := typed(sizeTo(New(p), 100, 30), "/xyz")
	m, cmd := pressKeyMsg(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	lines := plainLines(m)
	if s := searchIndex(lines); s < 1 || lines[s-1] != `Unknown command "/xyz".` || p.next != 0 || cmd != nil {
		t.Fatalf("message %q, Search calls %d, command %v; want the unknown-command line and no search", lines[searchIndex(lines)-1], p.next, cmd)
	}
}

// END: TestUnknownCommandMessage
