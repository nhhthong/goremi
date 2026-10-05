//go:build !windows

// Test that the user's own mpv settings do not change goremi's playback.
package player

import (
	"os"
	"path/filepath"
	"testing"
)

// START: TestUserMpvConfNotApplied

func TestUserMpvConfNotApplied(t *testing.T) {
	cfg := t.TempDir()
	t.Setenv("HOME", cfg)
	t.Setenv("XDG_CONFIG_HOME", cfg)
	if err := os.MkdirAll(filepath.Join(cfg, "mpv"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfg, "mpv", "mpv.conf"), []byte("pause=yes\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := startPlaying(t, 4)
	if p.Paused() {
		t.Fatal("Paused() = true: the pause=yes of the user's mpv.conf was applied")
	}
}

// END: TestUserMpvConfNotApplied
