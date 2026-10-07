//go:build !windows

// Test that scripts/release-build.sh stamps the version into the binaries (install task 20.42).
package install

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// START: TestReleaseBuildStampsVersion

func TestReleaseBuildStampsVersion(t *testing.T) {
	root := repoRoot(t)
	dist := filepath.Join(root, "dist")
	_ = os.RemoveAll(dist)
	t.Cleanup(func() { _ = os.RemoveAll(dist) })
	build := exec.Command("sh", filepath.Join(root, "scripts", "release-build.sh"))
	build.Dir = root
	build.Env = append(os.Environ(), "VERSION=2.3.4")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("release-build.sh error: %v\n%s", err, out)
	}
	info, err := exec.Command("go", "version", "-m", filepath.Join(dist, "goremi_linux_amd64")).CombinedOutput()
	if err != nil || !strings.Contains(string(info), "goremi/internal/app.Version=2.3.4") {
		t.Fatalf("go version -m: %v\n%s\nwant the ldflags to carry goremi/internal/app.Version=2.3.4", err, info)
	}
}

// END: TestReleaseBuildStampsVersion
