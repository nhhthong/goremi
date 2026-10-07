// The music notes of the mascot: where one fits, and how the notes rise and end.
package ui

// START: Note

// Note is one music note: the row and column of its glyph and its age in ticks.
type Note struct{ Row, Col, Age int }

// The note glyph is two braille rows, `⠀⢸⢳` over `⠸⠟⠀`: noteCells are the four cells it draws, as (row, column) offsets, and noteRunes what goes in them.
var (
	noteCells = [4][2]int{{0, 1}, {0, 2}, {1, 0}, {1, 1}}
	noteRunes = [4]rune{'⢸', '⢳', '⠸', '⠟'}
)

const (
	noteLife   = 8 // ticks a note lives
	noteMax    = 3 // notes at once
	noteSpawn  = 3 // a new note with a chance of 1 in noteSpawn per tick
	noteMinRow = 6 // a new note starts in these rows and below, and rises
)

// END: Note

// START: NoteSpots

// spotFree tells whether the four cells of a note glyph at (row, col) are inside the mascot area and blank in the art.
func spotFree(row, col int) bool {
	for _, c := range noteCells {
		r, x := row+c[0], col+c[1]
		if r < 0 || r >= len(mascotRows) || x < 0 || x >= MascotWidth || []rune(mascotRows[r])[x] != 0x2800 {
			return false
		}
	}
	return true
}

// NoteSpots lists every place where a note fits.
func NoteSpots() []Note {
	var out []Note
	for row := 0; row < len(mascotRows); row++ {
		for col := 0; col < MascotWidth; col++ {
			if spotFree(row, col) {
				out = append(out, Note{Row: row, Col: col})
			}
		}
	}
	return out
}

// END: NoteSpots

// START: StepNotes

// near tells whether two notes would touch: their glyphs are 2 rows and 3 columns wide.
func near(a, b Note) bool {
	return a.Row-b.Row > -3 && a.Row-b.Row < 3 && a.Col-b.Col > -4 && a.Col-b.Col < 4
}

// StepNotes is one tick: every note ages, rises one row every second tick and ends after noteLife ticks or when its place is not free; then, while fewer than noteMax show, a new note appears with a chance of 1 in noteSpawn at a random free place of the lower rows that touches no other note. intn(n) gives a number from 0 to n-1.
func StepNotes(notes []Note, intn func(n int) int) []Note {
	var out []Note
	for _, n := range notes {
		n.Age++
		if n.Age >= noteLife {
			continue
		}
		if n.Age%2 == 0 {
			n.Row--
		}
		if spotFree(n.Row, n.Col) {
			out = append(out, n)
		}
	}
	if len(out) >= noteMax || intn(noteSpawn) != 0 {
		return out
	}
	var spots []Note
	for _, s := range NoteSpots() {
		if s.Row < noteMinRow {
			continue
		}
		touches := false
		for _, n := range out {
			touches = touches || near(s, n)
		}
		if !touches {
			spots = append(spots, s)
		}
	}
	if len(spots) > 0 {
		out = append(out, spots[intn(len(spots))])
	}
	return out
}

// END: StepNotes
