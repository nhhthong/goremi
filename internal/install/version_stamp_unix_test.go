//go:build !windows

// Tests that `make build` stamps the version into the binary (ui task 3.4.5).
package install

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// START: TestMakeBuildStampsVersion

func TestMakeBuildStampsVersion(t *testing.T) {
	root := repoRoot(t)
	build := exec.Command("make", "build", "VERSION=1.2.3")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("make build failed: %v\n%s", err, out)
	}
	info := exec.Command("go", "version", "-m", filepath.Join(root, "bin", "goremi"))
	out, err := info.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "goremi/internal/app.Version=1.2.3") {
		t.Fatalf("go version -m: %v\n%s\nwant the ldflags to carry goremi/internal/app.Version=1.2.3", err, out)
	}
}

// END: TestMakeBuildStampsVersion

// START: TestMakeVersionDropsLeadingV

func TestMakeVersionDropsLeadingV(t *testing.T) {
	dir := t.TempDir()
	data, err := os.ReadFile(filepath.Join(repoRoot(t), "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "Makefile"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) {
		c := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@t"}, args...)...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("init", "-q")
	git("commit", "-q", "--allow-empty", "-m", "x")
	git("tag", "v1.2.3")
	dry := exec.Command("make", "-n", "build")
	dry.Dir = dir
	out, err := dry.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "goremi/internal/app.Version=1.2.3") || strings.Contains(string(out), "Version=v1.2.3") {
		t.Fatalf("make -n build: %v\n%s\nwant Version=1.2.3 without the leading v", err, out)
	}
}

// END: TestMakeVersionDropsLeadingV
