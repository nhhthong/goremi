// Tests of the v0.1.1 layout: the search bar at the bottom, the header before the first play, the message line, the hints and the playback hint (ui tasks 3.5.3 to 3.7.1, 10.14, 10.15, 19.12, 5.15.1).
package app

import (
	tea "charm.land/bubbletea/v2"
	"strings"
	"testing"
	"unicode/utf8"

	"goremi/internal/provider"
	"goremi/internal/ui/theme"
)

// START: helpers

// listAt is a model of the given size holding the tracks A, B and C found by one search; the list has the focus and nothing has played.
func listAt(width, height int) Model {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{trackA, trackB, trackC}}}}
	return search(sizeTo(New(p).WithPlayer(&recordingPlayer{}), width, height), "daft")
}

// playedList is listAt after a play request for B: the list has the focus and the panel shows the track.
func playedList(width, height int) Model {
	next, _ := listAt(width, height).Update(PlayMsg{Track: trackB})
	return next.(Model)
}

// searchIndex is the index of the top rule of the search bar, the line right above its `❯` line, or -1.
func searchIndex(lines []string) int {
	for i, l := range lines {
		if strings.HasPrefix(l, "❯") {
			return i - 1
		}
	}
	return -1
}

// runeCol is the column of the first rune of want in line, or -1.
func runeCol(line, want string) int {
	i := strings.Index(line, want)
	if i < 0 {
		return -1
	}
	return utf8.RuneCountInString(line[:i])
}

// END: helpers

// START: TestBottomBarIsLastLine

func TestBottomBarIsLastLine(t *testing.T) {
	lines := plainLines(playedList(100, 30))
	if len(lines) != 30 || !strings.HasPrefix(lines[28], "❯") {
		t.Fatalf("%d lines, prompt line %q; want 30 lines ending with the search bar", len(lines), lines[len(lines)-2])
	}
}

// END: TestBottomBarIsLastLine

// START: TestListAndPanelAboveBarSideBySide

func TestListAndPanelAboveBarSideBySide(t *testing.T) {
	lines := plainLines(playedList(100, 30))
	for _, l := range lines[:29] {
		if strings.Contains(l, "▶ A") && strings.Contains(l, "Loading…") {
			return
		}
	}
	t.Fatalf("no line above the bar holds the list row and the panel:\n%s", strings.Join(lines, "\n"))
}

// END: TestListAndPanelAboveBarSideBySide

// START: TestBarStaysLastWithoutHeight

func TestBarStaysLastWithoutHeight(t *testing.T) {
	lines := plainLines(playedList(100, 0))
	if !strings.HasPrefix(lines[len(lines)-2], "❯") {
		t.Fatalf("second last line %q, want the ❯ line", lines[len(lines)-2])
	}
	for i, l := range lines[:len(lines)-3] {
		if strings.TrimSpace(l) == "" {
			t.Fatalf("line %d is blank: with no height nothing pads the view", i)
		}
	}
}

// END: TestBarStaysLastWithoutHeight

// START: TestNarrowOrderPanelListBar

func TestNarrowOrderPanelListBar(t *testing.T) {
	lines := plainLines(playedList(60, 40))
	p, l, s := lineWith(lines, "Loading…"), lineWith(lines, "▶ A"), searchIndex(lines)
	if p < 0 || !(p < l && l < s) || s != len(lines)-3 {
		t.Fatalf("panel at %d, list at %d, bar at %d of %d lines; want panel, list, bar last", p, l, s, len(lines))
	}
}

// END: TestNarrowOrderPanelListBar

// START: header

func TestVersionDefaultsToDev(t *testing.T) {
	if Version != "dev" {
		t.Fatalf("Version = %q, want dev when unstamped", Version)
	}
}

func TestBadgeShowsVersion(t *testing.T) {
	old := Version
	Version = "9.9.9"
	defer func() { Version = old }()
	if got := plain(sizeTo(New(fakeProvider{}), 100, 30).View().Content); !strings.Contains(got, "Goremi v9.9.9") {
		t.Fatalf("view has no `Goremi v9.9.9`:\n%s", got)
	}
}

func TestHeaderAboveFullWidthList(t *testing.T) {
	lines := plainLines(listAt(100, 30))
	h, l := lineWith(lines, "Goremi v"), lineWith(lines, "▶ A")
	if h < 0 || l < 0 || h >= l {
		t.Fatalf("header at %d, list at %d; want the header above the list", h, l)
	}
	if strings.Contains(lines[l], "Loading…") || strings.Contains(strings.Join(lines, "\n"), "j -10s") {
		t.Fatalf("a panel shows before the first play:\n%s", strings.Join(lines, "\n"))
	}
}

func TestHeaderAloneWithoutList(t *testing.T) {
	got := plain(sizeTo(New(fakeProvider{}), 100, 30).View().Content)
	if !strings.Contains(got, "Current provider: YouTube") || strings.Contains(got, "▶") {
		t.Fatalf("view:\n%s\nwant the header and no list", got)
	}
}

func TestLogoArtNotDrawn(t *testing.T) {
	const art = "▄▀▀▀ ▄▀▀▄ █▀▀▄" // the first row of the old GOREMI art
	for name, m := range map[string]Model{"empty": sizeTo(New(fakeProvider{}), 100, 30), "list": listAt(100, 30), "narrow": sizeTo(New(fakeProvider{}), 60, 30)} {
		if got := plain(m.View().Content); strings.Contains(got, art) {
			t.Errorf("%s view still draws the GOREMI art:\n%s", name, got)
		}
	}
}

