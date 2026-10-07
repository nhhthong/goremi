// Tests of the in-app theme selector opened by /theme: preview, save, cancel and a failed save (theme tasks 8.15 to 8.19).
package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/ui/theme"
)

var themeNames = []string{"default", "catppuccin", "dracula", "gruvbox", "nord", "rosepine", "tokyonight"}

// START: helpers

// selectorAt opens the selector on a 100x30 model that draws with t and saves to path.
func selectorAt(t theme.Theme, path string) Model {
	m := sizeTo(New(fakeProvider{}).WithTheme(t).WithConfigPath(path), 100, 30)
	next, _ := m.Update(OpenThemeMsg{})
	return next.(Model)
}

func keyDown(m Model) Model { n, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown}); return n.(Model) }
func keyUp(m Model) Model   { n, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyUp}); return n.(Model) }
func keyEnter(m Model) Model {
	n, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	return n.(Model)
}
func keyEsc(m Model) Model { n, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape}); return n.(Model) }

// aboveBar returns the n lines above the Search: line, plain.
func aboveBar(m Model, n int) []string {
	lines := plainLines(m)
	s := searchIndex(lines)
	return lines[max(s-n, 0):s]
}

// highlighted is the name on the one line of the seven above the bar that carries a colour code.
func highlighted(m Model) string {
	lines := viewLines(m)
	s := searchIndex(plainLines(m))
	for _, l := range lines[s-7 : s] {
		if strings.Contains(l, "\x1b[") {
			return plain(l)
		}
	}
	return ""
}

// blocked is a config path below a regular file, so a save fails.
func blocked(t *testing.T) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	return filepath.Join(file, "config.toml")
}

// END: helpers

// START: TestSelectorListsThemes

func TestSelectorListsThemes(t *testing.T) {
	got := aboveBar(selectorAt(theme.Default(), filepath.Join(t.TempDir(), "c.toml")), 7)
	if strings.Join(got, ",") != strings.Join(themeNames, ",") {
		t.Fatalf("lines above the bar = %q, want the seven theme names", got)
	}
}

// END: TestSelectorListsThemes

// START: TestSelectorHighlightsCurrent

func TestSelectorHighlightsCurrent(t *testing.T) {
	for name, th := range map[string]theme.Theme{"default": theme.Default(), "dracula": theme.Dracula()} {
		if got := highlighted(selectorAt(th, filepath.Join(t.TempDir(), "c.toml"))); got != name {
			t.Errorf("current theme %s: highlighted line = %q", name, got)
		}
	}
}

// END: TestSelectorHighlightsCurrent

// START: TestSelectorArrowsPreview

func TestSelectorArrowsPreview(t *testing.T) {
	m := keyDown(selectorAt(theme.Default(), filepath.Join(t.TempDir(), "c.toml")))
	if got := highlighted(m); got != "catppuccin" {
		t.Fatalf("highlighted after Down = %q, want catppuccin", got)
	}
	lines := viewLines(m)
	if !strings.Contains(lines[len(lines)-2], "38;2;203;166;247") {
		t.Fatalf("Search: line %q is not drawn in the Catppuccin Accent 203;166;247", lines[len(lines)-2])
	}
}

// END: TestSelectorArrowsPreview

// START: TestSelectorUpStopsAtFirst

func TestSelectorUpStopsAtFirst(t *testing.T) {
	if got := highlighted(keyUp(selectorAt(theme.Default(), filepath.Join(t.TempDir(), "c.toml")))); got != "default" {
		t.Fatalf("highlighted after Up on the first = %q, want default", got)
	}
}

// END: TestSelectorUpStopsAtFirst

// START: TestSelectorEnterSaves

func TestSelectorEnterSaves(t *testing.T) {
	path := filepath.Join(t.TempDir(), "c.toml")
	m := keyEnter(keyDown(selectorAt(theme.Default(), path)))
	if got := LoadConfig(path).Theme; got != "catppuccin" {
		t.Fatalf("saved theme = %q, want catppuccin", got)
	}
	if m.Theme() != theme.Catppuccin() || strings.Contains(plain(m.View().Content), "tokyonight") {
		t.Fatalf("theme %v, selector still open or theme not kept", m.Theme())
	}
}

// END: TestSelectorEnterSaves

// START: TestSelectorEnterKeepsOtherKeys

func TestSelectorEnterKeepsOtherKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "c.toml")
	if err := os.WriteFile(path, []byte("[ui]\nmouse = false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	keyEnter(keyDown(selectorAt(theme.Default(), path)))
	if c := LoadConfig(path); c.Theme != "catppuccin" || c.Mouse {
		t.Fatalf("config after save = %+v, want theme catppuccin and mouse false kept", c)
	}
}

// END: TestSelectorEnterKeepsOtherKeys

// START: TestSelectorEscRestores

func TestSelectorEscRestores(t *testing.T) {
	path := filepath.Join(t.TempDir(), "c.toml")
	m := keyEsc(keyDown(selectorAt(theme.Default(), path)))
	if _, err := os.Stat(path); err == nil {
		t.Fatal("Esc wrote the config file")
	}
	if m.Theme() != theme.Default() || strings.Contains(plain(m.View().Content), "tokyonight") {
		t.Fatalf("theme %v, selector still open or theme not restored", m.Theme())
	}
}

// END: TestSelectorEscRestores

// START: TestSaveFailureKeepsTheme

func TestSaveFailureKeepsTheme(t *testing.T) {
	m := keyEnter(keyDown(selectorAt(theme.Default(), blocked(t))))
	lines := plainLines(m)
	if s := searchIndex(lines); s < 1 || !strings.HasPrefix(lines[s-1], "Cannot save the theme: ") {
		t.Fatalf("line above the bar = %q, want Cannot save the theme: <error>", lines[searchIndex(lines)-1])
	}
	if m.Theme() != theme.Catppuccin() || strings.Contains(strings.Join(lines, "\n"), "tokyonight") {
		t.Fatalf("theme %v or the selector is still open", m.Theme())
	}
}

// END: TestSaveFailureKeepsTheme

// START: TestSaveFailureAppKeepsRunning

func TestSaveFailureAppKeepsRunning(t *testing.T) {
	m := keyEnter(keyDown(selectorAt(theme.Default(), blocked(t))))
	m, _ = typed(m, "x")
	if m.Query() != "x" {
		t.Fatalf("Query() = %q after a failed save and a key, want x", m.Query())
	}
}

// END: TestSaveFailureAppKeepsRunning
