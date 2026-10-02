// Tests for the palette of each theme.
package theme

import "testing"

// START: TestDefaultPalette

func TestDefaultPalette(t *testing.T) {
	checkPalette(t, Default(), map[string]string{
		"Accent":   "#2dd4bf",
		"Selected": "#134e4a",
		"Muted":    "#8a8a8a",
		"Border":   "#6b7280",
		"Progress": "#3b82f6",
		"Spectrum": "#2dd4bf",
		"LogoFrom": "#2dd4bf",
		"LogoTo":   "#3b82f6",
	})
}

// END: TestDefaultPalette

// START: TestCatppuccinPalette

func TestCatppuccinPalette(t *testing.T) {
	checkPalette(t, Catppuccin(), map[string]string{
		"Accent":   "#cba6f7",
		"Selected": "#313244",
		"Muted":    "#6c7086",
		"Border":   "#6c7086",
		"Progress": "#f5c2e7",
		"Spectrum": "#89dceb",
		"LogoFrom": "#cba6f7",
		"LogoTo":   "#89b4fa",
	})
}

// END: TestCatppuccinPalette

// START: TestDraculaPalette

func TestDraculaPalette(t *testing.T) {
	checkPalette(t, Dracula(), map[string]string{
		"Accent":   "#bd93f9",
		"Selected": "#44475a",
		"Muted":    "#6272a4",
		"Border":   "#6272a4",
		"Progress": "#ff79c6",
		"Spectrum": "#8be9fd",
		"LogoFrom": "#bd93f9",
		"LogoTo":   "#ff79c6",
	})
}

// END: TestDraculaPalette

// START: TestGruvboxPalette

func TestGruvboxPalette(t *testing.T) {
	checkPalette(t, Gruvbox(), map[string]string{
		"Accent":   "#fabd2f",
		"Selected": "#504945",
		"Muted":    "#928374",
		"Border":   "#928374",
		"Progress": "#b8bb26",
		"Spectrum": "#fe8019",
		"LogoFrom": "#fabd2f",
		"LogoTo":   "#fe8019",
	})
}

// END: TestGruvboxPalette

// START: TestNordPalette

func TestNordPalette(t *testing.T) {
	checkPalette(t, Nord(), map[string]string{
		"Accent":   "#88c0d0",
		"Selected": "#434c5e",
		"Muted":    "#4c566a",
		"Border":   "#4c566a",
		"Progress": "#81a1c1",
		"Spectrum": "#8fbcbb",
		"LogoFrom": "#8fbcbb",
		"LogoTo":   "#5e81ac",
	})
}

// END: TestNordPalette

// START: TestRosePinePalette

func TestRosePinePalette(t *testing.T) {
	checkPalette(t, RosePine(), map[string]string{
		"Accent":   "#c4a7e7",
		"Selected": "#403d52",
		"Muted":    "#6e6a86",
		"Border":   "#6e6a86",
		"Progress": "#9ccfd8",
		"Spectrum": "#ebbcba",
		"LogoFrom": "#c4a7e7",
		"LogoTo":   "#ebbcba",
	})
}

// END: TestRosePinePalette

// START: TestTokyoNightPalette

func TestTokyoNightPalette(t *testing.T) {
	checkPalette(t, TokyoNight(), map[string]string{
		"Accent":   "#7aa2f7",
		"Selected": "#292e42",
		"Muted":    "#565f89",
		"Border":   "#565f89",
		"Progress": "#bb9af7",
		"Spectrum": "#7dcfff",
		"LogoFrom": "#7aa2f7",
		"LogoTo":   "#bb9af7",
	})
}

// END: TestTokyoNightPalette
