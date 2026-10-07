// Application model: one screen with a search input and a results list, two focus areas.
package app

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

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
	// height is the terminal height from the last WindowSizeMsg; 0 until the first one.
	height int
	theme  theme.Theme
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
	// notes are the music notes drawn around the mascot while a track plays; noting is true while their tick loop runs.
	notes  []ui.Note
	noting bool
	// listTop is the first line of the results list that is shown while the list scrolls.
	listTop int
	// cmdSel is the highlighted line of the command list; cmdClosed is true after Esc until the text changes.
	cmdSel    int
	cmdClosed bool
	// picking is true while the theme selector is open: picker holds it, savedTheme the theme to restore on Esc, configPath the file Enter saves to.
	picking    bool
	picker     ui.ThemePicker
	savedTheme theme.Theme
	configPath string
	// dragging is true while the mouse is held on the bar; dragCell is the cell it is on.
	dragging bool
	dragCell int
	// ticking is true while a tick loop runs, so a second play does not start another.
	ticking bool
	// log writes a line to the log file; nil writes nothing.
	log func(format string, a ...any)
	// spectrum is the config key show_spectrum; bars are the heights the spectrum draws now, framing tells that its frame tick runs, ended that the last track of the list ended.
	spectrum       bool
	bars           [ui.SpectrumBars]int
	idle           [ui.SpectrumBars]int
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

// Height is the terminal height from the last WindowSizeMsg; 0 until the first one.
func (m Model) Height() int { return m.height }

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

// WithConfigPath returns a copy whose theme selector saves to the config file at path.
func (m Model) WithConfigPath(path string) Model {
	m.configPath = path
	return m
}

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

// OpenThemeMsg asks the app to open the theme selector; the command `/theme` returns it.
type OpenThemeMsg struct{}

// PlayMsg asks the player to play the track; Enter on a track line returns it.
type PlayMsg struct{ Track provider.Track }

// searchedMsg carries the result of a search; the model starts using it with the next tasks.
type searchedMsg struct {
	query  string
	tracks []provider.Track
	err    error
}

// Update handles keys: Ctrl+C quits anywhere; in the input Tab moves to the list, Enter searches, other keys edit the query and Esc does nothing; in the list Tab and Esc return to the input, ↑/↓ select and Enter plays the track or loads more. The focus moves to the list when results arrive without error.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(noteTickMsg); ok {
		return m.onNoteTick()
	}
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
			m.results, m.listTop = ui.NewResults(m.provider, r.query, r.tracks), 0
			if len(r.tracks) == 0 {
				m.focus = FocusInput
				m.notice = `No results for "` + r.query + `".`
			}
		}
		return m, nil
	}
	if _, ok := msg.(OpenThemeMsg); ok {
		m.picking, m.savedTheme = true, m.theme
		m.picker = ui.NewThemePicker(theme.All(), m.themeName())
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
	if c, ok := msg.(tea.MouseMotionMsg); ok {
		return m.onMouseMove(c.Mouse()), nil
	}
	if c, ok := msg.(tea.MouseReleaseMsg); ok {
		return m.onMouseRelease(c.Mouse()), nil
	}
	if ev, ok := msg.(playerEventMsg); ok {
		return m.onPlayerEvent(ev)
	}
	if d, ok := msg.(detailsMsg); ok {
		return m.showDetails(d)
	}
	if w, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = w.Width
		m.height = w.Height
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
	if m.picking {
		return m.pickerKey(k)
	}
	if m.focus == FocusList {
		switch {
		case k.Code == tea.KeyTab || k.Code == tea.KeyEscape:
			m.focus = FocusInput
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
		case k.Code == tea.KeyDown && m.results.OnLoadMore(): // Down past the last line jumps to the search bar
			m.focus = FocusInput
		case k.Code == tea.KeyUp || k.Code == tea.KeyDown:
			m.results, _ = m.results.Update(k)
			m.listTop = m.visibleTop()
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
		if next, cmd, handled := m.commandKey(k); handled {
			return next, cmd
		}
		switch k.Code {
		case tea.KeyTab:
			m.focus = FocusList
			return m, nil
		case tea.KeyUp: // back to the list, on the line it was left
			if len(m.results.Tracks()) > 0 {
				m.focus = FocusList
			}
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
		before := m.input.Value()
		m.input = m.input.Update(k)
		if m.input.Value() != before { // new text: the command list starts again
			m.cmdSel, m.cmdClosed = 0, false
		}
	}
	return m, nil
}

// END: Update

// START: View

// sideBySideMin is the width, in columns, from which the list and the panel sit side by side.
const sideBySideMin = ui.WideMin

// escHint is the line above the search bar while the list has the focus.
const escHint = "Press Esc to return to search"

// View is the content, then the lines that belong to the search bar (the command list or the theme selector, the Esc hint, the message), then the search bar: a rule, the `❯` line and a rule, the last three lines. With the height known the content is padded so the bar ends the Height lines.
func (m Model) View() tea.View {
	prompt := lipgloss.NewStyle().Foreground(lipgloss.Color(m.theme.Accent))
	faint := lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color(m.theme.Muted))
	above := m.aboveLines()
	var lines []string
	if content := m.content(); content != "" {
		lines = strings.Split(content, "\n")
	}
	if limit := m.height - 3 - len(above); m.height > 0 && len(lines) > limit { // the bar stays fixed: the content is cut from the bottom
		lines = lines[:max(limit, 0)]
	}
	for m.height > 0 && len(lines)+len(above) < m.height-3 {
		lines = append(lines, "")
	}
	if m.focus != FocusInput {
		prompt = faint // the list has the focus: the prompt is dimmed
	}
	rule := lipgloss.NewStyle().Foreground(lipgloss.Color(m.theme.Border)).Render(strings.Repeat("─", m.barWidth()))
	lines = append(append(lines, above...), rule, prompt.Render("❯")+" "+m.input.Value(), rule)
	v := tea.NewView(strings.Join(lines, "\n"))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeNone
	if m.mouse {
		v.MouseMode = tea.MouseModeCellMotion
	}
	return v
}

// aboveLines are the lines between the content and the search bar: the command list or the theme selector, the Esc hint while the list has the focus, and the message.
func (m Model) aboveLines() []string {
	var above []string
	if m.picking {
		above = strings.Split(m.picker.View(), "\n")
	} else {
		above = m.commandLines()
	}
	if m.focus == FocusList && !m.picking {
		above = append(above, lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color(m.theme.Muted)).Render(escHint))
	}
	if m.notice != "" {
		above = append(above, lipgloss.NewStyle().Foreground(lipgloss.Color(m.theme.Error)).Render(m.notice))
	}
	return above
}

