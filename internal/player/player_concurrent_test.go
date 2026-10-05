// Tests for using one player from several goroutines while mpv events arrive.
package player

import (
	"sync"
	"testing"
	"time"
)

// START: TestConcurrentToggleKeepsParity

func TestConcurrentToggleKeepsParity(t *testing.T) {
	p := startPlaying(t, 30)
	start := make(chan struct{})
	var wg sync.WaitGroup
	errs := make(chan error, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs <- p.TogglePause()
		}()
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("TogglePause() error: %v", err)
		}
	}
	if p.Paused() {
		t.Fatal("Paused() = true after 20 toggles, want false: a toggle was lost")
	}
}

// END: TestConcurrentToggleKeepsParity

// START: TestReadWhileEventsUpdate

func TestReadWhileEventsUpdate(t *testing.T) {
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer p.Close()
	if err := p.Play(writeWAV(t, 1)); err != nil {
		t.Fatalf("Play() error: %v", err)
	}
	stop := make(chan struct{})
	var wg sync.WaitGroup
	bad := make(chan time.Duration, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				if pos := p.Position(); pos < 0 {
					bad <- pos
				}
				p.Paused()
				p.Duration()
			}
		}()
	}
	select {
	case <-p.Events():
	case <-time.After(6 * time.Second):
		t.Error("no event within 6 s")
	}
	close(stop)
	wg.Wait()
	close(bad)
	for pos := range bad {
		t.Fatalf("Position() = %v, want at least 0", pos)
	}
}

// END: TestReadWhileEventsUpdate

// START: TestCloseWhileCommandsPending

func TestCloseWhileCommandsPending(t *testing.T) {
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 50; j++ {
				_, _ = p.Property("idle-active")
			}
		}()
	}
	close(start)
	time.Sleep(5 * time.Millisecond)
	p.Close()
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("calls still pending 2 s after Close")
	}
}

// END: TestCloseWhileCommandsPending
