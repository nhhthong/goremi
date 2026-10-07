// Tests for a search that returns no track: the focus stays on the input.
package app

import "testing"

// zeroResults is the model after the query "daft" returned no track.
func zeroResults() Model {
	return search(New(&scriptedProvider{steps: []step{{}}}), "daft")
}

// START: TestZeroResultsKeepFocusOnInput

func TestZeroResultsKeepFocusOnInput(t *testing.T) {
	if got := zeroResults().Focus(); got != FocusInput {
		t.Fatalf("Focus() = %v, want FocusInput", got)
	}
}

// END: TestZeroResultsKeepFocusOnInput

// START: TestZeroResultsMessage

func TestZeroResultsMessage(t *testing.T) {
	if got, want := lineAboveSearch(viewLines(zeroResults())), `No results for "daft".`; got != want {
		t.Fatalf("line under Search: = %q, want %q", got, want)
	}
}

// END: TestZeroResultsMessage

// START: TestZeroResultsRetypeAtOnce

func TestZeroResultsRetypeAtOnce(t *testing.T) {
	next, _ := zeroResults().Update(key('x'))
	if got := next.(Model).Query(); got != "daftx" {
		t.Fatalf("Query() = %q, want %q (the input has the focus, so a key edits the query)", got, "daftx")
	}
}

// END: TestZeroResultsRetypeAtOnce
