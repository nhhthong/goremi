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
}

// NewResults starts at page 1 with the tracks already loaded.
func NewResults(p provider.Provider, query string, tracks []provider.Track) Results {
	return Results{provider: p, query: query, tracks: tracks, page: 1}
}

// Select returns a copy with the given line selected; len(tracks) is the "load more..." line.
func (r Results) Select(i int) Results {
	r.selected = i
	return r
}

func (r Results) Tracks() []provider.Track { return r.tracks }

// END: Results

// START: Update

// moreLoadedMsg carries the result of loading the next page.
type moreLoadedMsg struct {
	page   int
	tracks []provider.Track
	err    error
}

// Update loads the next page on Enter over "load more..." and appends it when it arrives.
func (r Results) Update(msg tea.Msg) (Results, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEnter && len(r.tracks) > 0 && r.selected == len(r.tracks) {
			p, q, page := r.provider, r.query, r.page+1
			return r, func() tea.Msg {
				tracks, err := p.Search(q, page)
				return moreLoadedMsg{page: page, tracks: tracks, err: err}
			}
		}
	case moreLoadedMsg:
		if msg.err == nil {
			r.tracks = append(append([]provider.Track{}, r.tracks...), msg.tracks...)
			r.page = msg.page
		}
	}
	return r, nil
}

// END: Update
