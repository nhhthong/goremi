// Tests for Play running off the screen loop: Update returns at once, plays keep their order, and a late failure still reaches the screen.
package app

import (
	"os/exec"
	"sync"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/player"
	"goremi/internal/provider"
)

// gatedPlayer is a Player whose first Play waits until gate is closed; it records the URLs it played and the pauses it was asked for.
type gatedPlayer struct {
	mu       sync.Mutex
	gate     chan struct{}
	started  bool
	urls     []string
	toggles  int
	playsErr error
}

func newGatedPlayer() *gatedPlayer { return &gatedPlayer{gate: make(chan struct{})} }

func (g *gatedPlayer) Play(url string) error {
	g.mu.Lock()
	first := !g.started
	g.started = true
	g.mu.Unlock()
	if first {
		<-g.gate
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.urls = append(g.urls, url)
	return g.playsErr
}

func (g *gatedPlayer) TogglePause() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.toggles++
	return nil
}

func (g *gatedPlayer) Position() time.Duration     { return 0 }
func (g *gatedPlayer) Duration() time.Duration     { return 0 }
func (g *gatedPlayer) Paused() bool                { return false }
func (g *gatedPlayer) Seek(float64) error          { return nil }
func (g *gatedPlayer) Events() <-chan player.Event { return nil }
func (g *gatedPlayer) played() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]string(nil), g.urls...)
}
func (g *gatedPlayer) pauses() int { g.mu.Lock(); defer g.mu.Unlock(); return g.toggles }

// resolved is the message of a resolved track.
func resolved(track provider.Track, url string) resolvedMsg {
	return resolvedMsg{track: track, url: url}
}

// START: TestUpdateReturnsWhilePlayBlocks

func TestUpdateReturnsWhilePlayBlocks(t *testing.T) {
	g := newGatedPlayer()
	m := New(fakeProvider{}).WithPlayer(g)
	start := time.Now()
	_, cmd := m.Update(resolved(trackA, "urlA"))
	elapsed, playedDuring := time.Since(start), g.played()
	close(g.gate)
	cmd()
	if elapsed > 100*time.Millisecond || len(playedDuring) != 0 {
		t.Fatalf("Update took %v with Play blocked and %v played, want under 100ms and nothing played yet", elapsed, playedDuring)
	}
}

// END: TestUpdateReturnsWhilePlayBlocks

// START: TestPlaysKeepRequestOrder

func TestPlaysKeepRequestOrder(t *testing.T) {
	g := newGatedPlayer()
	m := New(fakeProvider{}).WithPlayer(g)
	next, first := m.Update(resolved(trackA, "urlA"))
	_, second := next.Update(resolved(trackB, "urlB"))
	var wg sync.WaitGroup
	for _, c := range []tea.Cmd{first, second} {
		wg.Add(1)
		go func() { defer wg.Done(); c() }()
	}
	time.Sleep(5 * time.Millisecond) // both plays are queued while the first Play still blocks
	close(g.gate)
	wg.Wait()
	if got := g.played(); len(got) != 2 || got[0] != "urlA" || got[1] != "urlB" {
		t.Fatalf("Play order = %v, want [urlA urlB]", got)
	}
}

// END: TestPlaysKeepRequestOrder

// START: TestKeysWorkWhilePlayBlocks

func TestKeysWorkWhilePlayBlocks(t *testing.T) {
	g := newGatedPlayer()
	m := New(fakeProvider{}).WithPlayer(g)
	next, cmd := m.Update(resolved(trackA, "urlA"))
	next.(Model).WithFocus(FocusList).Update(tea.KeyPressMsg{Code: 'k', Text: "k"})
	during := g.pauses()
	close(g.gate)
	cmd()
	if during != 1 {
		t.Fatalf("TogglePause calls while Play blocked = %d, want 1", during)
	}
}

// END: TestKeysWorkWhilePlayBlocks

// START: TestLateMpvNotFoundNotice

func TestLateMpvNotFoundNotice(t *testing.T) {
	g := newGatedPlayer()
	g.playsErr = exec.ErrNotFound
	m, _ := New(fakeProvider{}).WithPlayer(g).Update(PlayMsg{Track: trackA})
	next, cmd := m.Update(resolved(trackA, "urlA"))
	close(g.gate)
	done, _ := next.Update(cmd())
	if got := noticeLine(done.(Model)); got != mpvMissing {
		t.Fatalf("line under Search: = %q, want %q", got, mpvMissing)
	}
}

// END: TestLateMpvNotFoundNotice
