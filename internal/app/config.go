// Config file for goremi: where it lives and what it holds.
package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/BurntSushi/toml"

	"goremi/internal/provider"
	"goremi/internal/ui/theme"
)

// START: ConfigPath

// ConfigPath is the config file: <user config dir>/goremi/config.toml.
func ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "goremi", "config.toml"), nil
}

// END: ConfigPath

// START: LoadConfig

// Config is what config.toml holds under [ui].
type Config struct {
	Theme        string
	ShowSpectrum bool
	Mouse        bool
}

// configFile is the TOML shape; the show keys are pointers to tell an absent key from false.
type configFile struct {
	UI struct {
		Theme        string `toml:"theme"`
		ShowSpectrum *bool  `toml:"show_spectrum"`
		Mouse        *bool  `toml:"mouse"`
	} `toml:"ui"`
}

// defaultConfig is what the app uses without a (valid) config file.
func defaultConfig() Config {
	return Config{Theme: "default", ShowSpectrum: true, Mouse: true}
}

// LoadConfig reads the file at path. A missing or unparsable file, a theme not in the list and an absent key all fall back to the defaults.
func LoadConfig(path string) Config {
	var f configFile
	if _, err := toml.DecodeFile(path, &f); err != nil {
		return defaultConfig()
	}
	c := defaultConfig()
	if _, ok := theme.ByName(f.UI.Theme); ok {
		c.Theme = f.UI.Theme
	}
	if f.UI.ShowSpectrum != nil {
		c.ShowSpectrum = *f.UI.ShowSpectrum
	}
	if f.UI.Mouse != nil {
		c.Mouse = *f.UI.Mouse
	}
	return c
}

// END: LoadConfig

// START: SaveTheme

// SaveTheme writes name as the theme of the config file at path, creating the directory and keeping the show keys.
func SaveTheme(path, name string) error {
	c := LoadConfig(path)
	var f configFile
	f.UI.Theme, f.UI.ShowSpectrum, f.UI.Mouse = name, &c.ShowSpectrum, &c.Mouse
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return toml.NewEncoder(file).Encode(f)
}

// END: SaveTheme

// START: NewFromConfig

// NewFromConfig starts the app with the theme named in the config file at path (dark when there is none).
func NewFromConfig(path string, p provider.Provider) Model {
	c := LoadConfig(path)
	t, _ := theme.ByName(c.Theme)
	return New(p).WithTheme(t).WithMouse(c.Mouse).WithSpectrum(c.ShowSpectrum).WithConfigPath(path)
}

// END: NewFromConfig

// START: Run

// Run is the goremi command. Any argument opens nothing and returns a usage error `Unknown command "<x>".` (exit code 2): `goremi theme` and `goremi help` are gone, the theme is chosen with /theme inside the app.
// Without arguments the app opens with the theme from the config at path. out is not written to any more; it stays so main and the tests keep one signature. run runs one model to its end and returns the final model.
func Run(args []string, path string, p provider.Provider, out io.Writer, run func(tea.Model) (tea.Model, error)) error {
	if len(args) > 0 {
		return usageError{fmt.Sprintf("Unknown command %q.", args[0])}
	}
	pl := newMpvPlayer()
	pl.spectrum = LoadConfig(path).ShowSpectrum // mpv gets the band filter only when the spectrum shows
	var logf func(string, ...any)
	if path, err := LogPath(); err == nil { // the log gets the play line, the stderr of yt-dlp and the warn and error messages of mpv
		lg, _ := OpenLog(path) // a log that cannot be opened discards (never a crash)
		pl.log, logf = lg.Printf, lg.Printf
		if yt, ok := p.(*provider.YouTubeProvider); ok && yt.Log == nil {
			yt.Log = lg.Printf
		}
	}
	defer pl.Close()
	_, err := run(NewFromConfig(path, p).WithPlayer(pl).WithLog(logf))
	return err
}

// usageError is a mistake in the command line, as opposed to a failure while running.
type usageError struct{ msg string }

func (e usageError) Error() string { return e.msg }

// ExitCode is the exit code of the command for the error Run returned: 0 for none, 2 for a mistake in the command line, 1 otherwise.
func ExitCode(err error) int {
	var u usageError
	switch {
	case err == nil:
		return 0
	case errors.As(err, &u):
		return 2
	}
	return 1
}

// END: Run
