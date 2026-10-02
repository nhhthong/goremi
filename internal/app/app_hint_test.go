// Test for the start-up hint that names the theme command.
package app

import (
	"strings"
	"testing"
)

// START: TestViewShowsThemeHint

func TestViewShowsThemeHint(t *testing.T) {
	if got := New(fakeProvider{}).View().Content; !strings.Contains(got, "goremi theme") {
		t.Fatalf("view %q: want it to contain the hint %q", got, "goremi theme")
	}
}

// END: TestViewShowsThemeHint
