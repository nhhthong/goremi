// The spectrum in the player: the mpv filter text that measures 32 bands.
package player

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// SpectrumBands is the number of bars the spectrum draws.
const SpectrumBands = 32

// START: SpectrumFilter

// bandCentre is the centre of band i in Hz: 32 bands on a log scale from 40 Hz to 16 kHz.
func bandCentre(i int) float64 {
	return 40 * math.Pow(16000.0/40, float64(i)/(SpectrumBands-1))
}

// SpectrumFilter is the value of the mpv option --af for the spectrum: one branch per band, each a band-pass of f/4 Hz,
// measured with astats, tagged with its number and printed, then dropped; the audio goes on unchanged through [o].
func SpectrumFilter() string {
	var b strings.Builder
	fmt.Fprintf(&b, "lavfi=[asplit=%d[o]", SpectrumBands+1)
	for i := 0; i < SpectrumBands; i++ {
		fmt.Fprintf(&b, "[b%d]", i)
	}
	for i := 0; i < SpectrumBands; i++ {
		f := bandCentre(i)
		fmt.Fprintf(&b, ";[b%d]bandpass=f=%.2f:width_type=h:w=%.2f,astats=metadata=1:reset=1:measure_perchannel=none:measure_overall=RMS_level,ametadata=mode=add:key=lavfi.band:value=%d,ametadata=mode=print,anullsink", i, f, f/4, i)
	}
	b.WriteString(";[o]anull]")
	return b.String()
}

// END: SpectrumFilter

// START: Bands

// floorDB is the level of a band nobody has reported, or whose level cannot be read.
const floorDB = -60

// bandState holds the latest level of every band and, per filter instance, the level line that waits for its band line.
type bandState struct {
	mu      sync.Mutex
	db      [SpectrumBands]float64
	seen    [SpectrumBands]bool
	pending map[string]float64
}

var (
	levelLine = regexp.MustCompile(`^(Parsed_ametadata_\d+): lavfi\.astats\.Overall\.RMS_level=(\S+)`)
	bandLine  = regexp.MustCompile(`^(Parsed_ametadata_\d+): lavfi\.band=(\d+)`)
)

// bandMessage reads one log message of mpv: the level line of a filter instance is kept, and the band line of the same instance gives that level to its band.
// A level that is -inf, not a number or unreadable counts as -60 dB; a band outside 0 to 31, or with no level before it, is ignored.
func (p *Player) bandMessage(prefix, text string) {
	if prefix != "ffmpeg" {
		return
	}
	b := &p.bands
	if m := levelLine.FindStringSubmatch(text); m != nil {
		v, err := strconv.ParseFloat(m[2], 64)
		if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
			v = floorDB
		}
		b.mu.Lock()
		if b.pending == nil {
			b.pending = map[string]float64{}
		}
		b.pending[m[1]] = v
		b.mu.Unlock()
		return
	}
	if m := bandLine.FindStringSubmatch(text); m != nil {
		i, err := strconv.Atoi(m[2])
		if err != nil || i < 0 || i >= SpectrumBands {
			return
		}
		b.mu.Lock()
		if v, ok := b.pending[m[1]]; ok {
			b.db[i], b.seen[i] = v, true
		}
		b.mu.Unlock()
	}
}

// Bands is the latest level of each of the 32 bands in dB; -60 for a band that has not reported.
func (p *Player) Bands() [SpectrumBands]float64 {
	b := &p.bands
	b.mu.Lock()
	defer b.mu.Unlock()
	var out [SpectrumBands]float64
	for i := range out {
		out[i] = floorDB
		if b.seen[i] {
			out[i] = b.db[i]
		}
	}
	return out
}

// Spectrum tells whether mpv runs with the band filter: false when it was not asked for, or when mpv refused it and started without.
func (p *Player) Spectrum() bool { return p.spectrum }

// END: Bands
