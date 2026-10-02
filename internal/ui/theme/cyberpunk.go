// The cyberpunk palette.
package theme

// START: Cyberpunk

// Cyberpunk returns the cyberpunk theme: very dark background, neon cyan and magenta accents.
func Cyberpunk() Theme {
	return Theme{
		Background: "#0a0a14",
		Foreground: "#e0e0ff",
		Muted:      "#5a5a7a",
		Border:     "#ff2bd6",
		Accent:     "#00f0ff",
		Selected:   "#1a1a3a",
		Progress:   "#00f0ff",
		Spectrum:   "#ff2bd6",
		LogoFrom:   logoFrom,
		LogoTo:     logoTo,
	}
}

// END: Cyberpunk
