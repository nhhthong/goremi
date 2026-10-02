// The theme picker: the list of themes the user chooses from, drawn in the colours of the selected one.
package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// START: ThemePicker

// ThemePicker lists the given themes, one name per line, with one selected.
type ThemePicker struct {
	entries  []theme.Entry
	selected int
}

// NewThemePicker builds a picker over entries, with the theme called current selected (the first line when none has that name).
func NewThemePicker(entries []theme.Entry, current string) ThemePicker {
	p := ThemePicker{entries: entries}
	for i, e := range entries {
		if e.Name == current {
			p.selected = i
		}
	}
	return p
}

// Selected is the name of the selected theme.
func (p ThemePicker) Selected() string {
	if len(p.entries) == 0 {
		return ""
	}
	return p.entries[p.selected].Name
}

// Update moves the selection one line on ↑ and ↓ and stops at the first and the last line.
func (p ThemePicker) Update(k tea.KeyPressMsg) ThemePicker {
	switch {
	case k.Code == tea.KeyUp && p.selected > 0:
		p.selected--
	case k.Code == tea.KeyDown && p.selected < len(p.entries)-1:
		p.selected++
	}
	return p
}

// View is one line per theme name, drawn in the colours of the selected theme.
func (p ThemePicker) View() string {
	if len(p.entries) == 0 {
		return ""
	}
	t := p.entries[p.selected].Theme
	plain := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Foreground))
	chosen := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(t.Accent)).Background(lipgloss.Color(t.Selected))
	lines := make([]string, len(p.entries))
	for i, e := range p.entries {
		if i == p.selected {
			lines[i] = chosen.Render(e.Name)
		} else {
			lines[i] = plain.Render(e.Name)
		}
	}
	return strings.Join(lines, "\n")
}

// END: ThemePicker
