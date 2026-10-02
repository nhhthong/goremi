// Tests for the colours of the message line and of the whole screen.
package app

import (
	"image/color"
	"testing"

	"goremi/internal/provider"
	"goremi/internal/ui/theme"
)

// failedSearchLine returns the raw message line under Search: after a network failure, drawn with t.
func failedSearchLine(t *testing.T, th theme.Theme) string {
	t.Helper()
	p := &scriptedProvider{steps: []step{{err: provider.ErrNetwork}}}
	return rowLine(t, search(New(p).WithTheme(th), "daft"), networkText)
}

// rgb returns the 8-bit red, green and blue of c.
func rgb(c color.Color) [3]uint8 {
	r, g, b, _ := c.RGBA()
	return [3]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

// START: TestNoticeForegroundDark

func TestNoticeForegroundDark(t *testing.T) {
	wantCodes(t, failedSearchLine(t, theme.Dark()), "38;2;205;214;244")
}

// END: TestNoticeForegroundDark

// START: TestNoticeForegroundLight

func TestNoticeForegroundLight(t *testing.T) {
	wantCodes(t, failedSearchLine(t, theme.Light()), "38;2;43;43;43")
}

// END: TestNoticeForegroundLight

// START: TestViewScreenColoursDark

func TestViewScreenColoursDark(t *testing.T) {
	v := New(fakeProvider{}).View()
	if bg, fg := rgb(v.BackgroundColor), rgb(v.ForegroundColor); bg != [3]uint8{30, 30, 46} || fg != [3]uint8{205, 214, 244} {
		t.Fatalf("background %v, foreground %v; want [30 30 46] and [205 214 244]", bg, fg)
	}
}

// END: TestViewScreenColoursDark

// START: TestViewScreenColoursLight

func TestViewScreenColoursLight(t *testing.T) {
	v := New(fakeProvider{}).WithTheme(theme.Light()).View()
	if bg, fg := rgb(v.BackgroundColor), rgb(v.ForegroundColor); bg != [3]uint8{250, 250, 250} || fg != [3]uint8{43, 43, 43} {
		t.Fatalf("background %v, foreground %v; want [250 250 250] and [43 43 43]", bg, fg)
	}
}

// END: TestViewScreenColoursLight