// listRows is the number of rows the results list may take: what is left under the header, the lines above the bar and the bar (and under the panel when it is stacked above the list); 0 means no limit, while the size is not known.
func (m Model) listRows() int {
	if m.width == 0 || m.height == 0 {
		return 0
	}
	rows := m.height - 3 - len(m.aboveLines()) - ui.HeaderRows(m.width) - ui.BannerRows(m.width)
	if m.artist != "" && m.width < sideBySideMin {
		rows -= strings.Count(m.panel(), "\n") + 1
	}
	return max(rows, 1)
}

// visibleTop is the first line of the list window: the stored top moved just enough that the selected line is inside the window.
func (m Model) visibleTop() int {
	rows, n, sel := m.listRows(), m.results.LineCount(), m.results.Selected()
	if rows <= 0 || n <= rows {
		return 0
	}
	top := m.listTop
	if sel < top {
		top = sel
	}
	if sel >= top+rows {
		top = sel - rows + 1
	}
	return max(0, min(top, n-rows))
}

// window cuts a painted list to the rows it may take, from the first visible line.
func (m Model) window(list string) string {
	rows := m.listRows()
	lines := strings.Split(list, "\n")
	if rows <= 0 || len(lines) <= rows {
		return list
	}
	top := m.visibleTop()
	return strings.Join(lines[top:top+rows], "\n")
}

// content is what sits above the search bar: with the width unknown only the list; else the header (the mascot and the badge, always), then the list at the full width before the first play, or the list and the player panel after it: side by side from 80 columns, panel then list below.
func (m Model) content() string {
	list := ""
	if len(m.results.Tracks()) > 0 {
		list = m.window(ui.PaintResults(m.theme, m.results.RenderWidth(m.listWidth()), m.results.Selected(), len(m.results.Tracks())))
	}
	if m.width == 0 { // size not known yet: no header and no panel
		return list
	}
	header := ui.Header(m.theme, m.width, Version, m.notes)
	if ui.BannerRows(m.width) > 0 { // the banner sits above the mascot and the badge
		header = ui.PaintBanner(m.theme) + "\n" + header
	}
	switch {
	case m.artist == "" && list == "":
		return header
	case m.artist == "":
		return header + "\n" + list
	}
	panel := m.panel()
	switch {
	case m.width >= sideBySideMin && list != "":
		return header + "\n" + ui.JoinPanes(list, panel)
	case m.width >= sideBySideMin || list == "":
		return header + "\n" + panel
	}
	return header + "\n" + panel + "\n" + list
}

// barWidth is the width of the rules of the search bar: the view width, the panel width until it is known.
func (m Model) barWidth() int {
	if m.width > 0 {
		return m.width
	}
	return ui.PanelWidth
}

