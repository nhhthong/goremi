// Tests for the theme colours the view uses.
package app

import (
	"strings"
	"testing"

	"goremi/internal/ui/theme"
)

// wantCodes fails when text lacks one of the colour codes.
func wantCodes(t *testing.T, text string, codes ...string) {
	t.Helper()
	for _, c := range codes {
		if !strings.Contains(text, c) {
			t.Errorf("%q lacks the colour code %s", text, c)
		}
	}
}

// START: TestPanelUsesModelLogoStops

func TestPanelUsesModelLogoStops(t *testing.T) {
	m := sized(New(fakeProvider{}).WithTheme(theme.Theme{LogoFrom: "#ff0000", LogoTo: "#0000ff"}), 100)
	wantCodes(t, m.View().Content, "38;2;255;0;0", "38;2;0;0;255")
}

// END: TestPanelUsesModelLogoStops

// START: TestPanelDefaultLogoColours

func TestPanelDefaultLogoColours(t *testing.T) {
	wantCodes(t, sized(New(fakeProvider{}), 100).View().Content, "38;2;45;212;191", "38;2;59;130;246")
}

// END: TestPanelDefaultLogoColours

// START: TestHintUsesMutedDark

func TestHintUsesMutedDark(t *testing.T) {
	first := viewLines(New(fakeProvider{}))[0]
	wantCodes(t, first, "38;2;108;112;134", "Ctrl+C: quit", "goremi theme")
}

// END: TestHintUsesMutedDark

// START: TestHintUsesMutedLight

func TestHintUsesMutedLight(t *testing.T) {
	wantCodes(t, viewLines(New(fakeProvider{}).WithTheme(theme.Light()))[0], "38;2;138;138;138")
}

// END: TestHintUsesMutedLight

// START: TestSearchLabelAccentDark

func TestSearchLabelAccentDark(t *testing.T) {
	wantCodes(t, viewLines(New(fakeProvider{}))[1], "38;2;137;180;250")
}

// END: TestSearchLabelAccentDark

// START: TestSearchLabelAccentLight

func TestSearchLabelAccentLight(t *testing.T) {
	wantCodes(t, viewLines(New(fakeProvider{}).WithTheme(theme.Light()))[1], "38;2;74;111;165")
}

// END: TestSearchLabelAccentLight
