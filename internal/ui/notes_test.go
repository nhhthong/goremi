// Tests of the generated music notes of the mascot: a clean art, the free places, the overlay and the steps (ui tasks 3.11.1 to 3.11.4).
package ui

import (
	"strings"
	"testing"
	"unicode/utf8"

	"goremi/internal/ui/theme"
)

// blank is the blank braille cell.
const blank = '⠀'

// grid is the mascot art as a grid of runes.
func grid() [][]rune {
	var g [][]rune
	for _, row := range strings.Split(Mascot(), "\n") {
		g = append(g, []rune(row))
	}
	return g
}

// START: TestMascotHasNoBakedNotes

func TestMascotHasNoBakedNotes(t *testing.T) {
	g := grid()
	seen := map[[2]int]bool{}
	parts := 0
	for y := range g {
		for x := range g[y] {
			if g[y][x] == blank || seen[[2]int{y, x}] {
				continue
			}
			parts++
			stack, onGround := [][2]int{{y, x}}, false
			seen[[2]int{y, x}] = true
			for len(stack) > 0 {
				c := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onGround = onGround || c[0] == len(g)-1
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						ny, nx := c[0]+dy, c[1]+dx
						if ny >= 0 && ny < len(g) && nx >= 0 && nx < len(g[ny]) && g[ny][nx] != blank && !seen[[2]int{ny, nx}] {
							seen[[2]int{ny, nx}] = true
							stack = append(stack, [2]int{ny, nx})
						}
					}
				}
			}
			if parts > 1 && !onGround {
				t.Errorf("a loose part starts at row %d column %d: a note or a sparkle is still in the art", y, x)
			}
		}
	}
}

// END: TestMascotHasNoBakedNotes

// START: TestMascotStillTwelveByThirtyTwo

func TestMascotStillTwelveByThirtyTwo(t *testing.T) {
	rows := strings.Split(Mascot(), "\n")
	if len(rows) != 12 {
		t.Fatalf("%d rows, want 12", len(rows))
	}
	for i, r := range rows {
		if utf8.RuneCountInString(r) != 32 {
			t.Errorf("row %d is %d columns, want 32", i, utf8.RuneCountInString(r))
		}
	}
}

// END: TestMascotStillTwelveByThirtyTwo

// START: TestNoteSpotsAreFree

func TestNoteSpotsAreFree(t *testing.T) {
	spots := NoteSpots()
	if len(spots) < 10 {
		t.Fatalf("%d spots, want at least 10", len(spots))
	}
	g := grid()
	for _, s := range spots {
		for _, c := range [][2]int{{s.Row, s.Col + 1}, {s.Row, s.Col + 2}, {s.Row + 1, s.Col}, {s.Row + 1, s.Col + 1}} {
			if c[0] < 0 || c[0] >= 12 || c[1] < 0 || c[1] >= 32 || g[c[0]][c[1]] != blank {
				t.Fatalf("spot %+v: cell %v is outside the area or holds art", s, c)
			}
		}
	}
}

// END: TestNoteSpotsAreFree

// START: TestOverlayDrawsTheGlyph

func TestOverlayDrawsTheGlyph(t *testing.T) {
	s := NoteSpots()[0]
	got := strings.Split(ansiCode.ReplaceAllString(PaintMascot(theme.Default(), []Note{s}), ""), "\n")
	base := strings.Split(Mascot(), "\n")
	want := map[[2]int]rune{{s.Row, s.Col + 1}: '⢸', {s.Row, s.Col + 2}: '⢳', {s.Row + 1, s.Col}: '⠸', {s.Row + 1, s.Col + 1}: '⠟'}
	for y := range base {
		for x, r := range []rune(base[y]) {
			w := r
			if n, ok := want[[2]int{y, x}]; ok {
				w = n
			}
			if g := []rune(got[y])[x]; g != w {
				t.Fatalf("row %d column %d = %q, want %q", y, x, g, w)
			}
		}
	}
}

// END: TestOverlayDrawsTheGlyph

// START: TestNoteCellsUseAccent

func TestNoteCellsUseAccent(t *testing.T) {
	th := theme.Default()
	s := NoteSpots()[0]
	rows := strings.Split(PaintMascot(th, []Note{s}), "\n")
	if !strings.Contains(rows[s.Row], code(th.Accent)) {
		t.Fatalf("row %d %q lacks the Accent colour %s of the note", s.Row, rows[s.Row], code(th.Accent))
	}
}

// END: TestNoteCellsUseAccent

// START: step

// always is a random source that always returns 0; never one that never returns 0.
func always(int) int  { return 0 }
func never(n int) int { return n - 1 }

func TestStepSpawnsAtAFreePlaceInTheLowerRows(t *testing.T) {
	got := StepNotes(nil, always)
	if len(got) != 1 || got[0].Row < 6 {
		t.Fatalf("StepNotes(nil) = %+v, want one note in rows 6 to 11", got)
	}
	free := false
	for _, s := range NoteSpots() {
		free = free || (s.Row == got[0].Row && s.Col == got[0].Col)
	}
	if !free {
		t.Fatalf("the new note %+v is not at a free place", got[0])
	}
}

func TestStepNeverShowsMoreThanThree(t *testing.T) {
	var notes []Note
	for i := 0; i < 100; i++ {
		notes = StepNotes(notes, always)
		if len(notes) > 3 {
			t.Fatalf("step %d: %d notes, want at most 3", i, len(notes))
		}
	}
}

func TestStepRisesOneRowEverySecondTick(t *testing.T) {
	notes := []Note{{Row: 8, Col: NoteSpots()[0].Col}}
	if !spotFree(8, notes[0].Col) || !spotFree(7, notes[0].Col) {
		t.Skip("the chosen column is not free in rows 7 and 8")
	}
	one := StepNotes(notes, never)
	two := StepNotes(one, never)
	if len(one) != 1 || one[0].Row != 8 || len(two) != 1 || two[0].Row != 7 {
		t.Fatalf("after one tick %+v, after two %+v; want row 8 then row 7", one, two)
	}
}

func TestStepEndsAfterEightTicks(t *testing.T) {
	notes := []Note{{Row: 9, Col: 0, Age: 7}}
	if !spotFree(9, 0) {
		t.Skip("column 0 is not free at row 9")
	}
	if got := StepNotes(notes, never); len(got) != 0 {
		t.Fatalf("the note is still there at its eighth tick: %+v", got)
	}
}

func TestStepSpawnsNothingOnAnotherDraw(t *testing.T) {
	if got := StepNotes(nil, never); len(got) != 0 {
		t.Fatalf("StepNotes(nil, never) = %+v, want none (the chance is 1 in 3)", got)
	}
}

// END: step
