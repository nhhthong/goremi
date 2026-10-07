// Tests of the three-line search bar and of the header that stays while a track plays (ui tasks 3.5.5, 3.5.6, 3.7.2).
package app

import (
	"strings"
	"testing"
)

// rule is a line of n `─`.
func rule(n int) string { return strings.Repeat("─", n) }

// START: TestBarIsThreeLines

func TestBarIsThreeLines(t *testing.T) {
	lines := plainLines(playedList(100, 30))
	if len(lines) != 30 || lines[27] != rule(100) || !strings.HasPrefix(lines[28], "❯ daft") || lines[29] != rule(100) {
		t.Fatalf("%d lines; last three: %q %q %q; want a rule, `❯ daft` and a rule", len(lines), lines[len(lines)-3], lines[len(lines)-2], lines[len(lines)-1])
	}
}

// END: TestBarIsThreeLines

// START: TestRuleIsViewWide

func TestRuleIsViewWide(t *testing.T) {
	for _, width := range []int{100, 60} {
		lines := plainLines(playedList(width, 30))
		if lines[len(lines)-1] != rule(width) || lines[len(lines)-3] != rule(width) {
			t.Errorf("width %d: rules %q and %q, want %d cells of ─", width, lines[len(lines)-3], lines[len(lines)-1], width)
		}
	}
}

// END: TestRuleIsViewWide

// START: TestNoSearchLabel

func TestNoSearchLabel(t *testing.T) {
	if got := plain(playedList(100, 30).View().Content); strings.Contains(got, "Search:") {
		t.Fatalf("the label Search: is still drawn:\n%s", got)
	}
}

// END: TestNoSearchLabel

// START: TestPromptAccentWithInputFocus

func TestPromptAccentWithInputFocus(t *testing.T) {
	lines := viewLines(listAt(100, 30).WithFocus(FocusInput))
	if got := lines[len(lines)-2]; !strings.Contains(got, "38;2;45;212;191") {
		t.Fatalf("prompt line %q lacks the Accent 45;212;191", got)
	}
}

// END: TestPromptAccentWithInputFocus

// START: TestPromptMutedWithListFocus

func TestPromptMutedWithListFocus(t *testing.T) {
	lines := viewLines(listAt(100, 30))
	if got := lines[len(lines)-2]; !strings.Contains(got, "\x1b[2;38;2;138;138;138m") {
		t.Fatalf("prompt line %q lacks the faint Muted 138;138;138", got)
	}
}

// END: TestPromptMutedWithListFocus

// START: TestRulesUseBorder

func TestRulesUseBorder(t *testing.T) {
	lines := viewLines(listAt(100, 30))
	for _, i := range []int{len(lines) - 3, len(lines) - 1} {
		if !strings.Contains(lines[i], "38;2;107;114;128") {
			t.Errorf("rule line %q lacks the Border 107;114;128", lines[i])
		}
	}
}

// END: TestRulesUseBorder

// START: TestHeaderWhilePlaying

func TestHeaderWhilePlaying(t *testing.T) {
	got := plain(playedList(100, 30).View().Content)
	if !strings.Contains(got, "Created by nhhthong") || !strings.Contains(got, "▶ A") || !strings.Contains(got, "Loading…") {
		t.Fatalf("view while a track plays:\n%s\nwant the badge, the list and the panel", got)
	}
}

// END: TestHeaderWhilePlaying

// START: TestHeaderAboveListWhilePlaying

func TestHeaderAboveListWhilePlaying(t *testing.T) {
	for _, width := range []int{100, 60} {
		lines := plainLines(playedList(width, 40))
		h, p, l := lineWith(lines, "Goremi v"), lineWith(lines, "Loading…"), lineWith(lines, "▶ A")
		if h < 0 || p < 0 || l < 0 || h >= p || h >= l {
			t.Errorf("width %d: header at %d, panel at %d, list at %d; want the header above both", width, h, p, l)
		}
	}
}

// END: TestHeaderAboveListWhilePlaying
