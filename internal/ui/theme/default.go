// The default palette.
package theme

// START: Default

// Default returns the goremi's own theme: teal to blue, the colours of its logo.
func Default() Theme {
	return Theme{
		Accent:   "#2dd4bf",
		Error:    "#ff5555",
		Selected: "#134e4a",
		Muted:    "#8a8a8a",
		Border:   "#6b7280",
		Progress: "#3b82f6",
		Spectrum: "#2dd4bf",
		LogoFrom: "#2dd4bf",
		LogoTo:   "#3b82f6",
	}
}

// END: Default
