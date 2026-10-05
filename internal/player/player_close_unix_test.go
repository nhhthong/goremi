//go:build !windows

// Tests for Close and for mpv ending on its own; they look at the private directory and the process from outside.
package player

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// START: helpers

// isolatedTmp points TMPDIR at a new directory and returns it, so the test sees only the directories Start makes.
func isolatedTmp(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("TMPDIR", dir)
	return dir
}

// endpointDirs lists the goremi-* directories inside tmp.
func endpointDirs(t *testing.T, tmp string) []string {
	t.Helper()
	dirs, err := filepath.Glob(filepath.Join(tmp, "goremi-*"))
	if err != nil {
		t.Fatal(err)
	}
	return dirs
}

// mpvPids returns the pids of the processes whose command line holds dir; none gives nil.
func mpvPids(dir string) []int {
	out, _ := exec.Command("pgrep", "-f", dir).Output()
	var pids []int
	for _, f := range strings.Fields(string(out)) {
		if pid, err := strconv.Atoi(f); err == nil {
			pids = append(pids, pid)
		}
	}
	return pids
}

// END: helpers

// START: TestCloseRemovesDir

func TestCloseRemovesDir(t *testing.T) {
	tmp := isolatedTmp(t)
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	p.Close()
	if dirs := endpointDirs(t, tmp); len(dirs) != 0 {
		t.Fatalf("directories left after Close: %v", dirs)
	}
}

// END: TestCloseRemovesDir

// START: TestCloseStopsMpv

func TestCloseStopsMpv(t *testing.T) {
	tmp := isolatedTmp(t)
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	dirs := endpointDirs(t, tmp)
	if len(dirs) != 1 || len(mpvPids(dirs[0])) == 0 {
		t.Fatalf("before Close: directories %v and no mpv process running", dirs)
	}
	p.Close()
	if pids := mpvPids(dirs[0]); len(pids) != 0 {
		t.Fatalf("mpv processes left after Close: %v", pids)
	}
}

// END: TestCloseStopsMpv

// START: TestFailedStartLeavesNoDir

func TestFailedStartLeavesNoDir(t *testing.T) {
	tmp := isolatedTmp(t)
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "mpv"), []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	p, err := Start()
	if p != nil || err == nil {
		t.Fatalf("Start() = %v, %v, want nil and an error", p, err)
	}
	if dirs := endpointDirs(t, tmp); len(dirs) != 0 {
		t.Fatalf("directories left after a failed Start: %v", dirs)
	}
}

// END: TestFailedStartLeavesNoDir

// START: TestMpvExitReportsDone

func TestMpvExitReportsDone(t *testing.T) {
	tmp := isolatedTmp(t)
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer p.Close()
	killMpv(t, endpointDirs(t, tmp)[0])
	select {
	case e := <-p.Events():
		if e != Done {
			t.Fatalf("event = %v, want Done", e)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no event within 2 s after mpv was killed")
	}
}

// killMpv kills the mpv process that uses dir.
func killMpv(t *testing.T, dir string) {
	t.Helper()
	pids := mpvPids(dir)
	if len(pids) == 0 {
		t.Fatal("no mpv process to kill")
	}
	for _, pid := range pids {
		_ = syscall.Kill(pid, syscall.SIGKILL)
	}
}

// END: TestMpvExitReportsDone

// START: TestCallsAfterDoneReturnError

func TestCallsAfterDoneReturnError(t *testing.T) {
	tmp := isolatedTmp(t)
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer p.Close()
	killMpv(t, endpointDirs(t, tmp)[0])
	select {
	case <-p.Events():
	case <-time.After(2 * time.Second):
		t.Fatal("no Done event within 2 s")
	}
	if _, err := p.Property("idle-active"); err == nil {
		t.Error("Property after Done returned no error")
	}
	if err := p.Play("file.wav"); err == nil {
		t.Error("Play after Done returned no error")
	}
}

// END: TestCallsAfterDoneReturnError
