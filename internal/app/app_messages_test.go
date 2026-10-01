// Tests for the one-line messages under the search box: empty query, missing yt-dlp, timeout, no results and clearing.
package app

import (
	"os/exec"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: helpers

// viewLines returns the lines of the model's view.
func viewLines(m Model) []string { return strings.Split(m.View().Content, "\n") }

// emptyEnter presses Enter on the given model's input.
func emptyEnter(m Model) Model {
	next, _ := m.Update(enter)
	return next.(Model)
}

// failedWith returns a model after one search that failed with err.
func failedWith(err error) Model {
	return search(New(&scriptedProvider{steps: []step{{err: err}}}), "daft")
}

// END: helpers

// START: empty query message

const emptyText = "Type something to search."

func TestEmptyQueryMessage(t *testing.T) {
	if got := lineAfterSearch(viewLines(emptyEnter(New(fakeProvider{})))); got != emptyText {
		t.Fatalf("line under Search: = %q, want %q", got, emptyText)
	}
}

func TestSpacesQueryMessage(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "   ")
	if got := lineAfterSearch(viewLines(emptyEnter(m))); got != emptyText {
		t.Fatalf("line under Search: = %q, want %q", got, emptyText)
	}
}

func TestEmptyEnterKeepsList(t *testing.T) {
	m := abcList().WithFocus(FocusInput)
	for range "daft" {
		next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
		m = next.(Model)
	}
	if got, want := titles(emptyEnter(m).Tracks()), []string{"A", "B", "C"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("tracks = %v, want %v", got, want)
	}
}

// END: empty query message

// START: yt-dlp missing message

const missingText = "yt-dlp not found. Install yt-dlp and try again."

func TestYtDlpMissingMessage(t *testing.T) {
	m := failedWith(&exec.Error{Name: "yt-dlp", Err: exec.ErrNotFound})
	if got := lineAfterSearch(viewLines(m)); got != missingText {
		t.Fatalf("line under Search: = %q, want %q", got, missingText)
	}
}

func TestYtDlpMissingEndToEnd(t *testing.T) {
	t.Setenv("PATH", "")
	m := search(New(&provider.YouTubeProvider{}), "daft")
	if got := lineAfterSearch(viewLines(m)); got != missingText {
		t.Fatalf("line under Search: = %q, want %q", got, missingText)
	}
}

// END: yt-dlp missing message

// START: timeout message

func TestTimeoutMessage(t *testing.T) {
	want := "Search timed out. Press Enter to retry."
	if got := lineAfterSearch(viewLines(failedWith(provider.ErrTimeout))); got != want {
		t.Fatalf("line under Search: = %q, want %q", got, want)
	}
}

func TestEnterRetriesAfterTimeout(t *testing.T) {
	p := &failingOnce{err: provider.ErrTimeout}
	m := search(New(p), "daft").WithFocus(FocusInput) // 18.3.1 will return the focus to the input by itself
	_, cmd := m.Update(enter)
	if cmd == nil {
		t.Fatal("want a command for the retry")
	}
	cmd()
	if want := []searchCall{{"daft", 1}, {"daft", 1}}; !reflect.DeepEqual(p.calls, want) {
		t.Fatalf("calls = %v, want %v", p.calls, want)
	}
}

// failingOnce fails its first Search with err, records every call and answers later ones with no tracks.
type failingOnce struct {
	fakeProvider
	err   error
	calls []searchCall
}

func (f *failingOnce) Search(q string, page int) ([]provider.Track, error) {
	f.calls = append(f.calls, searchCall{q, page})
	if len(f.calls) == 1 {
		return nil, f.err
	}
	return nil, nil
}

// END: timeout message

// START: no results message

func TestNoResultsMessage(t *testing.T) {
	want := `No results for "daft".`
	if got := lineAfterSearch(viewLines(search(New(fakeProvider{}), "daft"))); got != want {
		t.Fatalf("line under Search: = %q, want %q", got, want)
	}
}

func TestNoResultsUsesTrimmedQuery(t *testing.T) {
	want := `No results for "daft".`
	if got := lineAfterSearch(viewLines(search(New(fakeProvider{}), "  daft  "))); got != want {
		t.Fatalf("line under Search: = %q, want %q", got, want)
	}
}

func TestNoNoResultsMessageWhenFound(t *testing.T) {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}, {Title: "B"}}}}}
	if got := strings.Join(viewLines(search(New(p), "daft")), "\n"); strings.Contains(got, "No results") {
		t.Fatalf("view %q must not say No results", got)
	}
}

// END: no results message

// START: clearing

func TestMessageClearsOnTyping(t *testing.T) {
	m, _ := typed(failedWith(provider.ErrNetwork).WithFocus(FocusInput), "x")
	if got := strings.Join(viewLines(m), "\n"); strings.Contains(got, networkText) {
		t.Fatalf("view %q still shows the message", got)
	}
}

func TestMessageClearsOnBackspace(t *testing.T) {
	next, _ := failedWith(provider.ErrNetwork).WithFocus(FocusInput).Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	if got := strings.Join(viewLines(next.(Model)), "\n"); strings.Contains(got, networkText) {
		t.Fatalf("view %q still shows the message", got)
	}
}

func TestMessageClearsOnEnter(t *testing.T) {
	next, _ := failedWith(provider.ErrNetwork).WithFocus(FocusInput).Update(enter)
	if got := strings.Join(viewLines(next.(Model)), "\n"); strings.Contains(got, networkText) {
		t.Fatalf("view %q still shows the message", got)
	}
}

func TestEmptyMessageClearsOnTyping(t *testing.T) {
	m, _ := typed(emptyEnter(New(fakeProvider{})), "x")
	if got := strings.Join(viewLines(m), "\n"); strings.Contains(got, emptyText) {
		t.Fatalf("view %q still shows the message", got)
	}
}

// END: clearing
