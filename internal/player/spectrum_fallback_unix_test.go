//go:build !windows

// Tests for starting without the spectrum when mpv refuses the band filter: an mpv on PATH that exits when it gets --af stands for one that cannot run it.
package player

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// mpvThatRefusesFilter puts a script named mpv first on PATH: it exits 1 when an argument starts with --af, and runs the real mpv otherwise.
func mpvThatRefusesFilter(t *testing.T) {
	t.Helper()
	real, err := exec.LookPath("mpv")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	script := "#!/bin/sh\nfor a in \"$@\"; do case \"$a\" in --af=*) exit 1;; esac; done\nexec " + real + " \"$@\"\n"
	if err := os.WriteFile(filepath.Join(dir, "mpv"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// START: TestStartFallsBackWithoutFilter

func TestStartFallsBackWithoutFilter(t *testing.T) {
	mpvThatRefusesFilter(t)
	p, err := StartWith(context.Background(), Options{Spectrum: true})
	if err != nil {
		t.Fatalf("StartWith returned %v, want a player without the spectrum", err)
	}
	defer p.Close()
	if p.Spectrum() {
		t.Fatal("Spectrum() is true, want false after the fallback")
	}
}

// END: TestStartFallsBackWithoutFilter

// START: TestFallbackPlaysAndLogs

func TestFallbackPlaysAndLogs(t *testing.T) {
	mpvThatRefusesFilter(t)
	r := &recorder{}
	p, err := StartWith(context.Background(), Options{Spectrum: true, Log: r.logf})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Play("av://lavfi:sine=frequency=440:duration=3"); err != nil {
		t.Fatalf("Play returned %v, want nil", err)
	}
	af, err := p.Property("af")
	if list, ok := af.([]any); err != nil || !ok || len(list) != 0 {
		t.Fatalf("af = %#v (%v), want an empty list", af, err)
	}
	found := false
	for _, l := range r.get() {
		if strings.Contains(l, "did not start with the spectrum filter") {
			found = true
		}
	}
	if !found {
		t.Fatalf("the log has no line about the filter: %q", r.get())
	}
}

// END: TestFallbackPlaysAndLogs

// START: TestStartKeepsSpectrumWhenMpvAccepts

func TestStartKeepsSpectrumWhenMpvAccepts(t *testing.T) {
	p, err := StartWith(context.Background(), Options{Spectrum: true})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if !p.Spectrum() {
		t.Fatal("Spectrum() is false, want true with an mpv that accepts the filter")
	}
}

// END: TestStartKeepsSpectrumWhenMpvAccepts
