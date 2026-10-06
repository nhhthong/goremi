//go:build !windows

// Tests of install.sh: the asset it asks for, unsupported systems, where it writes, and what a failed download leaves.
package install

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const baseURL = "https://github.com/nhhthong/goremi/releases/latest/download/"

// START: asset names

// wantAsset runs the installer for one OS and architecture and checks the file that curl was asked for.
func wantAsset(t *testing.T, osName, arch, asset string) {
	t.Helper()
	res := runInstall(t, Run{OS: osName, Arch: arch}, "")
	if res.Exit != 0 {
		t.Fatalf("exit = %d, stderr = %q", res.Exit, res.Stderr)
	}
	if len(res.CurlURLs) != 1 || !strings.HasSuffix(res.CurlURLs[0], "/"+asset) {
		t.Fatalf("curl asked for %q, want one URL ending /%s", res.CurlURLs, asset)
	}
}

func TestAssetLinuxX86(t *testing.T)     { wantAsset(t, "Linux", "x86_64", "goremi_linux_amd64") }
func TestAssetLinuxAarch64(t *testing.T) { wantAsset(t, "Linux", "aarch64", "goremi_linux_arm64") }
func TestAssetLinuxArm64(t *testing.T)   { wantAsset(t, "Linux", "arm64", "goremi_linux_arm64") }
func TestAssetDarwinX86(t *testing.T)    { wantAsset(t, "Darwin", "x86_64", "goremi_darwin_amd64") }
func TestAssetDarwinArm64(t *testing.T)  { wantAsset(t, "Darwin", "arm64", "goremi_darwin_arm64") }

// END: asset names

// START: unsupported systems

// wantRefused checks an unsupported system: exit 1, the names in stderr, no download, nothing installed.
func wantRefused(t *testing.T, osName, arch string) {
	t.Helper()
	res := runInstall(t, Run{OS: osName, Arch: arch}, "")
	if res.Exit != 1 {
		t.Fatalf("exit = %d, want 1 (stderr %q)", res.Exit, res.Stderr)
	}
	if !strings.Contains(res.Stderr, osName) || !strings.Contains(res.Stderr, arch) {
		t.Fatalf("stderr = %q, want it to name %s and %s", res.Stderr, osName, arch)
	}
	if len(res.CurlURLs) != 0 {
		t.Fatalf("curl was called: %q", res.CurlURLs)
	}
	if _, err := os.Stat(filepath.Join(res.Home, ".local", "bin")); err == nil {
		t.Fatal("~/.local/bin exists after a refused run")
	}
}

func TestUnsupportedOS(t *testing.T)   { wantRefused(t, "FreeBSD", "amd64") }
func TestUnsupportedArch(t *testing.T) { wantRefused(t, "Linux", "riscv64") }

// END: unsupported systems

// START: download and file

func TestDownloadUrlAndFile(t *testing.T) {
	body := "#!/bin/sh\necho known body\n"
	res := runInstall(t, Run{OS: "Linux", Arch: "x86_64", Body: body}, "")
	if res.Exit != 0 {
		t.Fatalf("exit = %d, stderr = %q", res.Exit, res.Stderr)
	}
	if want := baseURL + "goremi_linux_amd64"; len(res.CurlURLs) != 1 || res.CurlURLs[0] != want {
		t.Fatalf("URLs = %q, want [%s]", res.CurlURLs, want)
	}
	got, err := os.ReadFile(filepath.Join(res.Home, ".local", "bin", "goremi"))
	if err != nil || string(got) != body {
		t.Fatalf("installed file = %q, %v; want the served body", got, err)
	}
}

func TestInstalledMode755(t *testing.T) {
	res := runInstall(t, Run{OS: "Linux", Arch: "x86_64"}, "")
	info, err := os.Stat(filepath.Join(res.Home, ".local", "bin", "goremi"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Fatalf("mode = %v, want 0755", info.Mode().Perm())
	}
}

func TestCreatesInstallDir(t *testing.T) {
	res := runInstall(t, Run{OS: "Linux", Arch: "x86_64"}, "")
	if res.Exit != 0 {
		t.Fatalf("exit = %d, stderr = %q", res.Exit, res.Stderr)
	}
	entries, err := os.ReadDir(filepath.Join(res.Home, ".local", "bin"))
	if err != nil || len(entries) != 1 || entries[0].Name() != "goremi" {
		t.Fatalf("~/.local/bin = %v, %v; want only goremi", entries, err)
	}
}

// END: download and file

// START: security

func TestHostileUnameNotExecuted(t *testing.T) {
	res := runInstall(t, Run{OS: "Linux", Arch: `x86_64; touch "$HOME/pwned"`}, "")
	if res.Exit != 1 {
		t.Fatalf("exit = %d, want 1", res.Exit)
	}
	if _, err := os.Stat(filepath.Join(res.Home, "pwned")); err == nil {
		t.Fatal("the output of uname was run as a command: pwned exists")
	}
	if len(res.CurlURLs) != 0 {
		t.Fatalf("curl was called: %q", res.CurlURLs)
	}
}

func TestHttpsOnlyNoSudo(t *testing.T) {
	res := runInstall(t, Run{OS: "Linux", Arch: "x86_64"}, "")
	if len(res.CurlURLs) == 0 {
		t.Fatal("curl was not called")
	}
	for _, u := range res.CurlURLs {
		if !strings.HasPrefix(u, "https://") {
			t.Fatalf("URL %q does not start with https://", u)
		}
	}
	if len(res.SudoCalls) != 0 {
		t.Fatalf("sudo was called: %q", res.SudoCalls)
	}
}

func TestWritesOnlyInstallDir(t *testing.T) {
	res := runInstall(t, Run{OS: "Linux", Arch: "x86_64"}, "")
	var files []string
	_ = filepath.WalkDir(res.Home, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(res.Home, p)
			files = append(files, rel)
		}
		return nil
	})
	if len(files) != 1 || files[0] != filepath.Join(".local", "bin", "goremi") {
		t.Fatalf("files under HOME = %q, want only .local/bin/goremi", files)
	}
}

