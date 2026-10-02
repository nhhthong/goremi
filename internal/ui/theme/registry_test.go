// Tests for the theme registry: the list of themes and the lookup by name.
package theme

import (
	"reflect"
	"testing"
)

// START: TestAllNames

func TestAllNames(t *testing.T) {
	var got []string
	for _, e := range All() {
		got = append(got, e.Name)
	}
	if want := []string{"light", "dark", "cyberpunk"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("All() names = %v, want %v", got, want)
	}
}

// END: TestAllNames

// START: TestByNameListed

func TestByNameListed(t *testing.T) {
	want := map[string]Theme{"light": Light(), "dark": Dark(), "cyberpunk": Cyberpunk()}
	for name, th := range want {
		got, ok := ByName(name)
		if !ok || got != th {
			t.Errorf("ByName(%q) = %v, %v; want the %s theme, true", name, got, ok, name)
		}
	}
}

// END: TestByNameListed

// START: TestByNameUnlisted

func TestByNameUnlisted(t *testing.T) {
	for _, name := range []string{"neon", "Dark", ""} {
		if _, ok := ByName(name); ok {
			t.Errorf("ByName(%q) found a theme, want not found", name)
		}
	}
}

// END: TestByNameUnlisted
