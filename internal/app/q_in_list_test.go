// Tests that `q` in the results list no longer quits (task 10.5.1).
package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: listModel

// listModel returns a model with two tracks, the focus on the list and the first track selected.
func listModel() Model {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}, {Title: "B"}}}}}
	return search(New(p), "daft").WithFocus(FocusList)
}

// END: listModel

// START: TestQInListChangesNothing

func TestQInListChangesNothing(t *testing.T) {
	before := listModel()
	next, _ := before.Update(key('q'))
	m := next.(Model)
	if m.Focus() != FocusList || m.Selected() != before.Selected() || m.Query() != before.Query() {
		t.Fatalf("focus %v, selected %d, query %q; want focus on the list, selected %d, query %q",
			m.Focus(), m.Selected(), m.Query(), before.Selected(), before.Query())
	}
}

// END: TestQInListChangesNothing

// START: TestQInListDoesNotQuit

func TestQInListDoesNotQuit(t *testing.T) {
	if _, cmd := listModel().Update(key('q')); quits(cmd) {
		t.Fatal("q in the list returned a quit command, want none")
	}
}

// END: TestQInListDoesNotQuit

// START: TestCtrlCStillQuitsFromList

func TestCtrlCStillQuitsFromList(t *testing.T) {
	if _, cmd := listModel().Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}); !quits(cmd) {
		t.Fatal("Ctrl+C in the list returned no quit command")
	}
}

// END: TestCtrlCStillQuitsFromList
