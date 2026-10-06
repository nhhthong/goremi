// Application model: one screen with a search input and a results list, two focus areas.
package app

import (
	"errors"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"goremi/internal/provider"
	"goremi/internal/ui"
	"goremi/internal/ui/theme"
)

// START: Focus

// Focus is the area that receives the keys.
type Focus int

const (
	FocusInput Focus = iota
	FocusList
)

// END: Focus

// START: Model

// Model is the whole application state.
type Model struct {
	provider provider.Provider
	focus    Focus
	input    ui.SearchInput
	results  ui.Results
	// notice is the one line shown under the search box: a failure, a hint or "no results".
	notice string
	// searching is true from Enter until the result arrives: no second search starts meanwhile.
	searching bool
	// width is the terminal width from the last WindowSizeMsg; 0 until the first one.
	width int
	theme theme.Theme
	// player plays the resolved URL; nil until WithPlayer sets one.
	player Player
	// plays queues Play calls for the player; set with it.
	plays *playQueue
	// mouse tells whether the view asks for mouse events.
	mouse bool
	// playing is the track last asked to play; artist is the artist line of the panel, empty until a track plays.
	playing provider.Track
	artist  string
	// paused tells the play/pause glyph; the keys toggle it, as the app does not hear mpv's pause.
	paused bool
	// elapsed and total are what the player last reported; total 0 means the track's own length.
	elapsed, total time.Duration
	// ticking is true while a tick loop runs, so a second play does not start another.
	ticking bool
	// log writes a line to the log file; nil writes nothing.
	log func(format string, a ...any)
	// spectrum is the config key show_spectrum; bars are the heights the spectrum draws now, framing tells that its frame tick runs, ended that the last track of the list ended.
	spectrum       bool
	bars           [ui.SpectrumBars]int
	framing, ended bool
	// intn is the random source of the idle bars (rand.Intn unless a test sets another).
	intn func(n int) int
}

// New starts with the focus on the search input (spec §3). The spectrum is off until WithSpectrum turns it on: NewFromConfig does, as show_spectrum is true unless the config says otherwise.
func New(p provider.Provider) Model {
	return Model{provider: p, focus: FocusInput, results: ui.NewResults(p, "", nil), theme: theme.Default(), mouse: true}
}

func (m Model) Focus() Focus  { return m.focus }
func (m Model) Query() string { return m.input.Value() }

// Selected is the selected line of the results list.
func (m Model) Selected() int { return m.results.Selected() }

// Tracks are the tracks shown in the results list.
func (m Model) Tracks() []provider.Track { return m.results.Tracks() }

// Theme is the colours the model draws with; dark until WithTheme sets another.
func (m Model) Theme() theme.Theme { return m.theme }

// WithTheme returns a copy that draws with t.
func (m Model) WithTheme(t theme.Theme) Model {
	m.theme = t
	return m
}

// WithSpectrum returns a copy that draws the spectrum at the top of the player panel, or not (the config key show_spectrum).
func (m Model) WithSpectrum(on bool) Model {
	m.spectrum = on
	return m
}

// Spectrum tells whether the spectrum shows at the top of the player panel.
func (m Model) Spectrum() bool { return m.spectrum }

// WithLog returns a copy that writes its log lines with f.
func (m Model) WithLog(f func(format string, a ...any)) Model {
	m.log = f
	return m
}

// WithMouse returns a copy that asks for mouse events, or not.
func (m Model) WithMouse(on bool) Model {
	m.mouse = on
	return m
}

// WithFocus returns a copy with the focus set.
func (m Model) WithFocus(f Focus) Model {
	m.focus = f
	return m
}

// END: Model

// START: Update

// Init waits for the events of the player, when there is one.
func (m Model) Init() tea.Cmd {
	if m.player == nil {
		return nil
	}
	return waitEvent(m.player.Events())
}

// PlayMsg asks the player to play the track; Enter on a track line returns it.
type PlayMsg struct{ Track provider.Track }

// searchedMsg carries the result of a search; the model starts using it with the next tasks.
type searchedMsg struct {
	query  string
	tracks []provider.Track
	err    error
}

