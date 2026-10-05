// Tests for reading the spectrum bands from the log messages of mpv, and for keeping those messages out of the log.
package player

import (
	"sync"
	"testing"
)

// level and band build the two ffmpeg messages a band filter instance prints, the level line first.
func level(inst int, db string) (string, string) {
	return "ffmpeg", "Parsed_ametadata_" + itoa(inst) + ": lavfi.astats.Overall.RMS_level=" + db + "\n"
}

func band(inst, i int) (string, string) {
	return "ffmpeg", "Parsed_ametadata_" + itoa(inst) + ": lavfi.band=" + itoa(i) + "\n"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for ; n > 0; n /= 10 {
		s = string(rune('0'+n%10)) + s
	}
	return s
}

// feed gives the messages to the player in order.
func feed(p *Player, msgs ...[2]string) {
	for _, m := range msgs {
		p.bandMessage(m[0], m[1])
	}
}

func msg(prefix, text string) [2]string { return [2]string{prefix, text} }

// allFloor reports whether every band is at -60.
func allFloor(b [SpectrumBands]float64) bool {
	for _, v := range b {
		if v != -60 {
			return false
		}
	}
	return true
}

// START: TestBandMessagePairSetsBand

func TestBandMessagePairSetsBand(t *testing.T) {
	p := &Player{}
	feed(p, msg(level(4, "-21.9")), msg(band(4, 0)))
	got := p.Bands()
	if got[0] != -21.9 {
		t.Fatalf("band 0 = %v, want -21.9", got[0])
	}
	got[0] = -60
	if !allFloor(got) {
		t.Fatalf("the other bands are %v, want all -60", got)
	}
}

// END: TestBandMessagePairSetsBand

// START: TestBandMessagesPairByInstance

func TestBandMessagesPairByInstance(t *testing.T) {
	p := &Player{}
	feed(p, msg(level(4, "-10")), msg(level(9, "-40")), msg(band(9, 1)), msg(band(4, 0)))
	got := p.Bands()
	if got[0] != -10 || got[1] != -40 {
		t.Fatalf("bands 0 and 1 = %v and %v, want -10 and -40", got[0], got[1])
	}
}

// END: TestBandMessagesPairByInstance

// START: TestBandLineWithoutLevelIsIgnored

func TestBandLineWithoutLevelIsIgnored(t *testing.T) {
	p := &Player{}
	feed(p, msg(band(4, 3)))
	if got := p.Bands(); !allFloor(got) {
		t.Fatalf("bands = %v, want all -60", got)
	}
}

// END: TestBandLineWithoutLevelIsIgnored

// START: TestBandMessageOtherPrefixIsIgnored

func TestBandMessageOtherPrefixIsIgnored(t *testing.T) {
	p := &Player{}
	_, l := level(4, "-10")
	_, b := band(4, 0)
	feed(p, msg("cplayer", l), msg("cplayer", b))
	if got := p.Bands(); !allFloor(got) {
		t.Fatalf("bands = %v, want all -60", got)
	}
}

// END: TestBandMessageOtherPrefixIsIgnored

// START: TestBandsReadWhileWritten

func TestBandsReadWhileWritten(t *testing.T) {
	p := &Player{}
	start := make(chan struct{})
	var wg sync.WaitGroup
	bad := make(chan float64, 1)
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 1000; i++ {
			db := "-20"
			if i%2 == 1 {
				db = "-30"
			}
			feed(p, msg(level(4, db)), msg(band(4, 5)))
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 5000; i++ {
			if v := p.Bands()[5]; v != -60 && v != -20 && v != -30 {
				select {
				case bad <- v:
				default:
				}
			}
		}
	}()
	close(start)
	wg.Wait()
	select {
	case v := <-bad:
		t.Fatalf("a reading of band 5 was %v, want -60, -20 or -30", v)
	default:
	}
	if got := p.Bands()[5]; got != -30 {
		t.Fatalf("band 5 after the writer ended = %v, want the last level written, -30", got)
	}
}

// END: TestBandsReadWhileWritten

// START: TestLogMessageBandLineAtVWritesNothing

func TestLogMessageBandLineAtVWritesNothing(t *testing.T) {
	r := &recorder{}
	(&Player{log: r.logf}).logMessage("v", "ffmpeg", "Parsed_ametadata_4: lavfi.band=0\n")
	if got := r.get(); len(got) != 0 {
		t.Fatalf("log = %q, want nothing", got)
	}
}

// END: TestLogMessageBandLineAtVWritesNothing

// START: TestLogMessageOtherVMessageWritesNothing

func TestLogMessageOtherVMessageWritesNothing(t *testing.T) {
	r := &recorder{}
	(&Player{log: r.logf}).logMessage("v", "cplayer", "Running hook: x\n")
	if got := r.get(); len(got) != 0 {
		t.Fatalf("log = %q, want nothing", got)
	}
}

// END: TestLogMessageOtherVMessageWritesNothing

// START: TestLogMessageWarnStillLoggedWithSpectrum

func TestLogMessageWarnStillLoggedWithSpectrum(t *testing.T) {
	r := &recorder{}
	p := &Player{log: r.logf}
	feed(p, msg(level(4, "-10")), msg(band(4, 0)))
	p.logMessage("warn", "ffmpeg", "https: HTTP error 403 Forbidden\n")
	if got := r.get(); len(got) != 1 || got[0] != "mpv [ffmpeg] https: HTTP error 403 Forbidden" {
		t.Fatalf("log = %q", got)
	}
}

// END: TestLogMessageWarnStillLoggedWithSpectrum
