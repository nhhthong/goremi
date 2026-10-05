// Tests for a track the player cannot play: the line under Search:, and playing again.
package app

import (
	"testing"

	"goremi/internal/player"
	"goremi/internal/provider"
)

// START: TestPlayerFailedNotice

func TestPlayerFailedNotice(t *testing.T) {
	pl := &recordingPlayer{events: make(chan player.Event, 1)}
	m, _ := playOne(New(&scriptedResolve{errs: []error{nil}}).WithPlayer(pl), provider.Track{ID: "a1", Title: "One"})
	cmd := m.Init()
	pl.events <- player.Failed
	next, _ := m.Update(cmd())
	if got, want := noticeLine(next.(Model)), `Cannot play "One".`; got != want {
		t.Fatalf("line under Search: = %q, want %q", got, want)
	}
}

// END: TestPlayerFailedNotice

// START: TestPlayRecoversAfterPlayerFailed

func TestPlayRecoversAfterPlayerFailed(t *testing.T) {
	pl := &recordingPlayer{events: make(chan player.Event, 1)}
	m, _ := playOne(New(&scriptedResolve{errs: []error{nil, nil}}).WithPlayer(pl), provider.Track{ID: "a1", Title: "One"})
	cmd := m.Init()
	pl.events <- player.Failed
	next, _ := m.Update(cmd())
	before := len(pl.urls)
	m, _ = playOne(next.(Model), provider.Track{ID: "a2", Title: "Two"})
	if len(pl.urls) != before+1 {
		t.Fatalf("Play calls went from %d to %d, want one more", before, len(pl.urls))
	}
	if got := noticeLine(m); got != "" {
		t.Fatalf("line under Search: = %q, want it gone", got)
	}
}

// END: TestPlayRecoversAfterPlayerFailed
