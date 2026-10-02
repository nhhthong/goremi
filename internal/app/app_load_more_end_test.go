// Test for the list of a first page that is not full: no load more line.
package app

import (
	"fmt"
	"strings"
	"testing"

	"goremi/internal/provider"
)

// START: TestViewNoLoadMoreUnderTen

func TestViewNoLoadMoreUnderTen(t *testing.T) {
	var nine []provider.Track
	for i := 0; i < 9; i++ {
		nine = append(nine, provider.Track{ID: fmt.Sprint(i), Title: fmt.Sprintf("T%d", i)})
	}
	p := &scriptedProvider{steps: []step{{tracks: nine}}}
	if got := search(New(p), "daft").View().Content; strings.Contains(got, "load more...") {
		t.Fatalf("view %q has a load more line after a page of 9 tracks", got)
	}
}

// END: TestViewNoLoadMoreUnderTen
