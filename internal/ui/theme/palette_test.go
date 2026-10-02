// Tests for the light and cyberpunk palettes.
package theme

import "testing"

// START: TestLightPalette

func TestLightPalette(t *testing.T) {
	checkPalette(t, Light(), map[string]string{
		"Background": "#fafafa",
		"Foreground": "#2b2b2b",
		"Muted":      "#8a8a8a",
		"Border":     "#d0d0d0",
		"Accent":     "#4a6fa5",
		"Selected":   "#e6e6e6",
		"Progress":   "#4a6fa5",
		"Spectrum":   "#6a8caf",
	})
}

// END: TestLightPalette

// START: TestCyberpunkPalette

func TestCyberpunkPalette(t *testing.T) {
	checkPalette(t, Cyberpunk(), map[string]string{
		"Background": "#0a0a14",
		"Foreground": "#e0e0ff",
		"Muted":      "#5a5a7a",
		"Border":     "#ff2bd6",
		"Accent":     "#00f0ff",
		"Selected":   "#1a1a3a",
		"Progress":   "#00f0ff",
		"Spectrum":   "#ff2bd6",
	})
}

// END: TestCyberpunkPalette
