// Tests for the Theme struct and the dark palette, with the palette check the other theme tests share.
package theme

import (
	"reflect"
	"testing"
)

// START: checkPalette

// checkPalette fails for each named field of got that is not the wanted hex value.
func checkPalette(t *testing.T, got Theme, want map[string]string) {
	t.Helper()
	v := reflect.ValueOf(got)
	for name, hex := range want {
		if g := v.FieldByName(name).String(); g != hex {
			t.Errorf("%s = %q, want %q", name, g, hex)
		}
	}
}

// END: checkPalette

// START: TestDarkPalette

func TestDarkPalette(t *testing.T) {
	checkPalette(t, Dark(), map[string]string{
		"Background": "#1e1e2e",
		"Foreground": "#cdd6f4",
		"Muted":      "#6c7086",
		"Border":     "#45475a",
		"Accent":     "#89b4fa",
		"Selected":   "#313244",
		"Progress":   "#89b4fa",
		"Spectrum":   "#74c7ec",
	})
}

// END: TestDarkPalette

// START: TestThemeFields

func TestThemeFields(t *testing.T) {
	want := []string{"Background", "Foreground", "Muted", "Border", "Accent", "Selected", "Progress", "Spectrum", "LogoFrom", "LogoTo"}
	typ := reflect.TypeOf(Theme{})
	if typ.NumField() != len(want) {
		t.Fatalf("Theme has %d fields, want %d", typ.NumField(), len(want))
	}
	for i, name := range want {
		f := typ.Field(i)
		if f.Name != name || f.Type.Kind() != reflect.String {
			t.Errorf("field %d = %s %s, want %s string", i, f.Name, f.Type, name)
		}
	}
}

// END: TestThemeFields

// START: TestLogoStops

func TestLogoStops(t *testing.T) {
	want := map[string]string{"LogoFrom": "#2dd4bf", "LogoTo": "#3b82f6"}
	for name, th := range map[string]Theme{"dark": Dark(), "light": Light(), "cyberpunk": Cyberpunk()} {
		t.Run(name, func(t *testing.T) { checkPalette(t, th, want) })
	}
}

// END: TestLogoStops
