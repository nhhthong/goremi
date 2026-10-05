// Tests for the playback keys while the list has focus: pause and seek back.
package app

import (
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// pressInList presses one key with the list focused and returns the fake player.
func pressInList(k tea.KeyPressMsg) *recordingPlayer {
	pl := &recordingPlayer{}
	New(&resolveProvider{}).WithPlayer(pl).WithFocus(FocusList).Update(k)
	return pl
}

// START: TestKeyKTogglesPause

func TestKeyKTogglesPause(t *testing.T) {
	if pl := pressInList(tea.KeyPressMsg{Code: 'k', Text: "k"}); pl.toggles != 1 {
		t.Fatalf("TogglePause calls after k = %d, want 1", pl.toggles)
	}
}

// END: TestKeyKTogglesPause

// START: TestSpaceTogglesPause

func TestSpaceTogglesPause(t *testing.T) {
	if pl := pressInList(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "}); pl.toggles != 1 {
		t.Fatalf("TogglePause calls after Space = %d, want 1", pl.toggles)
	}
}

// END: TestSpaceTogglesPause

// START: TestKeyJSeeksBack

func TestKeyJSeeksBack(t *testing.T) {
	if pl := pressInList(tea.KeyPressMsg{Code: 'j', Text: "j"}); !reflect.DeepEqual(pl.seeks, []float64{-10}) {
		t.Fatalf("Seek calls after j = %v, want [-10]", pl.seeks)
	}
}

// END: TestKeyJSeeksBack

// START: TestLeftSeeksBack

func TestLeftSeeksBack(t *testing.T) {
	if pl := pressInList(tea.KeyPressMsg{Code: tea.KeyLeft}); !reflect.DeepEqual(pl.seeks, []float64{-10}) {
		t.Fatalf("Seek calls after ← = %v, want [-10]", pl.seeks)
	}
}

// END: TestLeftSeeksBack

// START: TestKeyLSeeksForward

func TestKeyLSeeksForward(t *testing.T) {
	if pl := pressInList(tea.KeyPressMsg{Code: 'l', Text: "l"}); !reflect.DeepEqual(pl.seeks, []float64{10}) {
		t.Fatalf("Seek calls after l = %v, want [10]", pl.seeks)
	}
}

// END: TestKeyLSeeksForward

// START: TestRightSeeksForward

func TestRightSeeksForward(t *testing.T) {
	if pl := pressInList(tea.KeyPressMsg{Code: tea.KeyRight}); !reflect.DeepEqual(pl.seeks, []float64{10}) {
		t.Fatalf("Seek calls after → = %v, want [10]", pl.seeks)
	}
}

// END: TestRightSeeksForward