// END: security

// START: failed download

func TestDownloadFailsNoFile(t *testing.T) {
	res := runInstall(t, Run{OS: "Linux", Arch: "x86_64", Fail: true}, "")
	if res.Exit != 1 {
		t.Fatalf("exit = %d, want 1", res.Exit)
	}
	if strings.TrimSpace(res.Stderr) == "" {
		t.Fatal("stderr is empty, want a message")
	}
	if _, err := os.Stat(filepath.Join(res.Home, ".local", "bin", "goremi")); err == nil {
		t.Fatal("goremi exists after a failed download")
	}
}

func TestDownloadFailsKeepsOld(t *testing.T) {
	home := filepath.Join(t.TempDir(), "home")
	bin := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(bin, "goremi")
	if err := os.WriteFile(old, []byte("old version"), 0o755); err != nil {
		t.Fatal(err)
	}
	res := runInstall(t, Run{OS: "Linux", Arch: "x86_64", Fail: true}, home)
	if res.Exit != 1 {
		t.Fatalf("exit = %d, want 1", res.Exit)
	}
	got, err := os.ReadFile(old)
	info, _ := os.Stat(old)
	if err != nil || string(got) != "old version" || info.Mode().Perm() != 0o755 {
		t.Fatalf("old goremi = %q (%v), %v; want it untouched", got, info.Mode(), err)
	}
}

func TestDownloadFailsNoPartial(t *testing.T) {
	res := runInstall(t, Run{OS: "Linux", Arch: "x86_64", Fail: true, Partial: true}, "")
	if res.Exit != 1 {
		t.Fatalf("exit = %d, want 1", res.Exit)
	}
	entries, _ := os.ReadDir(filepath.Join(res.Home, ".local", "bin"))
	for _, e := range entries {
		t.Errorf("leftover in ~/.local/bin: %s", e.Name())
	}
}

// END: failed download

// START: reinstall

func TestReinstallReplaces(t *testing.T) {
	first := runInstall(t, Run{OS: "Linux", Arch: "x86_64", Body: "version one\n"}, "")
	if first.Exit != 0 {
		t.Fatalf("first run exit = %d, stderr = %q", first.Exit, first.Stderr)
	}
	second := runInstall(t, Run{OS: "Linux", Arch: "x86_64", Body: "version two\n"}, first.Home)
	if second.Exit != 0 {
		t.Fatalf("second run exit = %d, stderr = %q", second.Exit, second.Stderr)
	}
	bin := filepath.Join(first.Home, ".local", "bin")
	got, err := os.ReadFile(filepath.Join(bin, "goremi"))
	if err != nil || string(got) != "version two\n" {
		t.Fatalf("goremi = %q, %v; want the second body", got, err)
	}
	if entries, _ := os.ReadDir(bin); len(entries) != 1 {
		t.Fatalf("~/.local/bin has %d entries, want only goremi", len(entries))
	}
}

// END: reinstall

// START: dependency report

// reported runs a successful install and returns its stdout.
func reported(t *testing.T, r Run) string {
	t.Helper()
	res := runInstall(t, r, "")
	if res.Exit != 0 {
		t.Fatalf("exit = %d, want 0 (stderr %q)", res.Exit, res.Stderr)
	}
	return res.Stdout
}

func TestDepsBothMissing(t *testing.T) {
	out := reported(t, Run{OS: "Linux", Arch: "x86_64"})
	m, y := strings.Index(out, "mpv"), strings.Index(out, "yt-dlp")
	if m < 0 || y < 0 || m > y {
		t.Fatalf("stdout = %q, want a line for mpv, then one for yt-dlp", out)
	}
}

func TestDepsOnlyMpvMissing(t *testing.T) {
	out := reported(t, Run{OS: "Linux", Arch: "x86_64", YtDlp: true})
	if !strings.Contains(out, "mpv") || strings.Contains(out, "yt-dlp") {
		t.Fatalf("stdout = %q, want mpv only", out)
	}
}

