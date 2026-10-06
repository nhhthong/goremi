// Tests for the controls of a playing file: length, pause, seek and the end of the file.
package player

import (
	"testing"
	"time"
)

// START: settle

// settle returns the position of a paused player once mpv has caught up: after a pause mpv plays out its audio buffer and time-pos jumps to the real stop point late.
func settle(t *testing.T, p *Player) time.Duration {
	t.Helper()
	time.Sleep(500 * time.Millisecond)
	return p.Position()
}

// END: settle

// START: TestDurationIsFileLength

func TestDurationIsFileLength(t *testing.T) {
	p := startPlaying(t, 4)
	d := p.Duration()
	if diff := d - 4*time.Second; diff < -50*time.Millisecond || diff > 50*time.Millisecond {
		t.Fatalf("Duration() = %v, want 4s within 50ms", d)
	}
}

// END: TestDurationIsFileLength

// START: TestTogglePausePauses

func TestTogglePausePauses(t *testing.T) {
	p := startPlaying(t, 30)
	if err := p.TogglePause(); err != nil {
		t.Fatalf("TogglePause() error: %v", err)
	}
	before := settle(t, p)
	time.Sleep(300 * time.Millisecond)
	after := p.Position()
	if !p.Paused() || after != before {
		t.Fatalf("Paused() = %v, Position() %v then %v, want paused and unchanged", p.Paused(), before, after)
	}
}

// END: TestTogglePausePauses

// START: TestTogglePauseResumes

func TestTogglePauseResumes(t *testing.T) {
	p := startPlaying(t, 30)
	for i := 0; i < 2; i++ {
		if err := p.TogglePause(); err != nil {
			t.Fatalf("TogglePause() #%d error: %v", i+1, err)
		}
	}
	before := p.Position()
	time.Sleep(300 * time.Millisecond)
	if after := p.Position(); p.Paused() || after <= before {
		t.Fatalf("Paused() = %v, Position() %v then %v, want playing and growing", p.Paused(), before, after)
	}
}

// END: TestTogglePauseResumes

// START: TestSeekForward10

func TestSeekForward10(t *testing.T) {
	p := startPlaying(t, 30)
	if err := p.TogglePause(); err != nil { // paused, the position holds still between the two reads
		t.Fatalf("TogglePause() error: %v", err)
	}
	before := settle(t, p)
	if err := p.Seek(10); err != nil {
		t.Fatalf("Seek(10) error: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	after := p.Position()
	if diff := (after - before).Round(time.Millisecond); diff < 9700*time.Millisecond || diff >= 10500*time.Millisecond { // null audio output reports up to 0.3s early
		t.Fatalf("Position() %v then %v after Seek(10), want a step of 9.7s to 10.5s", before, after)
	}
}

// END: TestSeekForward10

// START: TestEndOfFileReportsEnded

func TestEndOfFileReportsEnded(t *testing.T) {
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer p.Close()
	if err := p.Play(writeWAV(t, 1)); err != nil {
		t.Fatalf("Play() error: %v", err)
	}
	select {
	case e := <-p.Events():
		if e != Ended {
			t.Fatalf("event = %v, want Ended", e)
		}
	case <-time.After(6 * time.Second):
		t.Fatal("no Ended event within 6 s")
	}
	select {
	case e := <-p.Events():
		t.Fatalf("event %v after Ended, want none", e)
	case <-time.After(300 * time.Millisecond):
	}
}

// END: TestEndOfFileReportsEnded
