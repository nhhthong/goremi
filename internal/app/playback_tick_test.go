// Tests for the tick that reads the player and moves the bar and the clock.
package app

import (
	"strings"

	"testing"
	"time"

	"goremi/internal/provider"
)

// tickedModel plays a track of 200 s with a fake player that says 151 s of 248 s, then runs one tick and feeds the reading back.
func tickedModel(t *testing.T) Model {
	t.Helper()
	pl := &recordingPlayer{position: 151 * time.Second, duration: 248 * time.Second}
	m := sized(New(&scriptedResolve{errs: []error{nil}, artist: "Daft Punk"}).WithPlayer(pl), 100)
	m, cmd := playOne(m, provider.Track{ID: "a1", Title: "One", Duration: 200 * time.Second})
	next, _ := m.Update(cmd())
	next, read := next.Update(tickMsg{})
	next, _ = next.Update(read())
	return next.(Model)
}

// START: TestTickShowsClock

func TestTickShowsClock(t *testing.T) {
	m := tickedModel(t)
	if lineWith(plainLines(m), "2:31 / 4:08") < 0 {
		t.Fatalf("no clock 2:31 / 4:08 in:\n%s", plain(m.View().Content))
	}
}

// END: TestTickShowsClock

// START: TestTickMovesBar

func TestTickMovesBar(t *testing.T) {
	m := tickedModel(t)
	if want := strings.Repeat("━", 23) + "●" + strings.Repeat("─", 14); lineWith(plainLines(m), want) < 0 {
		t.Fatalf("no bar %q in:\n%s", want, plain(m.View().Content))
	}
}

// END: TestTickMovesBar

// slowPosition is a fake player whose Position waits until released.
type slowPosition struct {
	recordingPlayer
	release chan struct{}
}

func (s *slowPosition) Position() time.Duration {
	<-s.release
	return 0
}

// START: TestTickReadsOffTheScreenLoop

func TestTickReadsOffTheScreenLoop(t *testing.T) {
	pl := &slowPosition{release: make(chan struct{})}
	m := New(fakeProvider{}).WithPlayer(pl)
	start := time.Now()
	_, cmd := m.Update(tickMsg{})
	elapsed := time.Since(start)
	close(pl.release)
	cmd()
	if elapsed > 100*time.Millisecond {
		t.Fatalf("Update took %v while Position blocked, want under 100ms", elapsed)
	}
}

// END: TestTickReadsOffTheScreenLoop
