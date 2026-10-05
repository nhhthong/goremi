// Tests that playing a track makes no artwork request: the artwork is dropped from the product.
package app

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"

	"goremi/internal/provider"
)

// START: noArtworkHelpers

// cmdMsgs runs cmd and returns the messages it yields, opening the batches.
func cmdMsgs(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}
	msg := cmd()
	if batch, ok := msg.(tea.BatchMsg); ok {
		var out []tea.Msg
		for _, c := range batch {
			out = append(out, cmdMsgs(c)...)
		}
		return out
	}
	return []tea.Msg{msg}
}

// playOn plays a track whose artwork URL is a test server on a terminal with this environment; it returns the track, the messages of the play command and the count of requests the server got.
func playOn(t *testing.T, env uv.Environ) (provider.Track, []tea.Msg, *atomic.Int32) {
	t.Helper()
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		w.Write([]byte("not an image"))
	}))
	t.Cleanup(srv.Close)
	track := provider.Track{ID: "noart1", Title: "One", Artwork: srv.URL + "/a.jpg"}
	m, _ := New(&resolveProvider{}).Update(tea.EnvMsg(env))
	_, cmd := m.Update(PlayMsg{Track: track})
	return track, cmdMsgs(cmd), &hits
}

// END: noArtworkHelpers

// START: TestPlayKittyMakesNoArtworkRequest

func TestPlayKittyMakesNoArtworkRequest(t *testing.T) {
	_, _, hits := playOn(t, uv.Environ{"TERM=xterm-256color", "KITTY_WINDOW_ID=1"})
	if n := hits.Load(); n != 0 {
		t.Fatalf("the artwork server got %d requests, want 0", n)
	}
}

// END: TestPlayKittyMakesNoArtworkRequest

// START: TestPlayPlainTerminalMakesNoArtworkRequest

func TestPlayPlainTerminalMakesNoArtworkRequest(t *testing.T) {
	_, _, hits := playOn(t, uv.Environ{"TERM=xterm-256color"})
	if n := hits.Load(); n != 0 {
		t.Fatalf("the artwork server got %d requests, want 0", n)
	}
}

// END: TestPlayPlainTerminalMakesNoArtworkRequest

// START: TestPlayKittyYieldsOnlyResolve

func TestPlayKittyYieldsOnlyResolve(t *testing.T) {
	track, msgs, _ := playOn(t, uv.Environ{"TERM=xterm-256color", "KITTY_WINDOW_ID=1"})
	if len(msgs) != 1 {
		t.Fatalf("the play command yielded %d messages, want 1: %#v", len(msgs), msgs)
	}
	r, ok := msgs[0].(resolvedMsg)
	if !ok || r.track.ID != track.ID {
		t.Fatalf("the message is %#v, want a resolvedMsg for track %q", msgs[0], track.ID)
	}
}

// END: TestPlayKittyYieldsOnlyResolve
