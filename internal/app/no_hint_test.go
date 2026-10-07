// Tests that the old hint line `Ctrl+C: quit · run goremi theme ...` is gone from the view (task 3.2.3).
package app

import (
	"strings"
	"testing"
)

// START: TestNoQuitHint

func TestNoQuitHint(t *testing.T) {
	if got := plain(sizeTo(New(fakeProvider{}), 100, 30).View().Content); strings.Contains(got, "Ctrl+C: quit") {
		t.Fatalf("view %q: want no Ctrl+C hint", got)
	}
}

// END: TestNoQuitHint

// START: TestNoThemeHint

func TestNoThemeHint(t *testing.T) {
	for _, width := range []int{0, 40, 100} {
		if got := plain(sizeTo(New(fakeProvider{}), width, 30).View().Content); strings.Contains(got, "goremi theme") {
			t.Fatalf("width %d, view %q: want no `goremi theme` text", width, got)
		}
	}
}

// END: TestNoThemeHint

// START: TestHintLineIsGone

func TestHintLineIsGone(t *testing.T) {
	const old = "Ctrl+C: quit · run goremi theme to choose a theme"
	for _, width := range []int{0, 40, 80, 100} {
		if got := plain(sizeTo(New(fakeProvider{}), width, 30).View().Content); strings.Contains(got, old) {
			t.Fatalf("width %d, view %q: want the old hint line gone", width, got)
		}
	}
}

// END: TestHintLineIsGone
