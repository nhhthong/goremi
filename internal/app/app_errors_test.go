// Tests for the message line under the search box when a search fails.
package app

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"goremi/internal/provider"
)

// START: helpers

const networkText = "Unable to search. Check your internet connection."

// afterSearch returns the lines of the view after a search that ran the given steps.
func afterSearch(steps ...step) []string {
	p := &scriptedProvider{steps: steps}
	m := New(p)
	for range steps {
		m = search(m, "daft")
	}
	return strings.Split(m.View().Content, "\n")
}

// ansiCode matches a colour or style code, so tests can read a styled line as plain text.
var ansiCode = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// plain removes the colour and style codes from a view line.
func plain(s string) string { return ansiCode.ReplaceAllString(s, "") }

// lineAboveSearch returns the view line right above the "Search:" line (the message line), or "" when there is none.
func lineAboveSearch(lines []string) string {
	for i, l := range lines {
		if strings.HasPrefix(plain(l), "❯") && i > 1 {
			return plain(lines[i-2])
		}
	}
	return ""
}

// END: helpers

// START: network error line

func TestNetworkErrorOneLineBelowSearch(t *testing.T) {
	if got := lineAboveSearch(afterSearch(step{err: provider.ErrNetwork})); got != networkText {
		t.Fatalf("line under Search: = %q, want %q", got, networkText)
	}
}

func TestNoMessageBeforeSearch(t *testing.T) {
	if got := lineAboveSearch(strings.Split(New(fakeProvider{}).View().Content, "\n")); got != "" {
		t.Fatalf("line under Search: = %q, want none", got)
	}
}

func TestMessageAboveKeptList(t *testing.T) {
	lines := afterSearch(step{tracks: []provider.Track{{Title: "A"}, {Title: "B"}}}, step{err: provider.ErrNetwork})
	var order []string
	for _, l := range lines {
		switch {
		case strings.HasPrefix(plain(l), "❯"):
			order = append(order, "search")
		case plain(l) == networkText:
			order = append(order, "message")
		case strings.Contains(l, "A") && strings.Contains(l, "▶"):
			order = append(order, "list")
		}
	}
	if want := []string{"list", "message", "search"}; !reflect.DeepEqual(order, want) {
		t.Fatalf("order = %v, want %v in %q", order, want, lines)
	}
}

func TestNetworkErrorRecovers(t *testing.T) {
	lines := afterSearch(step{err: provider.ErrNetwork}, step{tracks: []provider.Track{{Title: "A"}}})
	if got := strings.Join(lines, "\n"); strings.Contains(got, networkText) || !strings.Contains(got, "A") {
		t.Fatalf("view %q: want the message gone and track A shown", got)
	}
}

// END: network error line

// START: other failure line

func TestOtherFailureMessage(t *testing.T) {
	want := "Search failed. Press Enter to retry."
	if got := lineAboveSearch(afterSearch(step{err: errors.New("boom")})); got != want {
		t.Fatalf("line under Search: = %q, want %q", got, want)
	}
}

func TestOtherFailureNotNetworkText(t *testing.T) {
	if got := strings.Join(afterSearch(step{err: errors.New("boom")}), "\n"); strings.Contains(got, "Unable to search.") {
		t.Fatalf("view %q must not show the network text", got)
	}
}

// END: other failure line
