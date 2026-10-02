// The screen behind `goremi theme`: the theme picker, saving the choice to the config file.
package app

import (
	tea "charm.land/bubbletea/v2"

	"goremi/internal/ui"
	"goremi/internal/ui/theme"
)

// START: ThemeModel

// ThemeModel lets the user pick a theme from the list and saves it to the config file.
type ThemeModel struct {
	path   string
	picker ui.ThemePicker
	chosen string
}

// NewThemeModel starts on the theme the config file at path names (dark without a config).
func NewThemeModel(path string) ThemeModel {
	return ThemeModel{path: path, picker: ui.NewThemePicker(theme.All(), LoadConfig(path).Theme)}
}

// Chosen is the name of the theme that was saved, empty when the user cancelled or the save failed.
func (m ThemeModel) Chosen() string { return m.chosen }

func (m ThemeModel) Init() tea.Cmd { return nil }

// Update: Enter saves the selected theme and quits; Esc and Ctrl+C quit without saving; ↑ and ↓ move the selection.
func (m ThemeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch {
	case k.Code == tea.KeyEnter:
		if name := m.picker.Selected(); SaveTheme(m.path, name) == nil {
			m.chosen = name
		}
		return m, tea.Quit
	case k.Code == tea.KeyEscape, k.Code == 'c' && k.Mod&tea.ModCtrl != 0:
		return m, tea.Quit
	}
	m.picker = m.picker.Update(k)
	return m, nil
}

func (m ThemeModel) View() tea.View { return tea.NewView(m.picker.View()) }

// END: ThemeModel