// panel is the player panel of a track that has played: the spectrum when it is on, then the artist, the title, the bar, the clock, the controls and, while the list has the focus, the keys hint.
func (m Model) panel() string {
	elapsed, total := m.shown()
	panel := ui.ArtistLine(m.theme, m.artist) + "\n" + ui.TitleLine(m.theme, m.playing.Title) +
		"\n" + ui.BarLine(m.theme, ui.PanelWidth, elapsed, total) + "\n" + ui.ClockText(elapsed, total) + "\n" + ui.ControlsLine(m.theme, m.paused)
	if m.focus == FocusList {
		panel += "\n" + lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color(m.theme.Muted)).Render(playbackHint)
	}
	if m.spectrum { // the spectrum takes the top of the panel, above the artist line
		base := lipgloss.NewStyle().Foreground(lipgloss.Color(m.theme.Border)).Render(ui.SpectrumBaseline())
		panel = strings.Join(ui.PaintSpectrum(m.theme, ui.SpectrumRows(m.bars)), "\n") + "\n" + base + "\n" + panel
	}
	return panel
}

// END: View

// START: commands

// commands are the words the command box offers.
var commands = []string{"/quit", "/theme"}

// commandMatches are the commands the text of the search bar is the start of; none when the text does not begin with `/` or the list was closed with Esc.
func (m Model) commandMatches() []string {
	v := m.input.Value()
	if !strings.HasPrefix(v, "/") || m.cmdClosed {
		return nil
	}
	var out []string
	for _, c := range commands {
		if strings.HasPrefix(c, v) {
			out = append(out, c)
		}
	}
	return out
}

// commandLines draws the open command list, one line per command, the highlighted one marked and coloured.
func (m Model) commandLines() []string {
	matches := m.commandMatches()
	lines := make([]string, len(matches))
	for i, c := range matches {
		lines[i] = "  " + c
		if i == min(m.cmdSel, len(matches)-1) {
			lines[i] = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(m.theme.Accent)).Render("▶ " + c)
		}
	}
	return lines
}

// commandKey handles the keys of the command box while the search bar has the focus: Down and Up move the highlight, Tab completes, Esc closes the list, Enter runs the highlighted command; Enter on other text that begins with `/` is an unknown command and runs no search.
func (m Model) commandKey(k tea.KeyPressMsg) (Model, tea.Cmd, bool) {
	matches := m.commandMatches()
	if len(matches) == 0 {
		if text := strings.TrimSpace(m.input.Value()); k.Code == tea.KeyEnter && strings.HasPrefix(text, "/") {
			m.notice = fmt.Sprintf("Unknown command %q.", text)
			return m, nil, true
		}
		return m, nil, false
	}
	sel := min(m.cmdSel, len(matches)-1)
	switch k.Code {
	case tea.KeyDown:
		m.cmdSel = min(sel+1, len(matches)-1)
	case tea.KeyUp:
		m.cmdSel = max(sel-1, 0)
	case tea.KeyTab:
		m.input, m.cmdSel = m.input.WithValue(matches[sel]), 0
	case tea.KeyEscape:
		m.cmdClosed = true
	case tea.KeyEnter:
		m.input, m.cmdSel = m.input.WithValue(""), 0
		if matches[sel] == "/quit" {
			return m, tea.Quit, true
		}
		return m, func() tea.Msg { return OpenThemeMsg{} }, true
	default:
		return m, nil, false
	}
	return m, nil, true
}

// END: commands

// START: theme selector

// themeName is the name of the theme the model draws with, default when it is none of the listed ones.
func (m Model) themeName() string {
	for _, e := range theme.All() {
		if e.Theme == m.theme {
			return e.Name
		}
	}
	return "default"
}

// pickerKey handles the keys while the theme selector is open: Up and Down move it and draw the whole view in the highlighted theme, Enter saves the choice and closes it (a failed save is a message and the theme stays), Esc closes it and restores the theme.
func (m Model) pickerKey(k tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch k.Code {
	case tea.KeyUp, tea.KeyDown:
		m.picker = m.picker.Update(k)
		if t, ok := theme.ByName(m.picker.Selected()); ok {
			m.theme = t
		}
	case tea.KeyEnter:
		m.picking = false
		if err := SaveTheme(m.configPath, m.picker.Selected()); err != nil {
			m.notice = "Cannot save the theme: " + err.Error()
		}
	case tea.KeyEscape:
		m.picking, m.theme = false, m.savedTheme
	}
	return m, nil
}

// END: theme selector

// playbackHint tells the playback keys; it shows while a track plays and the list has focus.
const playbackHint = "j -10s  k pause  l +10s  p prev  n next"

// listWidth is the width a list line may have: beside the 40-column panel and its two-space gap from sideBySideMin once a track has played, the whole width before that and when stacked, 0 (no limit) until the width is known.
func (m Model) listWidth() int {
	if m.width >= sideBySideMin && m.artist != "" {
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
