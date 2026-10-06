//go:build !windows

// Tests of scripts/release-build.sh: the five release files, their headers, and that dist/ stays out of git.
package install

import (
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"
)

// START: buildDist

// buildDist runs the script from the repository root on a clean dist/ and returns that directory; dist/ is removed at the end.
func buildDist(t *testing.T) string {
	t.Helper()
	root := repoRoot(t)
	dist := filepath.Join(root, "dist")
	_ = os.RemoveAll(dist)
	t.Cleanup(func() { _ = os.RemoveAll(dist) })
	cmd := exec.Command("sh", filepath.Join(root, "scripts", "release-build.sh"))
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("release-build.sh error: %v\n%s", err, out)
	}
	return dist
}

// END: buildDist

// START: TestReleaseBuildFiveAssets

func TestReleaseBuildFiveAssets(t *testing.T) {
	entries, err := os.ReadDir(buildDist(t))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range entries {
		got = append(got, e.Name())
	}
	sort.Strings(got)
	want := []string{"goremi_darwin_amd64", "goremi_darwin_arm64", "goremi_linux_amd64", "goremi_linux_arm64", "goremi_windows_amd64.exe"}
	if len(got) != len(want) {
		t.Fatalf("dist/ = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("dist/ = %q, want %q", got, want)
		}
	}
}

// END: TestReleaseBuildFiveAssets

// START: TestReleaseBuildHeaders

func TestReleaseBuildHeaders(t *testing.T) {
	dist := buildDist(t)
	check := func(name string, ok func(path string) (bool, error)) {
		good, err := ok(filepath.Join(dist, name))
		if err != nil || !good {
			t.Errorf("%s has the wrong header (ok=%v, err=%v)", name, good, err)
		}
	}
	check("goremi_linux_amd64", func(p string) (bool, error) {
		f, err := elf.Open(p)
		if err != nil {
			return false, err
		}
		defer f.Close()
		return f.Machine == elf.EM_X86_64, nil
	})
	check("goremi_linux_arm64", func(p string) (bool, error) {
		f, err := elf.Open(p)
		if err != nil {
			return false, err
		}
		defer f.Close()
		return f.Machine == elf.EM_AARCH64, nil
	})
	check("goremi_darwin_amd64", func(p string) (bool, error) {
		f, err := macho.Open(p)
		if err != nil {
			return false, err
		}
		defer f.Close()
		return f.Cpu == macho.CpuAmd64, nil
	})
	check("goremi_darwin_arm64", func(p string) (bool, error) {
		f, err := macho.Open(p)
		if err != nil {
			return false, err
		}
		defer f.Close()
		return f.Cpu == macho.CpuArm64, nil
	})
	check("goremi_windows_amd64.exe", func(p string) (bool, error) {
		f, err := pe.Open(p)
		if err != nil {
			return false, err
		}
		defer f.Close()
		return f.Machine == pe.IMAGE_FILE_MACHINE_AMD64, nil
	})
}

// END: TestReleaseBuildHeaders

// START: TestDistIgnored

func TestDistIgnored(t *testing.T) {
	cmd := exec.Command("git", "check-ignore", "-q", "dist/goremi_linux_amd64")
	cmd.Dir = repoRoot(t)
	if err := cmd.Run(); err != nil {
		t.Fatalf("git check-ignore dist/goremi_linux_amd64: %v, want exit 0 (dist/ in .gitignore)", err)
	}
}

// END: TestDistIgnored

// START: TestLinuxBinariesStatic

func TestLinuxBinariesStatic(t *testing.T) {
	dist := buildDist(t)
	for _, name := range []string{"goremi_linux_amd64", "goremi_linux_arm64"} {
		f, err := elf.Open(filepath.Join(dist, name))
		if err != nil {
			t.Fatal(err)
		}
		libs, _ := f.ImportedLibraries()
		if f.Section(".interp") != nil || len(libs) > 0 {
			t.Errorf("%s is dynamically linked (interpreter section: %v, libraries: %v), want a static binary", name, f.Section(".interp") != nil, libs)
		}
		f.Close()
	}
}

// END: TestLinuxBinariesStatic
