// Tests for the frame tick of the spectrum: 30 frames a second, the bars follow the bands while audio plays, fall idle when paused or ended, and stay empty when mpv has no filter.
package app

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/player"
	"goremi/internal/ui"
)

// playingSpectrum is a model with a track on (artist set), the spectrum on and this player.
func playingSpectrum(pl Player) Model {
	m := specApp(pl, true)
	m.playing, m.artist = specTrack, specTrack.Artist
	return m
}

// frame sends one frame to the model and returns it.
func frame(m Model) (Model, tea.Cmd) {
	next, cmd := m.Update(frameMsg{})
	return next.(Model), cmd
}

// every returns a bar-heights value with all 32 bars at h.
func every(h int) [ui.SpectrumBars]int {
	var out [ui.SpectrumBars]int
	for i := range out {
		out[i] = h
	}
	return out
}

// START: TestFrameIs30PerSecond

func TestFrameIs30PerSecond(t *testing.T) {
	if frameEvery*30 > time.Second || frameEvery*31 <= time.Second {
		t.Fatalf("frameEvery = %v, want one thirtieth of a second", frameEvery)
	}
}

// END: TestFrameIs30PerSecond

// START: TestFrameTickRunsWithSpectrumOn

func TestFrameTickRunsWithSpectrumOn(t *testing.T) {
	_, msgs := playAndSettle(specApp(&recordingPlayer{}, true), specTrack)
	n := 0
	for _, msg := range msgs {
		if _, ok := msg.(frameMsg); ok {
			n++
		}
	}
	if n == 0 {
		t.Fatal("no frame tick ran during a play with the spectrum on")
	}
}

// END: TestFrameTickRunsWithSpectrumOn

// START: TestNoFrameTickWithSpectrumOff

func TestNoFrameTickWithSpectrumOff(t *testing.T) {
	_, msgs := playAndSettle(specApp(&recordingPlayer{}, false), specTrack)
	for _, msg := range msgs {
		if _, ok := msg.(frameMsg); ok {
			t.Fatal("a frame tick ran with the spectrum off")
		}
	}
}

// END: TestNoFrameTickWithSpectrumOff

// START: TestFrameBarsFollowBands

func TestFrameBarsFollowBands(t *testing.T) {
	pl := newBandsPlayer(0)
	m, cmd := frame(playingSpectrum(pl))
	if m.bars != every(64) {
		t.Fatalf("bars after a frame of 0 dB = %v, want all 64 (a bar rises at once)", m.bars)
	}
	if cmd == nil {
		t.Fatal("the frame asked for no next frame while a track is on")
	}
	pl.db = -60
	m, _ = frame(m)
	if m.bars != every(54) {
		t.Fatalf("bars after a frame of -60 dB = %v, want all 54 (85 %% of 64)", m.bars)
	}
}

// END: TestFrameBarsFollowBands

// START: TestFrameEndsWithNoTrack

func TestFrameEndsWithNoTrack(t *testing.T) {
	m := playingSpectrum(newBandsPlayer(0))
	m.artist = ""
	if _, cmd := frame(m); cmd != nil {
		t.Fatal("the frame tick went on with no track on")
	}
}

// END: TestFrameEndsWithNoTrack

// START: TestFrameBarsStayEmptyWithoutFilter

func TestFrameBarsStayEmptyWithoutFilter(t *testing.T) {
	pl := newBandsPlayer(0)
	pl.on = false
	m, _ := frame(playingSpectrum(pl))
	if m.bars != every(0) {
		t.Fatalf("bars with no filter = %v, want all 0", m.bars)
	}
}

// END: TestFrameBarsStayEmptyWithoutFilter

// START: TestPausedBarsAreIdle

func TestPausedBarsAreIdle(t *testing.T) {
	m := playingSpectrum(newBandsPlayer(0))
	m.paused = true
	m.intn = func(n int) int { return n - 1 }
	m, _ = frame(m)
	if m.bars != every(3) {
		t.Fatalf("bars while paused = %v, want the idle height 3 (not the bands, which are at 64)", m.bars)
	}
}

// END: TestPausedBarsAreIdle

// START: TestEndedBarsAreIdle

func TestEndedBarsAreIdle(t *testing.T) {
	m := playingSpectrum(newBandsPlayer(0))
	m.intn = func(n int) int { return n - 1 }
	next, _ := m.Update(playerEventMsg{e: player.Ended}) // the only track of the list ended: no next
	m, _ = frame(next.(Model))
	if m.bars != every(3) {
		t.Fatalf("bars after the last track ended = %v, want the idle height 3", m.bars)
	}
}

// END: TestEndedBarsAreIdle

// START: TestIdleBarsFallByTheRule

func TestIdleBarsFallByTheRule(t *testing.T) {
	m := playingSpectrum(newBandsPlayer(0))
	m.paused = true
	m.bars = every(64)
	m.intn = func(int) int { return 0 }
	m, _ = frame(m)
	if m.bars != every(54) {
		t.Fatalf("bars after pausing from 64 = %v, want 54: the idle target goes through the fall rule", m.bars)
	}
}

// END: TestIdleBarsFallByTheRule

// START: TestBarsReturnToBandsWhenPlayingAgain

func TestBarsReturnToBandsWhenPlayingAgain(t *testing.T) {
	m := playingSpectrum(newBandsPlayer(0))
	m.paused = true
	m.intn = func(int) int { return 0 }
	m, _ = frame(m)
	m.paused = false
	m, _ = frame(m)
	if m.bars != every(64) {
		t.Fatalf("bars after playing again = %v, want 64 (the bands)", m.bars)
	}
}

// END: TestBarsReturnToBandsWhenPlayingAgain

// START: TestPlayAfterEndedFollowsBands

func TestPlayAfterEndedFollowsBands(t *testing.T) {
	m := playingSpectrum(newBandsPlayer(0))
	next, _ := m.Update(playerEventMsg{e: player.Ended})
	next, _ = next.Update(PlayMsg{Track: specTrack})
	m, _ = frame(next.(Model))
	if m.bars != every(64) {
		t.Fatalf("bars after a new play = %v, want 64 (the bands)", m.bars)
	}
}

// END: TestPlayAfterEndedFollowsBands
