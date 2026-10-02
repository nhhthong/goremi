// The dark palette.
package theme

// START: Dark

// Dark returns the dark theme: minimal, terminal-native.
func Dark() Theme {
	return Theme{
		Background: "#1e1e2e",
		Foreground: "#cdd6f4",
		Muted:      "#6c7086",
		Border:     "#45475a",
		Accent:     "#89b4fa",
		Selected:   "#313244",
		Progress:   "#89b4fa",
		Spectrum:   "#74c7ec",
		LogoFrom:   logoFrom,
		LogoTo:     logoTo,
	}
}

// END: Dark
