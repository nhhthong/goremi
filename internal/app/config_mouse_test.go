// Tests for the [ui] mouse key of the config file.
package app

import "testing"

// START: TestLoadConfigMouseFalse

func TestLoadConfigMouseFalse(t *testing.T) {
	if LoadConfig(writeConfig(t, "[ui]\nmouse = false\n")).Mouse {
		t.Fatal("Mouse = true, want false")
	}
}

// END: TestLoadConfigMouseFalse

// START: TestLoadConfigMouseDefaultTrue

func TestLoadConfigMouseDefaultTrue(t *testing.T) {
	if !LoadConfig(writeConfig(t, "[ui]\ntheme = \"default\"\n")).Mouse {
		t.Fatal("Mouse = false, want true when the key is absent")
	}
}

// END: TestLoadConfigMouseDefaultTrue

// START: TestSaveThemeKeepsMouse

func TestSaveThemeKeepsMouse(t *testing.T) {
	path := writeConfig(t, "[ui]\nmouse = false\n")
	if err := SaveTheme(path, "nord"); err != nil {
		t.Fatal(err)
	}
	c := LoadConfig(path)
	if c.Mouse || c.Theme != "nord" {
		t.Fatalf("after SaveTheme: Mouse = %v, Theme = %q, want false, nord", c.Mouse, c.Theme)
	}
}

// END: TestSaveThemeKeepsMouse
