// Tests for the levels and band numbers mpv sends that cannot be read: they count as silence or are ignored, and nothing panics.
package player

import "testing"

// setBand1 gives band 0 the level -10 so a later unreadable level has something to replace.
func setBand0(p *Player) { feed(p, msg(level(4, "-10")), msg(band(4, 0))) }

// START: TestBandLevelMinusInfIsFloor

func TestBandLevelMinusInfIsFloor(t *testing.T) {
	p := &Player{}
	setBand0(p)
	feed(p, msg(level(4, "-inf")), msg(band(4, 0)))
	if got := p.Bands()[0]; got != -60 {
		t.Fatalf("band 0 after -inf = %v, want -60", got)
	}
}

// END: TestBandLevelMinusInfIsFloor

// START: TestBandLevelUnreadableIsFloor

func TestBandLevelUnreadableIsFloor(t *testing.T) {
	for _, text := range []string{"abc", "nan", "--5", "1e999x"} {
		p := &Player{}
		setBand0(p)
		feed(p, msg(level(4, text)), msg(band(4, 0)))
		if got := p.Bands()[0]; got != -60 {
			t.Fatalf("band 0 after level %q = %v, want -60", text, got)
		}
	}
}

// END: TestBandLevelUnreadableIsFloor

// START: TestBandOutsideRangeIsIgnored

func TestBandOutsideRangeIsIgnored(t *testing.T) {
	p := &Player{}
	setBand0(p)
	feed(p, msg(level(4, "-5")), msg(band(4, 32)), msg(band(4, 99)))
	got := p.Bands()
	if got[0] != -10 {
		t.Fatalf("band 0 = %v, want -10 untouched", got[0])
	}
	got[0] = -60
	if !allFloor(got) {
		t.Fatalf("bands = %v, want the others at -60", got)
	}
}

// END: TestBandOutsideRangeIsIgnored

// START: TestBandGarbageDoesNotPanic

func TestBandGarbageDoesNotPanic(t *testing.T) {
	p := &Player{}
	setBand0(p)
	for _, text := range []string{
		"Parsed_ametadata_4: lavfi.band=99999999999999999999\n",
		"Parsed_ametadata_4: lavfi.band=\n",
		"Parsed_ametadata_4: lavfi.astats.Overall.RMS_level=\n",
		"\x00\xff\xfe",
		"",
	} {
		p.bandMessage("ffmpeg", text)
		p.bandMessage("", text)
	}
	if got := p.Bands()[0]; got != -10 {
		t.Fatalf("band 0 = %v, want -10 untouched", got)
	}
}

// END: TestBandGarbageDoesNotPanic
