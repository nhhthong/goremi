// Tests for the line under Search: when the player cannot play a track.
package app

import (
	"errors"
	"os/exec"
	"testing"

	"goremi/internal/provider"
)

const mpvMissing = "mpv not found. Install mpv and try again."

// playWith plays track One through a fake player whose Play returns errs in turn and returns the model and the player.
func playWith(errs ...error) (Model, *recordingPlayer) {
	pl := &recordingPlayer{playErrs: errs}
	m, _ := playOne(New(&scriptedResolve{errs: []error{nil, nil}}).WithPlayer(pl), provider.Track{ID: "a1", Title: "One"})
	return m, pl
}

// START: TestMpvNotFoundNotice

func TestMpvNotFoundNotice(t *testing.T) {
	m, _ := playWith(exec.ErrNotFound)
	if got := noticeLine(m); got != mpvMissing {
		t.Fatalf("line under Search: = %q, want %q", got, mpvMissing)
	}
}

// END: TestMpvNotFoundNotice

// START: TestOtherPlayErrorNotice

func TestOtherPlayErrorNotice(t *testing.T) {
	m, _ := playWith(errors.New("boom"))
	if got, want := noticeLine(m), `Cannot play "One".`; got != want {
		t.Fatalf("line under Search: = %q, want %q", got, want)
	}
}

// END: TestOtherPlayErrorNotice

// START: TestPlayRecoversAfterMpvNotFound

func TestPlayRecoversAfterMpvNotFound(t *testing.T) {
	m, pl := playWith(exec.ErrNotFound, nil)
	m, _ = playOne(m, provider.Track{ID: "a2", Title: "Two"})
	if len(pl.urls) != 2 || noticeLine(m) != "" {
		t.Fatalf("Play calls = %d, line = %q, want 2 calls and no line", len(pl.urls), noticeLine(m))
	}
}

// END: TestPlayRecoversAfterMpvNotFound
