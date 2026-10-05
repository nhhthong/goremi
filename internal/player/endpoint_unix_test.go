//go:build !windows

// Tests for the private mpv endpoint: where the socket lies and who may enter its directory.
package player

import (
	"os"
	"path/filepath"
	"testing"
)

// START: TestEndpointDirMode0700

func TestEndpointDirMode0700(t *testing.T) {
	tmp := isolatedTmp(t)
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer p.Close()
	dirs := endpointDirs(t, tmp)
	if len(dirs) != 1 {
		t.Fatalf("directories = %v, want one", dirs)
	}
	info, err := os.Stat(dirs[0])
	if err != nil || info.Mode().Perm() != 0o700 {
		t.Fatalf("mode of %s = %v (%v), want 0700", dirs[0], info.Mode().Perm(), err)
	}
}

// END: TestEndpointDirMode0700

// START: TestEndpointPathDiffersPerStart

func TestEndpointPathDiffersPerStart(t *testing.T) {
	tmp := isolatedTmp(t)
	for i := 0; i < 2; i++ {
		p, err := Start()
		if err != nil {
			t.Fatalf("Start() #%d error: %v", i+1, err)
		}
		defer p.Close()
	}
	dirs := endpointDirs(t, tmp)
	if len(dirs) != 2 || dirs[0] == dirs[1] {
		t.Fatalf("directories = %v, want two different ones", dirs)
	}
}

// END: TestEndpointPathDiffersPerStart

// START: TestEndpointInPrivateDir

func TestEndpointInPrivateDir(t *testing.T) {
	tmp := isolatedTmp(t)
	p, err := Start()
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer p.Close()
	dirs := endpointDirs(t, tmp)
	if len(dirs) != 1 {
		t.Fatalf("directories = %v, want one", dirs)
	}
	entries, _ := filepath.Glob(filepath.Join(dirs[0], "*"))
	if len(entries) != 1 {
		t.Fatalf("entries in %s = %v, want one socket", dirs[0], entries)
	}
	if info, err := os.Stat(entries[0]); err != nil || info.Mode()&os.ModeSocket == 0 {
		t.Fatalf("%s is not a socket (%v)", entries[0], err)
	}
	if got, err := p.Property("idle-active"); err != nil || got != true {
		t.Fatalf("Property(idle-active) = %v, %v, want true, nil", got, err)
	}
}

// END: TestEndpointInPrivateDir
