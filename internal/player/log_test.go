// Tests for the mpv log messages that reach the log.
package player

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// START: recorder

// recorder collects log calls.
type recorder struct {
	mu    sync.Mutex
	lines []string
}

func (r *recorder) logf(format string, a ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lines = append(r.lines, fmt.Sprintf(format, a...))
}

func (r *recorder) get() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.lines...)
}

// END: recorder

// START: logMessage tests

func TestLogMessageWarn(t *testing.T) {
	r := &recorder{}
	(&Player{log: r.logf}).logMessage("warn", "ffmpeg", "https: HTTP error 403 Forbidden\n")
	if got := r.get(); len(got) != 1 || got[0] != "mpv [ffmpeg] https: HTTP error 403 Forbidden" {
		t.Fatalf("log = %q", got)
	}
}

func TestLogMessageError(t *testing.T) {
	r := &recorder{}
	(&Player{log: r.logf}).logMessage("error", "stream", "Failed to open x.\n")
	if got := r.get(); len(got) != 1 || got[0] != "mpv [stream] Failed to open x." {
		t.Fatalf("log = %q", got)
	}
}

func TestLogMessageOtherLevels(t *testing.T) {
	r := &recorder{}
	p := &Player{log: r.logf}
	p.logMessage("info", "cplayer", "playing\n")
	p.logMessage("v", "cplayer", "verbose\n")
	if got := r.get(); len(got) != 0 {
		t.Fatalf("log = %q, want nothing", got)
	}
}

func TestLogMessageSkipsBandLine(t *testing.T) {
	r := &recorder{}
	(&Player{log: r.logf}).logMessage("warn", "ffmpeg", "lavfi.band=3\n")
	if got := r.get(); len(got) != 0 {
		t.Fatalf("log = %q, want nothing", got)
	}
}

func TestLogMessageSkipsAstatsLine(t *testing.T) {
	r := &recorder{}
	(&Player{log: r.logf}).logMessage("warn", "ffmpeg", "lavfi.astats.Overall.RMS_level=-41.2\n")
	if got := r.get(); len(got) != 0 {
		t.Fatalf("log = %q, want nothing", got)
	}
}

func TestLogMessageKeepsMentionOfBand(t *testing.T) {
	r := &recorder{}
	(&Player{log: r.logf}).logMessage("warn", "ffmpeg", "filter lavfi.band=3 failed\n")
	if got := r.get(); len(got) != 1 {
		t.Fatalf("log = %q, want one line", got)
	}
}

// END: logMessage tests

// START: TestStartLoggedReportsMpvErrors

func TestStartLoggedReportsMpvErrors(t *testing.T) {
	r := &recorder{}
	p, err := StartLogged(t.Context(), r.logf)
	if err != nil {
		t.Skipf("mpv not available: %v", err)
	}
	defer p.Close()
	_ = p.Play("/nonexistent/file.mp3")
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); time.Sleep(50 * time.Millisecond) {
		for _, l := range r.get() {
			if strings.Contains(l, "nonexistent") {
				return
			}
		}
	}
	t.Fatalf("no log line about the missing file, log = %q", r.get())
}

// END: TestStartLoggedReportsMpvErrors
