// Tests for a failed load of the next page: the message shows, the list stays, and a retry asks for the same page.
package app

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// pageRecorder answers each Search call with the next scripted step and records the page asked for.
type pageRecorder struct {
	fakeProvider
	steps []step
	next  int
	pages []int
}

func (p *pageRecorder) Search(_ string, page int) ([]provider.Track, error) {
	p.pages = append(p.pages, page)
	st := p.steps[p.next]
	p.next++
	return st.tracks, st.err
}

// onLoadMore returns a model whose list holds A to J with the load more line selected; later searches follow later.
func onLoadMore(later ...step) (Model, *pageRecorder) {
	p := &pageRecorder{steps: append([]step{{tracks: tenTracks()}}, later...)}
	var m tea.Model = search(New(p), "daft")
	for i := 0; i < 10; i++ {
		m, _ = m.Update(down)
	}
	return m.(Model), p
}

// loadMore presses Enter on the model and feeds the result of the command back.
func loadMore(m Model) Model {
	next, cmd := m.Update(enter)
	after, _ := next.Update(cmd())
	return after.(Model)
}

// START: TestLoadMoreFailureMessage

func TestLoadMoreFailureMessage(t *testing.T) {
	m, _ := onLoadMore(step{err: provider.ErrNetwork})
	if got := lineAboveSearch(viewLines(loadMore(m))); got != networkText {
		t.Fatalf("line under Search: = %q, want %q", got, networkText)
	}
}

// END: TestLoadMoreFailureMessage

// START: TestLoadMoreFailureKeepsList

func TestLoadMoreFailureKeepsList(t *testing.T) {
	m, _ := onLoadMore(step{err: provider.ErrNetwork})
	after := loadMore(m)
	if len(after.Tracks()) != 10 || !strings.Contains(plain(after.View().Content), "load more...") {
		t.Fatalf("%d tracks, view %q; want 10 tracks and the load more line", len(after.Tracks()), after.View().Content)
	}
}

// END: TestLoadMoreFailureKeepsList

// START: TestLoadMoreFailureRecovers

func TestLoadMoreFailureRecovers(t *testing.T) {
	m, p := onLoadMore(step{err: provider.ErrNetwork}, step{tracks: []provider.Track{{Title: "K"}, {Title: "L"}}})
	after := loadMore(loadMore(m))
	if got := p.pages; len(got) != 3 || got[1] != 2 || got[2] != 2 {
		t.Fatalf("pages asked for = %v, want 1, 2, 2", got)
	}
	if got := titles(after.Tracks()); len(got) != 12 || got[10] != "K" || got[11] != "L" {
		t.Fatalf("tracks = %v, want A to J then K, L", got)
	}
}

// END: TestLoadMoreFailureRecovers
