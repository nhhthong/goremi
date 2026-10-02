// Tests for the end of the results, the loading line and duplicates when more tracks load.
package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// numbered returns n tracks titled prefix0, prefix1, … with the same text as ID.
func numbered(prefix string, n int) []provider.Track {
	var out []provider.Track
	for i := 0; i < n; i++ {
		out = append(out, provider.Track{ID: fmt.Sprintf("%s%d", prefix, i), Title: fmt.Sprintf("%s%d", prefix, i)})
	}
	return out
}

// tenTracks is a full page: the tracks A to J, without IDs.
func tenTracks() []provider.Track {
	var out []provider.Track
	for _, c := range "ABCDEFGHIJ" {
		out = append(out, provider.Track{Title: string(c)})
	}
	return out
}

// lastLine is the last line of the rendered list.
func lastLine(r Results) string {
	lines := strings.Split(r.Render(), "\n")
	return lines[len(lines)-1]
}

// START: TestNoLoadMoreUnderTenFirstPage

func TestNoLoadMoreUnderTenFirstPage(t *testing.T) {
	if out := NewResults(&fakeProvider{}, "q", numbered("a", 9)).Render(); strings.Contains(out, "load more...") {
		t.Fatalf("list %q has a load more line after a first page of 9 tracks", out)
	}
}

// END: TestNoLoadMoreUnderTenFirstPage

// START: TestLoadMoreAtTen

func TestLoadMoreAtTen(t *testing.T) {
	if got := lastLine(NewResults(&fakeProvider{}, "q", numbered("a", 10))); !strings.Contains(got, "load more...") {
		t.Fatalf("last line = %q, want load more... after a full page of 10", got)
	}
}

// END: TestLoadMoreAtTen

// START: TestDownStopsAtLastTrackWithoutMore

func TestDownStopsAtLastTrackWithoutMore(t *testing.T) {
	if got := pressed(numbered("a", 9), 8, tea.KeyDown); got != 8 {
		t.Fatalf("Selected() = %d, want 8 (the last track, there is no load more line)", got)
	}
}

// END: TestDownStopsAtLastTrackWithoutMore

// loadedFrom returns the list after Enter on its load more line, with the provider answering page 2 with page2.
func loadedFrom(page2 []provider.Track) Results {
	f := &fakeProvider{pages: map[int][]provider.Track{2: page2}}
	r, _ := pressEnter(NewResults(f, "q", numbered("a", 10)).Select(10))
	return r
}

// START: TestEndOfResultsAfterLoad

func TestEndOfResultsAfterLoad(t *testing.T) {
	r := loadedFrom(numbered("b", 3))
	if len(r.Tracks()) != 13 || strings.Contains(r.Render(), "load more...") {
		t.Fatalf("%d tracks, list %q; want 13 tracks and no load more line after a page of 3", len(r.Tracks()), r.Render())
	}
}

// END: TestEndOfResultsAfterLoad

// START: TestDownStopsAfterEnd

func TestDownStopsAfterEnd(t *testing.T) {
	r := loadedFrom(numbered("b", 3)).Select(12)
	r, _ = r.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if r.Selected() != 12 {
		t.Fatalf("Selected() = %d, want 12 (the last track)", r.Selected())
	}
}

// END: TestDownStopsAfterEnd

// START: TestEmptyPageClampsSelection

func TestEmptyPageClampsSelection(t *testing.T) {
	r := loadedFrom(nil)
	if r.Selected() != 9 || strings.Contains(r.Render(), "load more...") {
		t.Fatalf("Selected() = %d, list %q; want the last track (9) and no load more line", r.Selected(), r.Render())
	}
}

// END: TestEmptyPageClampsSelection

// START: TestLoadingLine

func TestLoadingLine(t *testing.T) {
	f := &fakeProvider{pages: map[int][]provider.Track{2: numbered("b", 10)}}
	r, _ := NewResults(f, "q", numbered("a", 10)).Select(10).Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if got := lastLine(r); !strings.Contains(got, "Loading...") || strings.Contains(got, "load more...") {
		t.Fatalf("last line = %q, want Loading... and not load more...", got)
	}
}

// END: TestLoadingLine

// START: TestLoadMoreLineBackAfterFullPage

func TestLoadMoreLineBackAfterFullPage(t *testing.T) {
	if got := lastLine(loadedFrom(numbered("b", 10))); !strings.Contains(got, "load more...") {
		t.Fatalf("last line = %q, want load more... again after a full page", got)
	}
}

// END: TestLoadMoreLineBackAfterFullPage

// duplicatePage is a page of 10 whose first two tracks are already in the list.
func duplicatePage() []provider.Track {
	return append([]provider.Track{{ID: "a8", Title: "a8"}, {ID: "a9", Title: "a9"}}, numbered("b", 8)...)
}

// START: TestDuplicatesListedOnce

func TestDuplicatesListedOnce(t *testing.T) {
	r := loadedFrom(duplicatePage())
	seen := map[string]bool{}
	for _, tr := range r.Tracks() {
		if seen[tr.ID] {
			t.Errorf("track %q is listed twice", tr.ID)
		}
		seen[tr.ID] = true
	}
	if len(r.Tracks()) != 18 {
		t.Fatalf("%d tracks, want 18 (10 + 8 new)", len(r.Tracks()))
	}
}

// END: TestDuplicatesListedOnce

// START: TestMoreCountsReturnedPage

func TestMoreCountsReturnedPage(t *testing.T) {
	if got := lastLine(loadedFrom(duplicatePage())); !strings.Contains(got, "load more...") {
		t.Fatalf("last line = %q, want load more... (the page returned 10 tracks, 2 of them duplicates)", got)
	}
}

// END: TestMoreCountsReturnedPage

// START: TestNoIDNotDeduplicated

func TestNoIDNotDeduplicated(t *testing.T) {
	f := &fakeProvider{pages: map[int][]provider.Track{2: {{Title: "K"}, {Title: "L"}}}}
	r, _ := pressEnter(NewResults(f, "q", tenTracks()).Select(10))
	if len(r.Tracks()) != 12 {
		t.Fatalf("%d tracks, want 12 (tracks without an ID are never taken for duplicates)", len(r.Tracks()))
	}
}

// END: TestNoIDNotDeduplicated
