// Tests for the play line in the log: a play names its track before its stream is resolved.
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"goremi/internal/provider"
)

// START: playLogHelpers

// orderProvider records "resolve" in the shared event list when asked for a stream.
type orderProvider struct {
	fakeProvider
	events *[]string
}

func (o orderProvider) Resolve(provider.Track) (string, error) {
	*o.events = append(*o.events, "resolve")
	return "http://stream/1", nil
}

// END: playLogHelpers

// START: TestPlayLogsPlayLineBeforeResolve

func TestPlayLogsPlayLineBeforeResolve(t *testing.T) {
	var events []string
	logf := func(format string, a ...any) { events = append(events, strings.TrimSpace(fmt.Sprintf(format, a...))) }
	m := New(orderProvider{events: &events}).WithLog(logf)
	_, cmd := m.Update(PlayMsg{Track: provider.Track{ID: "p1", Title: "One"}})
	cmdMsgs(cmd)
	want := []string{`play p1 "One"`, "resolve"}
	if len(events) != 2 || events[0] != want[0] || events[1] != want[1] {
		t.Fatalf("events = %q, want %q", events, want)
	}
}

// END: TestPlayLogsPlayLineBeforeResolve

// START: TestPlayWithoutLogStillResolves

func TestPlayWithoutLogStillResolves(t *testing.T) {
	var events []string
	_, cmd := New(orderProvider{events: &events}).Update(PlayMsg{Track: provider.Track{ID: "p1", Title: "One"}})
	msgs := cmdMsgs(cmd)
	if len(msgs) != 1 {
		t.Fatalf("the play command yielded %d messages, want 1", len(msgs))
	}
	if _, ok := msgs[0].(resolvedMsg); !ok {
		t.Fatalf("the message is %#v, want a resolvedMsg", msgs[0])
	}
}

// END: TestPlayWithoutLogStillResolves

// START: TestPlayLineReachesLogFile

func TestPlayLineReachesLogFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "goremi.log")
	lg, err := OpenLog(path)
	if err != nil {
		t.Fatal(err)
	}
	closeLog(t, lg)
	m := New(&resolveProvider{}).WithLog(lg.Printf)
	_, cmd := m.Update(PlayMsg{Track: provider.Track{ID: "p1", Title: "One"}})
	cmdMsgs(cmd)
	text, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if want := fmt.Sprintf("%d play p1 \"One\"\n", os.Getpid()); string(text) != want {
		t.Fatalf("log = %q, want %q", text, want)
	}
}

// END: TestPlayLineReachesLogFile

// START: TestPlayLineCannotForgeLogLines

func TestPlayLineCannotForgeLogLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "goremi.log")
	lg, err := OpenLog(path)
	if err != nil {
		t.Fatal(err)
	}
	closeLog(t, lg)
	m := New(&resolveProvider{}).WithLog(lg.Printf)
	_, cmd := m.Update(PlayMsg{Track: provider.Track{ID: "p1", Title: "x\"\n999 mpv [ffmpeg] fake"}})
	cmdMsgs(cmd)
	text, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSuffix(string(text), "\n"), "\n")
	if len(lines) != 1 || !strings.HasPrefix(lines[0], fmt.Sprintf("%d play p1 ", os.Getpid())) {
		t.Fatalf("the play wrote %d lines, want 1 play line:\n%s", len(lines), text)
	}
}

// END: TestPlayLineCannotForgeLogLines
