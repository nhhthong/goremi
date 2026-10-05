// Tests for a file mpv cannot play: the player says so and keeps working.
package player

import (
	"testing"
	"time"
)

// START: TestUnplayableReportsFailed

func TestUnplayableReportsFailed(t *testing.T) {
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer p.Close()
	_ = p.Play("/no/such/file.wav")
	select {
	case e := <-p.Events():
		if e != Failed {
			t.Fatalf("event = %v, want Failed", e)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no Failed event within 3 s")
	}
}

// END: TestUnplayableReportsFailed

// START: TestPlayRecoversAfterFailed

func TestPlayRecoversAfterFailed(t *testing.T) {
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer p.Close()
	_ = p.Play("/no/such/file.wav")
	select {
	case <-p.Events():
	case <-time.After(3 * time.Second):
		t.Fatal("no Failed event within 3 s")
	}
	if err := p.Play(writeWAV(t, 4)); err != nil {
		t.Fatalf("Play() after Failed error: %v", err)
	}
	waitPosition(t, p)
}

// END: TestPlayRecoversAfterFailed
