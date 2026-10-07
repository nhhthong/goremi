// Regression: the GOREMI logo art of the old panel is drawn nowhere (ui task 3.4.8). It lives in its own package, reading only the public seam of app, so `clio test red` can run it against the commit before the fix.
package regress

import (
	"regexp"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/app"
	"goremi/internal/provider"
)

type nothing struct{}

func (nothing) Search(string, int) ([]provider.Track, error)     { return nil, nil }
func (nothing) Details(t provider.Track) (provider.Track, error) { return t, nil }
func (nothing) Resolve(provider.Track) (string, error)           { return "", nil }

var ansiCode = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// START: TestLogoArtIsGone

func TestLogoArtIsGone(t *testing.T) {
	const art = "▄▀▀▀ ▄▀▀▄ █▀▀▄" // the first row of the old GOREMI art
	for _, width := range []int{0, 40, 80, 100} {
		next, _ := app.New(nothing{}).Update(tea.WindowSizeMsg{Width: width, Height: 30})
		if got := ansiCode.ReplaceAllString(next.(app.Model).View().Content, ""); strings.Contains(got, art) {
			t.Errorf("width %d still draws the GOREMI art:\n%s", width, got)
		}
	}
}

// END: TestLogoArtIsGone
