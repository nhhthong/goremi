// Tests for playing a track: the app resolves it, then hands the URL to the player.
package app

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/player"
	"goremi/internal/provider"
)

// resolveProvider records the tracks it resolves and answers with a fixed URL.
type resolveProvider struct {
	fakeProvider
	resolved []provider.Track
}

func (r *resolveProvider) Resolve(t provider.Track) (string, error) {
	r.resolved = append(r.resolved, t)
	return "http://stream/1", nil
}

// recordingPlayer records the URLs it is asked to play.
type recordingPlayer struct {
	urls    []string
	toggles int
	seeks   []float64
	// seekTos are the positions SeekTo was asked for.
	seekTos []time.Duration
	events  chan player.Event
	// playErrs are the errors Play returns, one per call in turn; nil or none means success.
	playErrs []error
	// position and duration are what Position and Duration report.
	position, duration time.Duration
	// pausedNow is what Paused reports.
	pausedNow bool
}

func (r *recordingPlayer) Play(url string) error {
	r.urls = append(r.urls, url)
	if len(r.playErrs) == 0 {
		return nil
	}
	err := r.playErrs[0]
	r.playErrs = r.playErrs[1:]
	return err
}

func (r *recordingPlayer) Events() <-chan player.Event { return r.events }
func (r *recordingPlayer) Position() time.Duration     { return r.position }
func (r *recordingPlayer) Duration() time.Duration     { return r.duration }
func (r *recordingPlayer) Paused() bool                { return r.pausedNow }

func (r *recordingPlayer) TogglePause() error {
	r.toggles++
	return nil
}

func (r *recordingPlayer) SeekTo(position time.Duration) error {
	r.seekTos = append(r.seekTos, position)
	return nil
}

func (r *recordingPlayer) Seek(seconds float64) error {
	r.seeks = append(r.seeks, seconds)
	return nil
}

// START: TestPlayMsgResolvesTrack

func TestPlayMsgResolvesTrack(t *testing.T) {
	p, track := &resolveProvider{}, provider.Track{ID: "a1", Title: "One"}
	_, cmd := New(p).WithPlayer(&recordingPlayer{}).Update(PlayMsg{Track: track})
	if cmd == nil {
		t.Fatal("PlayMsg returned no command")
	}
	cmd()
	if len(p.resolved) != 1 || p.resolved[0] != track {
		t.Fatalf("Resolve calls = %v, want one with %v", p.resolved, track)
	}
}

// END: TestPlayMsgResolvesTrack

// START: TestResolvedURLIsPlayed

func TestResolvedURLIsPlayed(t *testing.T) {
	pl := &recordingPlayer{}
	m := New(&resolveProvider{}).WithPlayer(pl)
	_, cmd := m.Update(PlayMsg{Track: provider.Track{ID: "a1"}})
	var msg tea.Msg = cmd()
	_, play := m.Update(msg)
	play() // Play runs in the command Update returns
	if len(pl.urls) != 1 || pl.urls[0] != "http://stream/1" {
		t.Fatalf("Play calls = %v, want one with http://stream/1", pl.urls)
	}
}

// END: TestResolvedURLIsPlayed
