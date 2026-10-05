// Tests of the spectrum filter in a real mpv: it plays a tone, every band reports, the nearest band is the loudest.
package player

import (
	"bufio"
	"encoding/json"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"testing"
	"time"
)

// bandReading is what a real mpv reported for a 440 Hz tone: the highest level seen per band, and the number of messages of level error.
type bandReading struct {
	levels map[int]float64
	errors int
}

// readBands starts mpv with the spectrum filter, plays a 440 Hz tone and reads the band lines for 1.5 seconds.
func readBands(t *testing.T) bandReading {
	t.Helper()
	path, err := exec.LookPath("mpv")
	if err != nil {
		t.Fatal(err)
	}
	endpoint, cleanup, err := newEndpoint()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(path, "--no-config", "--idle=yes", "--no-terminal", "--no-video", "--ao=null", "--input-ipc-server="+endpoint, "--af="+SpectrumFilter())
	if err := cmd.Start(); err != nil {
		cleanup()
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait(); cleanup() })
	var conn interface {
		Write([]byte) (int, error)
		Read([]byte) (int, error)
		SetReadDeadline(time.Time) error
	}
	for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); time.Sleep(50 * time.Millisecond) {
		if c, err := dial(endpoint); err == nil {
			conn = c
			break
		}
	}
	if conn == nil {
		t.Fatal("mpv endpoint did not open")
	}
	for _, c := range []string{
		`{"command":["request_log_messages","v"]}` + "\n",
		`{"command":["loadfile","av://lavfi:sine=frequency=440:duration=3"]}` + "\n",
	} {
		if _, err := conn.Write([]byte(c)); err != nil {
			t.Fatal(err)
		}
	}
	rms := regexp.MustCompile(`^(Parsed_ametadata_\d+): lavfi\.astats\.Overall\.RMS_level=(\S+)`)
	band := regexp.MustCompile(`^(Parsed_ametadata_\d+): lavfi\.band=(\d+)`)
	last := map[string]float64{}
	got := bandReading{levels: map[int]float64{}}
	end := time.Now().Add(1500 * time.Millisecond)
	_ = conn.SetReadDeadline(end.Add(500 * time.Millisecond))
	r := bufio.NewReader(readerOf{conn})
	for time.Now().Before(end) {
		line, err := r.ReadBytes('\n')
		if err != nil {
			break
		}
		var m struct{ Event, Level, Text string }
		if json.Unmarshal(line, &m) != nil || m.Event != "log-message" {
			continue
		}
		if m.Level == "error" {
			got.errors++
		}
		if x := rms.FindStringSubmatch(m.Text); x != nil {
			v, _ := strconv.ParseFloat(x[2], 64)
			last[x[1]] = v
		}
		if x := band.FindStringSubmatch(m.Text); x != nil {
			i, _ := strconv.Atoi(x[2])
			v, ok := last[x[1]]
			if cur, seen := got.levels[i]; ok && (!seen || v > cur) {
				got.levels[i] = v
			}
		}
	}
	return got
}

// readerOf makes the read half of the connection an io.Reader.
type readerOf struct {
	c interface{ Read([]byte) (int, error) }
}

func (r readerOf) Read(b []byte) (int, error) { return r.c.Read(b) }

// START: TestRealMpvPlaysWithBandFilter

func TestRealMpvPlaysWithBandFilter(t *testing.T) {
	got := readBands(t)
	if got.errors != 0 {
		t.Fatalf("mpv sent %d messages of level error", got.errors)
	}
	if len(got.levels) != 32 {
		t.Fatalf("%d bands reported, want 32", len(got.levels))
	}
}

// END: TestRealMpvPlaysWithBandFilter

// START: TestRealMpvLoudestBandIsNearTone

func TestRealMpvLoudestBandIsNearTone(t *testing.T) {
	got := readBands(t)
	loudest, level := -1, math.Inf(-1)
	for i, v := range got.levels {
		if v > level {
			loudest, level = i, v
		}
	}
	if loudest != 12 && loudest != 13 {
		t.Fatalf("loudest band = %d (%v dB), want 12 or 13 (levels %v)", loudest, level, got.levels)
	}
}

// END: TestRealMpvLoudestBandIsNearTone
