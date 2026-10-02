// Config file for goremi: where it lives and what it holds.
package app

import (
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
	ShowArtwork  bool
	ShowSpectrum bool
}

// configFile is the TOML shape; the show keys are pointers to tell an absent key from false.
type configFile struct {
	UI struct {
		Theme        string `toml:"theme"`
		ShowArtwork  *bool  `toml:"show_artwork"`
		ShowSpectrum *bool  `toml:"show_spectrum"`
	} `toml:"ui"`
}

// defaultConfig is what the app uses without a (valid) config file.
func defaultConfig() Config { return Config{Theme: "dark", ShowArtwork: true, ShowSpectrum: true} }

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
	if f.UI.ShowArtwork != nil {
		c.ShowArtwork = *f.UI.ShowArtwork
	}
	if f.UI.ShowSpectrum != nil {
		c.ShowSpectrum = *f.UI.ShowSpectrum
	}
	return c
}

// END: LoadConfig

// START: SaveTheme

// SaveTheme writes name as the theme of the config file at path, creating the directory and keeping the show keys.
func SaveTheme(path, name string) error {
	c := LoadConfig(path)
	var f configFile
	f.UI.Theme, f.UI.ShowArtwork, f.UI.ShowSpectrum = name, &c.ShowArtwork, &c.ShowSpectrum
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
	t, _ := theme.ByName(LoadConfig(path).Theme)
	return New(p).WithTheme(t)
}

// END: NewFromConfig

// START: Run

// Run is the goremi command. With the argument "theme" it first shows the theme picker and, once a theme is saved, opens the app with it;
// a cancelled picker ends the command. Without it the app opens with the theme from the config at path.
// run runs one model to its end and returns the final model.
func Run(args []string, path string, p provider.Provider, run func(tea.Model) (tea.Model, error)) error {
	if len(args) > 0 && args[0] == "theme" {
		final, err := run(NewThemeModel(path))
		if err != nil {
			return err
		}
		if final.(ThemeModel).Chosen() == "" {
			return nil
		}
	}
	_, err := run(NewFromConfig(path, p))
	return err
}

// END: Run