// Update handles keys: Ctrl+C quits anywhere; in the input Esc quits, Tab moves to the list, Enter searches, other keys edit the query; in the list Tab and Esc return to the input, q quits, ↑/↓ select and Enter plays the track or loads more. The focus moves to the list when results arrive without error.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(frameMsg); ok {
		return m.onFrame()
	}
	if r, ok := msg.(searchedMsg); ok {
		m.searching = false
		if r.err != nil {
			m.notice = searchMessage(r.err)
			m.focus = FocusInput
		} else {
			m.focus = FocusList
			m.notice = ""
			m.results = ui.NewResults(m.provider, r.query, r.tracks)
			if len(r.tracks) == 0 {
				m.focus = FocusInput
				m.notice = `No results for "` + r.query + `".`
			}
		}
		return m, nil
	}
	if p, ok := msg.(PlayMsg); ok {
		return m.startPlay(p.Track)
	}
	if r, ok := msg.(resolvedMsg); ok {
		return m.playResolved(r)
	}
	if p, ok := msg.(playedMsg); ok {
		return m.played(p)
	}
	if _, ok := msg.(tickMsg); ok {
		return m.onTick()
	}
	if p, ok := msg.(progressMsg); ok {
		return m.onProgress(p)
	}
	if c, ok := msg.(tea.MouseClickMsg); ok {
		return m.click(c.Mouse())
	}
	if ev, ok := msg.(playerEventMsg); ok {
		return m.onPlayerEvent(ev)
	}
	if d, ok := msg.(detailsMsg); ok {
		return m.showDetails(d)
	}
	if w, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = w.Width
		return m, nil
	}
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		var cmd tea.Cmd
		m.results, cmd = m.results.Update(msg)
		if err := m.results.LoadErr(); err != nil {
			m.notice = searchMessage(err)
			m.results = m.results.ClearLoadErr()
		}
		return m, cmd
	}
	if k.Code == 'c' && k.Mod&tea.ModCtrl != 0 {
		return m, tea.Quit
	}
	m.notice = ""
	if m.focus == FocusList {
		switch {
		case k.Code == tea.KeyTab || k.Code == tea.KeyEscape:
			m.focus = FocusInput
		case k.Code == 'q':
			return m, tea.Quit
		case k.Code == 'k' || k.Code == tea.KeySpace:
			m = m.pauseKey()
		case k.Code == 'j' || k.Code == tea.KeyLeft:
			m.seekKey(-seekStep)
		case k.Code == 'l' || k.Code == tea.KeyRight:
			m.seekKey(seekStep)
		case k.Code == 'n':
			return m, m.stepTrack(1)
		case k.Code == 'p':
			return m, m.stepTrack(-1)
		case k.Code == tea.KeyUp || k.Code == tea.KeyDown:
			m.results, _ = m.results.Update(k)
		case k.Code == tea.KeyEnter:
			if tracks, i := m.results.Tracks(), m.results.Selected(); i < len(tracks) {
				return m, func() tea.Msg { return PlayMsg{Track: tracks[i]} }
			}
			var cmd tea.Cmd
			m.results, cmd = m.results.Update(k)
			return m, cmd
		}
		return m, nil
	}
	if m.focus == FocusInput {
		switch k.Code {
		case tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyTab:
			m.focus = FocusList
			return m, nil
		case tea.KeyEnter:
			if m.searching {
				return m, nil
			}
			p, q := m.provider, strings.TrimSpace(m.input.Value())
			if q == "" {
				m.notice = "Type something to search."
				return m, nil
			}
			m.searching = true
			return m, func() tea.Msg {
				tracks, err := p.Search(q, 1)
				return searchedMsg{query: q, tracks: tracks, err: err}
			}
		}
		m.input = m.input.Update(k)
	}
	return m, nil
}

// END: Update

// START: View

// sideBySideMin is the width, in columns, from which the list and the panel sit side by side.
const sideBySideMin = 80

// View shows the hints, the search input, the message of a failed search, then the list and the player panel: side by side from 80 columns, stacked (panel, list) below, and without a panel until the width is known.
func (m Model) View() tea.View {
	hint := lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color(m.theme.Muted))
	label := lipgloss.NewStyle().Foreground(lipgloss.Color(m.theme.Accent))
	hintText := hint.Render("Ctrl+C: quit · run goremi theme to choose a theme")
	if m.width > 0 { // a hint wider than the terminal would wrap and move the Search: line
		hintText = ansi.Truncate(hintText, m.width, "…")
	}
	out := hintText + "\n" + label.Render("Search:") + " " + m.input.Value()
	if m.notice != "" {
		out += "\n" + m.notice
	}
	panel := ui.PlayerPanel(m.theme)
	if m.artist != "" { // a track has played: the logo goes
		total := m.playing.Duration
		if m.total > 0 {
			total = m.total
		}
		panel = ui.ArtistLine(m.theme, m.artist) + "\n" + ui.TitleLine(m.theme, m.playing.Title) +
			"\n" + ui.BarLine(m.theme, ui.PanelWidth, m.elapsed, total) + "\n" + ui.ClockText(m.elapsed, total) + "\n" + ui.ControlsLine(m.theme, m.paused)
		if m.spectrum { // the spectrum takes the top of the panel, above the artist line
			panel = strings.Join(ui.PaintSpectrum(m.theme, ui.SpectrumRows(m.bars)), "\n") + "\n" + panel
		}
	}
	list := ""
	if len(m.results.Tracks()) > 0 {
		list = ui.PaintResults(m.theme, m.results.RenderWidth(m.listWidth()), m.results.Selected(), len(m.results.Tracks()))
	}
	switch {
	case m.width == 0: // size not known yet: no panel
		if list != "" {
			out += "\n" + list
		}
	case m.width >= sideBySideMin && list != "":
		out += "\n" + ui.JoinPanes(list, panel)
	case m.width >= sideBySideMin:
		out += "\n" + panel
	default: // narrow: search, panel, list from top to bottom
		out += "\n" + panel
		if list != "" {
			out += "\n" + list
		}
	}
	if m.focus == FocusList && m.artist != "" { // a track plays and the list has focus
		out += "\n" + hint.Render(playbackHint)
	}
	v := tea.NewView(out)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeNone
	if m.mouse {
		v.MouseMode = tea.MouseModeCellMotion
	}
	return v
}

// playbackHint tells the playback keys; it shows while a track plays and the list has focus.
const playbackHint = "j -10s  k pause  l +10s  p prev  n next"

// listWidth is the width a list line may have: beside the 40-column panel and its two-space gap from sideBySideMin, the whole width when stacked, 0 (no limit) until the width is known.
func (m Model) listWidth() int {
	if m.width >= sideBySideMin {
		return m.width - 42
	}
	return m.width
}

// searchMessage is the one-line text under the search box for a failed search.
func searchMessage(err error) string {
	switch {
	case errors.Is(err, provider.ErrNetwork):
		return "Unable to search. Check your internet connection."
	case errors.Is(err, exec.ErrNotFound):
		return "yt-dlp not found. Install yt-dlp and try again."
	case errors.Is(err, provider.ErrTimeout):
		return "Search timed out. Press Enter to retry."
	}
	return "Search failed. Press Enter to retry."
}

// END: View
