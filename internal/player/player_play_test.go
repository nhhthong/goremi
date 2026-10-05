// Tests for playing a file: it starts, and the position grows.
package player

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// START: writeWAV

// writeWAV writes a mono 8 kHz 440 Hz sine of the given length in seconds as a 16-bit WAV file and returns its path.
func writeWAV(t *testing.T, seconds int) string {
	t.Helper()
	const rate = 8000
	n := rate * seconds
	data := make([]byte, 0, 44+2*n)
	le := binary.LittleEndian
	data = append(data, "RIFF"...)
	data = le.AppendUint32(data, uint32(36+2*n))
	data = append(data, "WAVEfmt "...)
	data = le.AppendUint32(data, 16)
	data = le.AppendUint16(data, 1) // PCM
	data = le.AppendUint16(data, 1) // mono
	data = le.AppendUint32(data, rate)
	data = le.AppendUint32(data, 2*rate)
	data = le.AppendUint16(data, 2)
	data = le.AppendUint16(data, 16)
	data = append(data, "data"...)
	data = le.AppendUint32(data, uint32(2*n))
	for i := 0; i < n; i++ {
		data = le.AppendUint16(data, uint16(int16(8000*math.Sin(2*math.Pi*440*float64(i)/rate))))
	}
	path := filepath.Join(t.TempDir(), "sine.wav")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// END: writeWAV

// START: playing

// waitPosition waits until the position is above 0, and fails after 3 seconds.
func waitPosition(t *testing.T, p *Player) {
	t.Helper()
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); time.Sleep(50 * time.Millisecond) {
		if p.Position() > 0 {
			return
		}
	}
	t.Fatal("Position() still 0 after 3 s")
}

// startPlaying starts a player on a sine file of the given length and waits until the position is above 0.
func startPlaying(t *testing.T, seconds int) *Player {
	t.Helper()
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	t.Cleanup(p.Close)
	if err := p.Play(writeWAV(t, seconds)); err != nil {
		t.Fatalf("Play() error: %v", err)
	}
	waitPosition(t, p)
	return p
}

// END: playing

// START: TestPlayStartsPlayback

func TestPlayStartsPlayback(t *testing.T) {
	p := startPlaying(t, 4)
	if p.Paused() {
		t.Fatal("Paused() = true after Play, want false")
	}
}

// END: TestPlayStartsPlayback

// START: TestPositionGrows

func TestPositionGrows(t *testing.T) {
	p := startPlaying(t, 4)
	before := p.Position()
	time.Sleep(500 * time.Millisecond)
	if after := p.Position(); after <= before {
		t.Fatalf("Position() %v then %v, want it to grow", before, after)
	}
}

// END: TestPositionGrows
