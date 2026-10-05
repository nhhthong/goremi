// Tests for a track that ends: the next track of the results list plays, and at the last one playback stops.
package app

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/player"
	"goremi/internal/provider"
)

// collect runs a command and returns the messages it gives within 100 ms; a command that waits is left running, and a batch is run command by command.
func collect(cmd tea.Cmd) []tea.Msg {
	out := make(chan tea.Msg, 8)
	var run func(c tea.Cmd)
	run = func(c tea.Cmd) {
		if c == nil {
			return
		}
		go func() {
			msg := c()
			if batch, ok := msg.(tea.BatchMsg); ok {
				for _, b := range batch {
					run(b)
				}
				return
			}
			if msg != nil {
				out <- msg
			}
		}()
	}
	run(cmd)
	var msgs []tea.Msg
	for timeout := time.After(100 * time.Millisecond); ; {
		select {
		case msg := <-out:
			msgs = append(msgs, msg)
		case <-timeout:
			return msgs
		}
	}
}

// endedMessages plays the track in a list of A, B, C, lets the player report Ended and returns the messages the app's commands give.
func endedMessages(playing provider.Track) []tea.Msg {
	pl := &recordingPlayer{events: make(chan player.Event, 1)}
	m := playingModelWith(pl, playing)
	cmd := m.Init()
	pl.events <- player.Ended
	_, next := m.Update(cmd())
	return collect(next)
}

// hasPlayMsg tells whether msgs holds a PlayMsg for the track.
func hasPlayMsg(msgs []tea.Msg, track provider.Track) bool {
	for _, msg := range msgs {
		if p, ok := msg.(PlayMsg); ok && p.Track == track {
			return true
		}
	}
	return false
}

// START: TestEndedPlaysNext

func TestEndedPlaysNext(t *testing.T) {
	if msgs := endedMessages(trackA); !hasPlayMsg(msgs, trackB) {
		t.Fatalf("messages after Ended on A = %v, want a PlayMsg for B", msgs)
	}
}

// END: TestEndedPlaysNext

// START: TestEndedAtLastStops

func TestEndedAtLastStops(t *testing.T) {
	for _, msg := range endedMessages(trackC) {
		if _, ok := msg.(PlayMsg); ok {
			t.Fatalf("PlayMsg %v after Ended on the last track, want none", msg)
		}
	}
}

// END: TestEndedAtLastStops
