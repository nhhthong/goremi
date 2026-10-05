//go:build !windows

// Test that Run with no mpv on PATH tells the user under Search: instead of failing.
package app

import (
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: TestRunWithoutMpvNotice

func TestRunWithoutMpvNotice(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	got := ""
	run := func(m tea.Model) (tea.Model, error) {
		next, cmd := m.Update(PlayMsg{Track: provider.Track{ID: "a1", Title: "One"}})
		next, cmd = next.Update(cmd())
		next, _ = next.Update(cmd())
		got = noticeLine(next.(Model))
		return next, nil
	}
	if err := Run(nil, filepath.Join(t.TempDir(), "config.toml"), wavProvider{path: "x.wav"}, os.Stderr, run); err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if got != mpvMissing {
		t.Fatalf("line under Search: = %q, want %q", got, mpvMissing)
	}
}

// END: TestRunWithoutMpvNotice
