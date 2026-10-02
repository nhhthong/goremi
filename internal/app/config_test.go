// Tests for the config file location.
package app

import (
	"runtime"
	"testing"
)

// START: TestConfigPath

func TestConfigPath(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("XDG_CONFIG_HOME decides the config dir on Linux only")
	}
	t.Setenv("XDG_CONFIG_HOME", "/tmp/cfg")
	got, err := ConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if want := "/tmp/cfg/goremi/config.toml"; got != want {
		t.Fatalf("ConfigPath() = %q, want %q", got, want)
	}
}

// END: TestConfigPath
