// Results model: the track list state and the Enter key on the "load more..." line.
package ui

import (
	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: Results

// Results holds the loaded tracks, the selected line and the page loaded so far.
type Results struct {
	provider provider.Provider
	query    string
	tracks   []provider.Track
	selected int
	page     int
	// more is true while the list ends with a "load more..." line: the last page returned a full page of tracks.
	more bool
	// loading is true from Enter on "load more..." until the next page arrives or fails.
	loading bool
	// loadErr is the error of the last failed load of the next page, until the app takes it.
	loadErr error
}

// NewResults starts at page 1 with the tracks already loaded.
func NewResults(p provider.Provider, query string, tracks []provider.Track) Results {
	return Results{provider: p, query: query, tracks: tracks, page: 1, more: len(tracks) >= provider.PageSize}
}

// Select returns a copy with the given line selected; len(tracks) is the "load more..." line.
func (r Results) Select(i int) Results {
	r.selected = i
	return r
}

func (r Results) Tracks() []provider.Track { return r.tracks }
func (r Results) Selected() int            { return r.selected }

// LoadErr is the error of the last failed load of the next page, nil when the last load worked or the error was taken.
func (r Results) LoadErr() error { return r.loadErr }

// ClearLoadErr returns a copy without the load error.
func (r Results) ClearLoadErr() Results {
	r.loadErr = nil
	return r
}

// lastLine is the last line the selection can reach: the "load more..." line while there is one, else the last track.
func (r Results) lastLine() int {
	if r.more {
		return len(r.tracks)
	}
	return len(r.tracks) - 1
}

// END: Results

// START: Update

// moreLoadedMsg carries the result of loading the next page.
type moreLoadedMsg struct {
	page   int
	tracks []provider.Track
	err    error
}

// Update moves the selection on ↑/↓ and loads the next page on Enter over "load more..." and appends it when it arrives.
func (r Results) Update(msg tea.Msg) (Results, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.Code == tea.KeyDown && len(r.tracks) > 0 && r.selected < r.lastLine() {
			r.selected++
		}
		if msg.Code == tea.KeyUp && r.selected > 0 {
			r.selected--
		}
		if msg.Code == tea.KeyEnter && r.more && len(r.tracks) > 0 && r.selected == len(r.tracks) {
			r.loading = true
			p, q, page := r.provider, r.query, r.page+1
			return r, func() tea.Msg {
				tracks, err := p.Search(q, page)
				return moreLoadedMsg{page: page, tracks: tracks, err: err}
			}
		}
	case moreLoadedMsg:
		r.loading = false
		r.loadErr = msg.err
		if msg.err == nil {
			r.tracks = append(append([]provider.Track{}, r.tracks...), r.unseen(msg.tracks)...)
			r.page = msg.page
			r.more = len(msg.tracks) >= provider.PageSize
			r.selected = min(r.selected, r.lastLine())
		}
	}
	return r, nil
}

// unseen returns the tracks whose ID is not in the list yet; a track without an ID is never taken for a duplicate.
func (r Results) unseen(tracks []provider.Track) []provider.Track {
	have := map[string]bool{}
	for _, t := range r.tracks {
		have[t.ID] = true
	}
	var out []provider.Track
	for _, t := range tracks {
		if t.ID != "" && have[t.ID] {
			continue
		}
		out = append(out, t)
	}
	return out
}

// END: Update
