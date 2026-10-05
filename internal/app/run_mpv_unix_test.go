//go:build !windows

// Tests that Run connects mpv to the app: it starts at the first play and stops when the command ends.
package app

import (
	"encoding/binary"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// wavProvider resolves every track to one WAV file.
type wavProvider struct {
	fakeProvider
	path string
}

func (w wavProvider) Resolve(provider.Track) (string, error) { return w.path, nil }

// START: helpers

// writeSilence writes a one-second silent WAV file and returns its path.
func writeSilence(t *testing.T) string {
	t.Helper()
	const rate = 8000
	data := []byte("RIFF")
	le := binary.LittleEndian
	data = le.AppendUint32(data, 36+2*rate)
	data = append(data, "WAVEfmt "...)
	data = le.AppendUint32(data, 16)
	data = le.AppendUint16(data, 1)
	data = le.AppendUint16(data, 1)
	data = le.AppendUint32(data, rate)
	data = le.AppendUint32(data, 2*rate)
	data = le.AppendUint16(data, 2)
	data = le.AppendUint16(data, 16)
	data = append(data, "data"...)
	data = le.AppendUint32(data, 2*rate)
	data = append(data, make([]byte, 2*rate)...)
	path := filepath.Join(t.TempDir(), "silence.wav")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// runWith runs the program with a fake run that calls do on the model; TMPDIR is a temp directory the test reads.
func runWith(t *testing.T, do func(m Model, tmp string)) (tmp string) {
	t.Helper()
	tmp = t.TempDir()
	t.Setenv("TMPDIR", tmp)
	p := wavProvider{path: writeSilence(t)}
	run := func(m tea.Model) (tea.Model, error) {
		do(m.(Model), tmp)
		return m, nil
	}
	if err := Run(nil, filepath.Join(t.TempDir(), "config.toml"), p, os.Stderr, run); err != nil {
		t.Fatalf("Run error: %v", err)
	}
	return tmp
}

// play sends PlayMsg for a track and feeds the resolved URL back, as the event loop does.
func play(m Model) {
	_, cmd := m.Update(PlayMsg{Track: provider.Track{ID: "a1", Title: "One"}})
	_, cmd = m.Update(cmd())
	cmd() // Play runs in this command
}

// endpointDirs lists the goremi-* directories in tmp.
func endpointDirs(t *testing.T, tmp string) []string {
	t.Helper()
	dirs, err := filepath.Glob(filepath.Join(tmp, "goremi-*"))
	if err != nil {
		t.Fatal(err)
	}
	return dirs
}

// mpvRunning tells whether a process has dir in its command line.
func mpvRunning(dir string) bool {
	return exec.Command("pgrep", "-f", dir).Run() == nil
}

// END: helpers

// START: TestRunStartsNoMpvUntilPlay

func TestRunStartsNoMpvUntilPlay(t *testing.T) {
	inside := -1
	tmp := runWith(t, func(m Model, tmp string) { inside = len(endpointDirs(t, tmp)) })
	if after := endpointDirs(t, tmp); inside != 0 || len(after) != 0 {
		t.Fatalf("directories inside Run: %d, after: %v, want none", inside, after)
	}
}

// END: TestRunStartsNoMpvUntilPlay

// START: TestRunClosesMpvAtEnd

func TestRunClosesMpvAtEnd(t *testing.T) {
	var dirs []string
	tmp := runWith(t, func(m Model, tmp string) {
		play(m)
		dirs = endpointDirs(t, tmp)
		if len(dirs) != 1 || !mpvRunning(dirs[0]) {
			t.Errorf("during Run: directories %v, want one with mpv running", dirs)
		}
	})
	if left := endpointDirs(t, tmp); len(left) != 0 || (len(dirs) == 1 && mpvRunning(dirs[0])) {
		t.Fatalf("after Run: directories %v, mpv still running: %v", left, len(dirs) == 1 && mpvRunning(dirs[0]))
	}
}

// END: TestRunClosesMpvAtEnd

// START: TestRunKeepsOneMpv

func TestRunKeepsOneMpv(t *testing.T) {
	runWith(t, func(m Model, tmp string) {
		play(m)
		play(m)
		out, _ := exec.Command("pgrep", "-f", tmp).Output()
		if dirs := endpointDirs(t, tmp); len(dirs) != 1 {
			t.Errorf("directories after two plays = %v, want one", dirs)
		}
		if n := len(strings.Fields(string(out))); n != 1 {
			t.Errorf("mpv processes after two plays = %d, want 1", n)
		}
	})
}

// END: TestRunKeepsOneMpv
