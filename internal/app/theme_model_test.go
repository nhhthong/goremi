// Tests for the `goremi theme` screen: choosing saves the theme, cancelling saves nothing.
package app

import (
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
)

var (
	thDown  = tea.KeyPressMsg{Code: tea.KeyDown}
	thEnter = tea.KeyPressMsg{Code: tea.KeyEnter}
	thEsc   = tea.KeyPressMsg{Code: tea.KeyEscape}
	thCtrlC = tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
)

// press sends the keys to the model in order and returns the model and the command of the last key.
func press(m ThemeModel, keys ...tea.KeyPressMsg) (ThemeModel, tea.Cmd) {
	var cmd tea.Cmd
	for _, k := range keys {
		next, c := m.Update(k)
		m, cmd = next.(ThemeModel), c
	}
	return m, cmd
}

// START: TestThemeEnterChooses

func TestThemeEnterChooses(t *testing.T) {
	m, cmd := press(NewThemeModel(writeConfig(t, "[ui]\ntheme = \"light\"\n")), thDown, thEnter)
	if m.Chosen() != "dark" || !quits(cmd) {
		t.Fatalf("Chosen() = %q, quits = %v; want dark and a quit", m.Chosen(), quits(cmd))
	}
}

// END: TestThemeEnterChooses

// START: TestThemeEnterSaves

func TestThemeEnterSaves(t *testing.T) {
	path := writeConfig(t, "[ui]\ntheme = \"dark\"\n")
	press(NewThemeModel(path), thDown, thEnter)
	if got := LoadConfig(path).Theme; got != "cyberpunk" {
		t.Fatalf("saved Theme = %q, want cyberpunk", got)
	}
}

// END: TestThemeEnterSaves

// cancelled checks that the keys end the screen with no choice and no file.
func cancelled(t *testing.T, keys ...tea.KeyPressMsg) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	m, cmd := press(NewThemeModel(path), keys...)
	if m.Chosen() != "" || !quits(cmd) {
		t.Errorf("Chosen() = %q, quits = %v; want empty and a quit", m.Chosen(), quits(cmd))
	}
	if _, err := os.Stat(path); err == nil {
		t.Errorf("the config file %s exists, want nothing saved", path)
	}
}

// START: TestThemeEscCancels

func TestThemeEscCancels(t *testing.T) { cancelled(t, thDown, thEsc) }

// END: TestThemeEscCancels

// START: TestThemeCtrlCCancels

func TestThemeCtrlCCancels(t *testing.T) { cancelled(t, thDown, thCtrlC) }

// END: TestThemeCtrlCCancels
