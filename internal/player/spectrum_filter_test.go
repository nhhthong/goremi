// Tests for the mpv filter text of the spectrum: the branches, the log-spaced centres, the widths and what each branch measures.
package player

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// centres returns the f= value of every bandpass in the filter, in order.
func centres(t *testing.T) []float64 {
	t.Helper()
	var out []float64
	for _, m := range regexp.MustCompile(`bandpass=f=([0-9.]+):`).FindAllStringSubmatch(SpectrumFilter(), -1) {
		f, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, f)
	}
	return out
}

// START: TestSpectrumFilterHas32Branches

func TestSpectrumFilterHas32Branches(t *testing.T) {
	f := SpectrumFilter()
	if n := strings.Count(f, "bandpass="); n != 32 {
		t.Fatalf("bandpass count = %d, want 32", n)
	}
	if n := strings.Count(f, "ametadata=mode=print,anullsink"); n != 32 {
		t.Fatalf("branches ending in print and anullsink = %d, want 32", n)
	}
	for i := 0; i < 32; i++ {
		if want := fmt.Sprintf("ametadata=mode=add:key=lavfi.band:value=%d,", i); !strings.Contains(f, want) {
			t.Fatalf("filter has no %q", want)
		}
	}
}

// END: TestSpectrumFilterHas32Branches

// START: TestSpectrumFilterCentresAreLogSpaced

func TestSpectrumFilterCentresAreLogSpaced(t *testing.T) {
	c := centres(t)
	if len(c) != 32 || c[0] != 40 || c[31] != 16000 {
		t.Fatalf("centres = %v, want 32 from 40 to 16000", c)
	}
	ratio := math.Pow(400, 1.0/31)
	for i := 1; i < len(c); i++ {
		if got := c[i] / c[i-1]; math.Abs(got-ratio)/ratio > 0.01 {
			t.Fatalf("centre %d / centre %d = %.4f, want %.4f within 1 %%", i, i-1, got, ratio)
		}
	}
}

// END: TestSpectrumFilterCentresAreLogSpaced

// START: TestSpectrumFilterBandwidthIsQuarterOfCentre

func TestSpectrumFilterBandwidthIsQuarterOfCentre(t *testing.T) {
	ms := regexp.MustCompile(`bandpass=f=([0-9.]+):width_type=h:w=([0-9.]+),`).FindAllStringSubmatch(SpectrumFilter(), -1)
	if len(ms) != 32 {
		t.Fatalf("bandpass with width_type=h = %d, want 32", len(ms))
	}
	for _, m := range ms {
		f, _ := strconv.ParseFloat(m[1], 64)
		w, _ := strconv.ParseFloat(m[2], 64)
		if math.Abs(w-f/4) > 0.1 {
			t.Fatalf("f = %v, w = %v, want w = f/4", f, w)
		}
	}
}

// END: TestSpectrumFilterBandwidthIsQuarterOfCentre

// START: TestSpectrumFilterMeasuresRMSOnly

func TestSpectrumFilterMeasuresRMSOnly(t *testing.T) {
	const astats = "astats=metadata=1:reset=1:measure_perchannel=none:measure_overall=RMS_level,"
	if n := strings.Count(SpectrumFilter(), astats); n != 32 {
		t.Fatalf("branches with %q = %d, want 32", astats, n)
	}
}

// END: TestSpectrumFilterMeasuresRMSOnly
