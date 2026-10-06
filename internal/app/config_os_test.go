// Test of the config file location on each system: the directory comes from the variable that system's os.UserConfigDir reads.
package app

import (
	"path/filepath"
	"runtime"
	"testing"
)

// START: TestConfigPathPerOS

func TestConfigPathPerOS(t *testing.T) {
	base := t.TempDir()
	var want string
	switch runtime.GOOS {
	case "linux":
		t.Setenv("XDG_CONFIG_HOME", base)
		want = filepath.Join(base, "goremi", "config.toml")
	case "darwin":
		t.Setenv("HOME", base)
		want = filepath.Join(base, "Library", "Application Support", "goremi", "config.toml")
	case "windows":
		t.Setenv("AppData", base)
		want = filepath.Join(base, "goremi", "config.toml")
	default:
		t.Fatalf("no expectation for %s", runtime.GOOS)
	}
	got, err := ConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("ConfigPath() = %q, want %q", got, want)
	}
}

// END: TestConfigPathPerOS