// END: header

// START: Esc hint and Down

func TestDownOnLoadMoreJumpsToBar(t *testing.T) {
	var tracks []provider.Track
	for i := 0; i < provider.PageSize; i++ {
		tracks = append(tracks, provider.Track{ID: string(rune('a' + i)), Title: string(rune('A' + i))})
	}
	m := search(sizeTo(New(&scriptedProvider{steps: []step{{tracks: tracks}}}), 100, 30), "daft")
	for i := 0; i < provider.PageSize; i++ {
		next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		m = next.(Model)
		if m.Focus() != FocusList {
			t.Fatalf("focus left the list after %d downs, before the load more line", i+1)
		}
	}
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if got := next.(Model).Focus(); got != FocusInput {
		t.Fatalf("focus after Down on the load more line = %v, want the search bar", got)
	}
}

func TestDownOnLastTrackStaysWithoutLoadMore(t *testing.T) {
	m := listAt(100, 30)
	for i := 0; i < 5; i++ {
		next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		m = next.(Model)
	}
	if m.Focus() != FocusList || m.Selected() != 2 {
		t.Fatalf("focus %v, selected %d; want the list and the last track (2)", m.Focus(), m.Selected())
	}
}

const wantEscHint = "Press Esc to return to search"

func TestEscHintAboveBar(t *testing.T) {
	lines := plainLines(listAt(100, 30))
	if s := searchIndex(lines); s < 1 || lines[s-1] != wantEscHint {
		t.Fatalf("line above the bar = %q, want %q", lines[searchIndex(lines)-1], wantEscHint)
	}
}

func TestNoEscHintInInput(t *testing.T) {
	if got := plain(listAt(100, 30).WithFocus(FocusInput).View().Content); strings.Contains(got, wantEscHint) {
		t.Fatalf("view holds the hint while the input has the focus:\n%s", got)
	}
}

func TestEscHintIsFaintMuted(t *testing.T) {
	lines := viewLines(listAt(100, 30))
	if got := lines[len(lines)-4]; !strings.Contains(got, "\x1b[2;38;2;138;138;138m") {
		t.Fatalf("hint line %q lacks the faint code with Muted 138;138;138", got)
	}
}

func TestEscHintAboveMessage(t *testing.T) {
	m := listAt(100, 30)
	m.notice = "boom"
	lines := plainLines(m)
	s := searchIndex(lines)
	if lines[s-1] != "boom" || lines[s-2] != wantEscHint {
		t.Fatalf("lines above the bar = %q and %q, want the hint then the message", lines[s-2], lines[s-1])
	}
}

// END: Esc hint and Down

// START: message line

func failedAt(err error, th theme.Theme) Model {
	return sizeTo(failedWith(err).WithTheme(th), 100, 30)
}

func TestMessageAboveBar(t *testing.T) {
	lines := plainLines(failedAt(provider.ErrNetwork, theme.Default()))
	if s := searchIndex(lines); s < 1 || lines[s-1] != "Unable to search. Check your internet connection." {
		t.Fatalf("line above the bar = %q", lines[searchIndex(lines)-1])
	}
}

func TestMessageIsRedDefault(t *testing.T) {
	lines := viewLines(failedAt(provider.ErrNetwork, theme.Default()))
	if got := lines[len(lines)-4]; !strings.Contains(got, "38;2;255;85;85") {
		t.Fatalf("message line %q lacks the Error colour 255;85;85", got)
	}
}

func TestMessageIsRedNord(t *testing.T) {
	lines := viewLines(failedAt(provider.ErrNetwork, theme.Nord()))
	if got := lines[len(lines)-4]; !strings.Contains(got, "38;2;191;97;106") {
		t.Fatalf("message line %q lacks the Nord Error colour 191;97;106", got)
	}
}

func TestNoResultsLineIsRedAboveBar(t *testing.T) {
	m := sizeTo(search(New(&scriptedProvider{steps: []step{{}}}), "daft"), 100, 30)
	lines := viewLines(m)
	got := lines[len(lines)-4]
	if plain(got) != `No results for "daft".` || !strings.Contains(got, "38;2;255;85;85") {
		t.Fatalf("line above the bar = %q, want the red no-results message", got)
	}
}

// END: message line

// START: playback hint

const keysHint = "j -10s  k pause  l +10s  p prev  n next"

func keysHintBelowControls(t *testing.T, width int) {
	t.Helper()
	lines := plainLines(playedList(width, 40))
	c := lineWith(lines, "◀◀")
	if c < 0 || c+1 >= len(lines) || !strings.Contains(lines[c+1], keysHint) {
		t.Fatalf("width %d: the line below the controls is not the keys hint:\n%s", width, strings.Join(lines, "\n"))
	}
	if runeCol(lines[c], "◀◀") != runeCol(lines[c+1], "j -10s") {
		t.Fatalf("width %d: the hint starts at column %d, the controls at %d", width, runeCol(lines[c+1], "j -10s"), runeCol(lines[c], "◀◀"))
	}
}

func TestKeysHintBelowControls(t *testing.T)       { keysHintBelowControls(t, 100) }
func TestKeysHintBelowControlsNarrow(t *testing.T) { keysHintBelowControls(t, 60) }

func TestKeysHintAbsentInInput(t *testing.T) {
	if got := plain(playedList(100, 30).WithFocus(FocusInput).View().Content); strings.Contains(got, keysHint) {
		t.Fatalf("the keys hint shows while the input has the focus:\n%s", got)
	}
}

// END: playback hint
