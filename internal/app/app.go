// Application model: one screen with a search input and a results list, two focus areas.
package app

import (
	"errors"
	"os/exec"
	"strings"

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
	theme theme.Theme
}

// New starts with the focus on the search input (spec §3).
func New(p provider.Provider) Model {
	return Model{provider: p, focus: FocusInput, results: ui.NewResults(p, "", nil), theme: theme.Default()}
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

// WithFocus returns a copy with the focus set.
func (m Model) WithFocus(f Focus) Model {
	m.focus = f
	return m
}

// END: Model

// START: Update

func (m Model) Init() tea.Cmd { return nil }

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
	out := hint.Render("Ctrl+C: quit · run goremi theme to choose a theme") + "\n" + label.Render("Search:") + " " + m.input.Value()
	if m.notice != "" {
		out += "\n" + m.notice
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
		out += "\n" + ui.JoinPanes(list, ui.PlayerPanel(m.theme))
	case m.width >= sideBySideMin:
		out += "\n" + ui.PlayerPanel(m.theme)
	default: // narrow: search, panel, list from top to bottom
		out += "\n" + ui.PlayerPanel(m.theme)
		if list != "" {
			out += "\n" + list
		}
	}
	v := tea.NewView(out)
	v.AltScreen = true
	return v
}

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
