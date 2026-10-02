// Tests for the theme the model holds.
package app

import (
	"testing"

	"goremi/internal/ui/theme"
)

// START: TestDefaultThemeIsDark

func TestDefaultThemeIsDark(t *testing.T) {
	if got := New(fakeProvider{}).Theme(); got != theme.Dark() {
		t.Fatalf("Theme() = %v, want the dark theme", got)
	}
}

// END: TestDefaultThemeIsDark

// START: TestWithTheme

func TestWithTheme(t *testing.T) {
	if got := New(fakeProvider{}).WithTheme(theme.Light()).Theme(); got != theme.Light() {
		t.Fatalf("Theme() = %v, want the light theme", got)
	}
}

// END: TestWithTheme