func TestDepsOnlyYtdlpMissing(t *testing.T) {
	out := reported(t, Run{OS: "Linux", Arch: "x86_64", Mpv: true})
	if !strings.Contains(out, "yt-dlp") || strings.Contains(out, "mpv") {
		t.Fatalf("stdout = %q, want yt-dlp only", out)
	}
}

func TestDepsNoneMissing(t *testing.T) {
	out := reported(t, Run{OS: "Linux", Arch: "x86_64", Mpv: true, YtDlp: true})
	if strings.Contains(out, "mpv") || strings.Contains(out, "yt-dlp") {
		t.Fatalf("stdout = %q, want neither name", out)
	}
}

func TestMacHintsBoth(t *testing.T) {
	out := reported(t, Run{OS: "Darwin", Arch: "arm64"})
	if !strings.Contains(out, "brew install mpv") || !strings.Contains(out, "brew install yt-dlp") {
		t.Fatalf("stdout = %q, want both brew commands", out)
	}
}

func TestMacHintOnlyMissing(t *testing.T) {
	out := reported(t, Run{OS: "Darwin", Arch: "arm64", Mpv: true})
	if !strings.Contains(out, "brew install yt-dlp") || strings.Contains(out, "brew install mpv") {
		t.Fatalf("stdout = %q, want only the yt-dlp command", out)
	}
}

func TestLinuxNamesOnly(t *testing.T) {
	out := reported(t, Run{OS: "Linux", Arch: "x86_64"})
	if !strings.Contains(out, "mpv") || !strings.Contains(out, "yt-dlp") || strings.Contains(out, "brew") {
		t.Fatalf("stdout = %q, want both names and no brew", out)
	}
}

// END: dependency report

// START: piped install

func TestPipedInstallRuns(t *testing.T) {
	res := runInstall(t, Run{OS: "Linux", Arch: "x86_64", Stdin: true, Body: "#!/bin/sh\necho goremi runs\n"}, "")
	if res.Exit != 0 {
		t.Fatalf("exit = %d, stderr = %q", res.Exit, res.Stderr)
	}
	out, err := exec.Command(filepath.Join(res.Home, ".local", "bin", "goremi")).Output()
	if err != nil || strings.TrimSpace(string(out)) != "goremi runs" {
		t.Fatalf("goremi output = %q, %v; want \"goremi runs\"", out, err)
	}
}

func TestPipedInstallFailure(t *testing.T) {
	res := runInstall(t, Run{OS: "Linux", Arch: "x86_64", Stdin: true, Fail: true}, "")
	if res.Exit != 1 {
		t.Fatalf("exit = %d, want 1", res.Exit)
	}
	if _, err := os.Stat(filepath.Join(res.Home, ".local", "bin", "goremi")); err == nil {
		t.Fatal("goremi exists after a failed piped install")
	}
}

// END: piped install

// START: truncated script

func TestTruncatedScriptRunsNothing(t *testing.T) {
	text, err := os.ReadFile(filepath.Join(repoRoot(t), "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	body := strings.TrimRight(string(text), "\n")
	lastLine := strings.LastIndex(body, "\n") + 1 // every cut before the last line
	sealed := sealedPath(t, false, false)
	for cut := 1; cut < lastLine; cut++ {
		res := runInstall(t, Run{OS: "Linux", Arch: "x86_64", Stdin: true, Cut: cut, Sealed: sealed}, "")
		written, _ := os.ReadDir(res.Home)
		if len(res.CurlURLs) != 0 || len(written) != 0 {
			t.Fatalf("script cut at byte %d (after %q): curl calls %q, %d entries under HOME; want nothing run", cut, body[max(0, cut-30):cut], res.CurlURLs, len(written))
		}
	}
}

// END: truncated script

// START: hints

func TestPathHintWhenMissing(t *testing.T) {
	out := reported(t, Run{OS: "Linux", Arch: "x86_64", Mpv: true, YtDlp: true})
	if !strings.Contains(out, "~/.local/bin") || !strings.Contains(out, "PATH") {
		t.Fatalf("stdout = %q, want a line naming ~/.local/bin and PATH", out)
	}
}

func TestNoPathHintWhenPresent(t *testing.T) {
	out := reported(t, Run{OS: "Linux", Arch: "x86_64", Mpv: true, YtDlp: true, PathHome: true})
	if strings.Contains(out, "PATH") {
		t.Fatalf("stdout = %q, want no PATH line", out)
	}
}

func TestLinuxMpvAptHint(t *testing.T) {
	out := reported(t, Run{OS: "Linux", Arch: "x86_64", YtDlp: true})
	if !strings.Contains(out, "sudo apt install mpv") {
		t.Fatalf("stdout = %q, want sudo apt install mpv", out)
	}
}

func TestLinuxYtdlpBinaryHint(t *testing.T) {
	out := reported(t, Run{OS: "Linux", Arch: "x86_64", Mpv: true})
	if !strings.Contains(out, "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp") || !strings.Contains(out, "~/.local/bin/yt-dlp") {
		t.Fatalf("stdout = %q, want the yt-dlp release binary command", out)
	}
}

// END: hints
