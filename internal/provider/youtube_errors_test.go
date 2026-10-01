// Tests for how Search reports yt-dlp failures: connection loss, other failures and timeout.
package provider

import (
	"errors"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"
)

// START: helpers

// failing returns a Runner that fails like yt-dlp does: an ExitError carrying the stderr.
func failing(stderr string) Runner {
	return func(...string) ([]byte, error) {
		return nil, &exec.ExitError{Stderr: []byte(stderr)}
	}
}

// searchWith runs one Search through a provider with the given Runner.
func searchWith(r Runner) ([]Track, error) {
	return (&YouTubeProvider{Runner: r}).Search("daft punk", 1)
}

// END: helpers

// START: connection loss

func TestSearchNewConnectionError(t *testing.T) {
	_, err := searchWith(failing("NewConnectionError"))
	if !errors.Is(err, ErrNetwork) {
		t.Fatalf("err = %v, want ErrNetwork", err)
	}
}

func TestSearchNameResolutionError(t *testing.T) {
	_, err := searchWith(failing("NameResolutionError"))
	if !errors.Is(err, ErrNetwork) {
		t.Fatalf("err = %v, want ErrNetwork", err)
	}
}

// observedOffline is the last stderr line yt-dlp printed on 2026-10-02 through a closed proxy port.
const observedOffline = `ERROR: Unable to download API page: ('Unable to connect to proxy', NewConnectionError("HTTPSConnection(host='127.0.0.1', port=1): Failed to establish a new connection: [Errno 111] Connection refused"))`

func TestSearchObservedOffline(t *testing.T) {
	tracks, err := searchWith(failing(observedOffline))
	if !errors.Is(err, ErrNetwork) || len(tracks) != 0 {
		t.Fatalf("tracks = %v, err = %v; want no tracks and ErrNetwork", tracks, err)
	}
}

func TestSearchRecoversAfterNetwork(t *testing.T) {
	calls := 0
	p := &YouTubeProvider{Runner: func(...string) ([]byte, error) {
		calls++
		if calls == 1 {
			return nil, &exec.ExitError{Stderr: []byte("NameResolutionError")}
		}
		return []byte(jsonLines(2)), nil
	}}
	if _, err := p.Search("daft punk", 1); !errors.Is(err, ErrNetwork) {
		t.Fatalf("first err = %v, want ErrNetwork", err)
	}
	tracks, err := p.Search("daft punk", 1)
	if err != nil || len(tracks) != 2 {
		t.Fatalf("second: %d tracks, err = %v; want 2 tracks and no error", len(tracks), err)
	}
}

// END: connection loss

// START: other failures

func TestSearchHTTP429NotNetwork(t *testing.T) {
	_, err := searchWith(failing("ERROR: Unable to download API page: HTTP Error 429: Too Many Requests"))
	if err == nil || errors.Is(err, ErrNetwork) {
		t.Fatalf("err = %v, want an error that is not ErrNetwork", err)
	}
}

func TestSearchEmptyStderrNotNetwork(t *testing.T) {
	_, err := searchWith(failing(""))
	if err == nil || errors.Is(err, ErrNetwork) {
		t.Fatalf("err = %v, want an error that is not ErrNetwork", err)
	}
}

// END: other failures

// START: timeout

// hanging returns a Runner that blocks until release is closed.
func hanging(release chan struct{}) Runner {
	return func(...string) ([]byte, error) {
		<-release
		return nil, nil
	}
}

func TestSearchTimesOut(t *testing.T) {
	release := make(chan struct{})
	defer close(release)
	p := &YouTubeProvider{Runner: hanging(release), Timeout: 50 * time.Millisecond}
	start := time.Now()
	tracks, err := p.Search("daft punk", 1)
	if !errors.Is(err, ErrTimeout) || len(tracks) != 0 || time.Since(start) > 2*time.Second {
		t.Fatalf("tracks = %v, err = %v after %v; want no tracks and ErrTimeout within 2s", tracks, err, time.Since(start))
	}
}

func TestSearchTimeoutIsTenSeconds(t *testing.T) {
	if SearchTimeout != 10*time.Second {
		t.Fatalf("SearchTimeout = %v, want 10s", SearchTimeout)
	}
}

func TestSearchRecoversAfterTimeout(t *testing.T) {
	release := make(chan struct{})
	defer close(release)
	var calls atomic.Int32
	p := &YouTubeProvider{Timeout: 50 * time.Millisecond, Runner: func(...string) ([]byte, error) {
		if calls.Add(1) == 1 {
			<-release
			return nil, nil
		}
		return []byte(jsonLines(2)), nil
	}}
	if _, err := p.Search("daft punk", 1); !errors.Is(err, ErrTimeout) {
		t.Fatalf("first err = %v, want ErrTimeout", err)
	}
	tracks, err := p.Search("daft punk", 1)
	if err != nil || len(tracks) != 2 {
		t.Fatalf("second: %d tracks, err = %v; want 2 tracks and no error", len(tracks), err)
	}
}

// END: timeout
