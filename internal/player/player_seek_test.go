// Tests for seeking back: a step of 10 seconds, and the stop at the start of the file.
package player

import (
	"testing"
	"time"
)

// START: TestSeekBack10

func TestSeekBack10(t *testing.T) {
	p := startPlaying(t, 30)
	if err := p.TogglePause(); err != nil {
		t.Fatalf("TogglePause() error: %v", err)
	}
	settle(t, p)
	for i := 0; i < 2; i++ {
		if err := p.Seek(10); err != nil {
			t.Fatalf("Seek(10) error: %v", err)
		}
	}
	time.Sleep(100 * time.Millisecond)
	before := p.Position()
	if err := p.Seek(-10); err != nil {
		t.Fatalf("Seek(-10) error: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	after := p.Position()
	if diff := (before - after).Round(time.Millisecond); diff < 9700*time.Millisecond || diff >= 10500*time.Millisecond { // null audio output reports up to 0.3s early
		t.Fatalf("Position() %v then %v after Seek(-10), want a step back of 9.7s to 10.5s", before, after)
	}
}

// END: TestSeekBack10

// START: TestSeekBackStopsAtZero

func TestSeekBackStopsAtZero(t *testing.T) {
	p := startPlaying(t, 30)
	if err := p.TogglePause(); err != nil {
		t.Fatalf("TogglePause() error: %v", err)
	}
	settle(t, p)
	if err := p.Seek(-10); err != nil {
		t.Fatalf("Seek(-10) error: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	if pos := p.Position(); pos < -300*time.Millisecond || pos > 50*time.Millisecond { // null audio output reports up to 0.3s early
		t.Fatalf("Position() = %v after Seek(-10) near the start, want 0 (-0.3s to 50ms)", pos)
	}
}

// END: TestSeekBackStopsAtZero
