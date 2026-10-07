// Tests of the Error colour every theme carries (task 9.1.1).
package theme

import (
	"reflect"
	"testing"
)

// START: TestThemeHasErrorField

func TestThemeHasErrorField(t *testing.T) {
	f, ok := reflect.TypeOf(Theme{}).FieldByName("Error")
	if !ok || f.Type.Kind() != reflect.String {
		t.Fatalf("Theme has no string field Error (found %v)", ok)
	}
}

// END: TestThemeHasErrorField

// START: TestEveryThemeSetsError

func TestEveryThemeSetsError(t *testing.T) {
	want := map[string]string{
		"default":    "#ff5555",
		"catppuccin": "#f38ba8",
		"dracula":    "#ff5555",
		"gruvbox":    "#fb4934",
		"nord":       "#bf616a",
		"rosepine":   "#eb6f92",
		"tokyonight": "#f7768e",
	}
	for _, e := range All() {
		got := reflect.ValueOf(e.Theme).FieldByName("Error")
		if !got.IsValid() || got.String() != want[e.Name] {
			t.Errorf("theme %s: Error = %q, want %q", e.Name, got, want[e.Name])
		}
	}
}

// END: TestEveryThemeSetsError
