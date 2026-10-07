// Tests for the Theme struct, with the palette check the theme tests share.
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

// START: TestThemeFields

func TestThemeFields(t *testing.T) {
	want := []string{"Background", "Foreground", "Muted", "Border", "Accent", "Error", "Selected", "Progress", "Spectrum", "LogoFrom", "LogoTo"}
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
