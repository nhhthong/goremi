// Tests that ↑ in the search bar returns to the results list (ui task 10.16).
package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: TestUpFromBarReturnsToList

func TestUpFromBarReturnsToList(t *testing.T) {
	var tracks []provider.Track
	for i := 0; i < provider.PageSize; i++ {
		tracks = append(tracks, provider.Track{ID: string(rune('a' + i)), Title: string(rune('A' + i))})
	}
	m := search(sizeTo(New(&scriptedProvider{steps: []step{{tracks: tracks}}}), 100, 30), "daft")
	for i := 0; i <= provider.PageSize; i++ { // the last Down jumps from the load more line to the bar
		next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		m = next.(Model)
	}
	if m.Focus() != FocusInput {
		t.Fatalf("focus = %v, want the bar after the jump", m.Focus())
	}
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if got := next.(Model); got.Focus() != FocusList || got.Selected() != provider.PageSize {
		t.Fatalf("focus %v, selected %d; want the list with the load more line selected (%d)", got.Focus(), got.Selected(), provider.PageSize)
	}
}

// END: TestUpFromBarReturnsToList

// START: TestUpInBarWithoutResultsDoesNothing

func TestUpInBarWithoutResultsDoesNothing(t *testing.T) {
	next, _ := New(fakeProvider{}).Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if got := next.(Model).Focus(); got != FocusInput {
		t.Fatalf("focus = %v with no results, want the bar", got)
	}
}

// END: TestUpInBarWithoutResultsDoesNothing

// START: TestDownInBarDoesNothing

func TestDownInBarDoesNothing(t *testing.T) {
	m := listAt(100, 30).WithFocus(FocusInput)
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if got := next.(Model); got.Focus() != FocusInput || got.Selected() != m.Selected() {
		t.Fatalf("focus %v, selected %d; want the bar and the same line", got.Focus(), got.Selected())
	}
}

// END: TestDownInBarDoesNothing
