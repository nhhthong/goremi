// Theme definition: the semantic colour fields every component reads (palettes live in dark.go, light.go, cyberpunk.go).
package theme

// START: Theme

// Theme holds the colours of one theme as hex strings; components use these fields, never literal colours.
type Theme struct {
	Background string
	Foreground string
	Muted      string
	Border     string
	Accent     string
	Selected   string
	Progress   string
	Spectrum   string
	LogoFrom   string
	LogoTo     string
}

// END: Theme

// START: logo stops

// The logo gradient runs teal to blue; every theme uses the same pair.
const (
	logoFrom = "#2dd4bf"
	logoTo   = "#3b82f6"
)

// END: logo stops

// START: registry

// Entry is a theme with the name the user picks it by.
type Entry struct {
	Name  string
	Theme Theme
}

// all is the one list of themes; the picker and ByName read it. A new theme is one entry here.
var all = []Entry{
	{"light", Light()},
	{"dark", Dark()},
	{"cyberpunk", Cyberpunk()},
}

// All returns the themes in display order.
func All() []Entry { return all }

// ByName returns the theme called name; ok is false when no theme has that name.
func ByName(name string) (t Theme, ok bool) {
	for _, e := range all {
		if e.Name == name {
			return e.Theme, true
		}
	}
	return Theme{}, false
}

// END: registry
