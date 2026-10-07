// Tests for the tick loop: a tick comes every second while a track plays, once, and not before the first play or after a failed one.
package app

import (
	"os/exec"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// timed runs a command and returns its message and how long it took.
func timed(cmd tea.Cmd) (tea.Msg, time.Duration) {
	start := time.Now()
	msg := cmd()
	return msg, time.Since(start)
}

// checkTick fails unless cmd gives a tickMsg after 0.9 s to 1.5 s.
func checkTick(t *testing.T, cmd tea.Cmd) {
	t.Helper()
	if cmd == nil {
		t.Fatal("no tick command")
	}
	msg, d := timed(cmd) // a tea.Tick command works once: its timer is made when the command is
	if _, ok := msg.(tickMsg); ok && d >= 900*time.Millisecond && d <= 1500*time.Millisecond {
		return
	}
	if batch, ok := msg.(tea.BatchMsg); ok { // the note tick of the mascot comes in the same batch
		for _, c := range batch {
			if m, d := timed(c); true {
				if _, ok := m.(tickMsg); ok && d >= 900*time.Millisecond && d <= 1500*time.Millisecond {
					return
				}
			}
		}
	}
	t.Fatal("no command gives a tickMsg after about 1 s")
}

// playedOnce plays one track and returns the model and the command its details give.
func playedOnce(t *testing.T) (Model, tea.Cmd) {
	t.Helper()
	m, details := playOne(New(&scriptedResolve{errs: []error{nil, nil}, artist: "Daft Punk"}).WithPlayer(&recordingPlayer{}), provider.Track{ID: "a1", Title: "One"})
	next, tick := m.Update(details())
	_ = tick
	return next.(Model), tick
}

// START: TestTickStartsAfterPlay

func TestTickStartsAfterPlay(t *testing.T) {
	m, details := playOne(New(&scriptedResolve{errs: []error{nil}, artist: "Daft Punk"}).WithPlayer(&recordingPlayer{}), provider.Track{ID: "a1", Title: "One"})
	_, tick := m.Update(details())
	checkTick(t, tick)
}

// END: TestTickStartsAfterPlay

// START: TestTickRepeats

func TestTickRepeats(t *testing.T) {
	m, _ := playedOnce(t)
	_, again := m.Update(progressMsg{elapsed: time.Second})
	checkTick(t, again)
}

// END: TestTickRepeats

// START: TestTickDoesNotDouble

func TestTickDoesNotDouble(t *testing.T) {
	m, _ := playedOnce(t)
	m, details := playOne(m, provider.Track{ID: "a2", Title: "Two"})
	if _, second := m.Update(details()); second != nil {
		t.Fatal("the second play started a second tick loop")
	}
}

// END: TestTickDoesNotDouble

// START: TestNoTickAfterFailedPlay

func TestNoTickAfterFailedPlay(t *testing.T) {
	m, _ := playWith(exec.ErrNotFound)
	if _, cmd := m.Update(progressMsg{}); cmd != nil {
		t.Fatal("a tick was asked for after a failed play")
	}
}

// END: TestNoTickAfterFailedPlay

// START: TestNoTickBeforePlay

func TestNoTickBeforePlay(t *testing.T) {
	if _, cmd := New(fakeProvider{}).WithPlayer(&recordingPlayer{}).Update(progressMsg{}); cmd != nil {
		t.Fatal("a tick was asked for before any track played")
	}
}

// END: TestNoTickBeforePlay
