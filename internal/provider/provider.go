// Provider contract and the Track model shared by every music source.
package provider

import "time"

// START: Track model

// Track is one song as the UI sees it (spec §12).
type Track struct {
	ID       string
	Title    string
	Artist   string
	Album    string
	Duration time.Duration
	Artwork  string
	Source   string
}

// END: Track model

// START: Provider interface

// Provider is a music source: Search a page, fill Details after selection, Resolve a stream URL.
type Provider interface {
	Search(query string, page int) ([]Track, error)
	Details(track Track) (Track, error)
	Resolve(track Track) (string, error)
}

// END: Provider interface
