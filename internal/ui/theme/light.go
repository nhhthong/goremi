// The light palette.
package theme

// START: Light

// Light returns the light theme: clean, minimal, bright.
func Light() Theme {
	return Theme{
		Background: "#fafafa",
		Foreground: "#2b2b2b",
		Muted:      "#8a8a8a",
		Border:     "#d0d0d0",
		Accent:     "#4a6fa5",
		Selected:   "#e6e6e6",
		Progress:   "#4a6fa5",
		Spectrum:   "#6a8caf",
		LogoFrom:   logoFrom,
		LogoTo:     logoTo,
	}
}

// END: Light
