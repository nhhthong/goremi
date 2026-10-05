//go:build !windows

// Tests that closing the program while mpv is still starting neither waits for it nor leaves it running.
package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: helpers

// slowMpv puts a fake mpv first on PATH that never opens its socket, isolates TMPDIR and returns that directory.
func slowMpv(t *testing.T) string {
	t.Helper()
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "mpv"), []byte("#!/bin/sh\nexec /bin/sleep 987\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	t.Cleanup(func() { _ = exec.Command("pkill", "-f", "sleep 98[7]").Run() })
	return tmp
}

// fakeMpvRunning tells whether the fake mpv process exists.
func fakeMpvRunning() bool { return exec.Command("pgrep", "-f", "sleep 98[7]").Run() == nil }

// runPlaying runs the program with a fake run that starts Play in a goroutine and returns; with waitStart it first waits until the fake mpv is running. It returns when the fake run ended and when Run returned.
func runPlaying(t *testing.T, waitStart bool) (bodyEnd, runEnd time.Time) {
	t.Helper()
	run := func(m tea.Model) (tea.Model, error) {
		_, cmd := m.Update(PlayMsg{Track: provider.Track{ID: "a1", Title: "One"}})
		_, play := m.Update(cmd())
		go play()
		if waitStart {
			for deadline := time.Now().Add(3 * time.Second); !fakeMpvRunning() && time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
			}
			if !fakeMpvRunning() {
				t.Error("the fake mpv never started")
			}
		}
		bodyEnd = time.Now()
		return m, nil
	}
	if err := Run(nil, filepath.Join(t.TempDir(), "config.toml"), wavProvider{path: "x.wav"}, os.Stderr, run); err != nil {
		t.Fatalf("Run error: %v", err)
	}
	return bodyEnd, time.Now()
}

// nothingLeft waits up to 2 s for no fake mpv process and no goremi-* directory in tmp.
func nothingLeft(t *testing.T, tmp string) bool {
	t.Helper()
	for deadline := time.Now().Add(2 * time.Second); time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
		if !fakeMpvRunning() && len(endpointDirs(t, tmp)) == 0 {
			return true
		}
	}
	return false
}

// END: helpers

// START: TestRunClosesWhileMpvStarts

func TestRunClosesWhileMpvStarts(t *testing.T) {
	slowMpv(t)
	bodyEnd, runEnd := runPlaying(t, true)
	if d := runEnd.Sub(bodyEnd); d > 500*time.Millisecond {
		t.Fatalf("Run took %v after the fake run ended while mpv was starting, want under 500ms", d)
	}
}

// END: TestRunClosesWhileMpvStarts

// START: TestRunLeavesNoMpvAfterSlowStart

func TestRunLeavesNoMpvAfterSlowStart(t *testing.T) {
	tmp := slowMpv(t)
	runPlaying(t, true)
	if !nothingLeft(t, tmp) {
		t.Fatalf("after Run: fake mpv running = %v, directories %v, want none within 2s", fakeMpvRunning(), endpointDirs(t, tmp))
	}
}

// END: TestRunLeavesNoMpvAfterSlowStart

// START: TestCloseRacesPlay

func TestCloseRacesPlay(t *testing.T) {
	tmp := slowMpv(t)
	start := time.Now()
	runPlaying(t, false)
	if d := time.Since(start); d > 2*time.Second {
		t.Fatalf("Run took %v while Close raced Play, want under 2s", d)
	}
	if !nothingLeft(t, tmp) {
		t.Fatalf("after Run: fake mpv running = %v, directories %v, want none within 2s", fakeMpvRunning(), endpointDirs(t, tmp))
	}
}

// END: TestCloseRacesPlay
