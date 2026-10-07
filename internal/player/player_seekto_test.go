// Tests of the absolute seek of a playing file (task 5.32).
package player

import (
	"testing"
	"time"
)

// START: TestSeekToMovesToPosition

func TestSeekToMovesToPosition(t *testing.T) {
	p := startPlaying(t, 30)
	if err := p.TogglePause(); err != nil {
		t.Fatalf("TogglePause() error: %v", err)
	}
	settle(t, p)
	if err := p.SeekTo(20 * time.Second); err != nil {
		t.Fatalf("SeekTo(20s) error: %v", err)
	}
	time.Sleep(200 * time.Millisecond)
	if got := p.Position(); got < 19500*time.Millisecond || got > 21500*time.Millisecond {
		t.Fatalf("Position() = %v after SeekTo(20s), want 19.5s to 21.5s", got)
	}
}

// END: TestSeekToMovesToPosition
