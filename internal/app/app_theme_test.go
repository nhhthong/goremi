// Tests for the theme the model holds.
package app

import (
	"testing"

	"goremi/internal/ui/theme"
)

// START: TestDefaultThemeIsDefault

func TestDefaultThemeIsDefault(t *testing.T) {
	if got := New(fakeProvider{}).Theme(); got != theme.Default() {
		t.Fatalf("Theme() = %v, want the default theme", got)
	}
}

// END: TestDefaultThemeIsDefault

// START: TestWithTheme

func TestWithTheme(t *testing.T) {
	if got := New(fakeProvider{}).WithTheme(theme.Nord()).Theme(); got != theme.Nord() {
		t.Fatalf("Theme() = %v, want the nord theme", got)
	}
}

// END: TestWithTheme
