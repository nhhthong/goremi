//go:build !windows

// Tests of install.sh offering to install mpv and yt-dlp: the list, the one question, the answer, the package managers, the yt-dlp binary and its checksum (tasks 20.23 to 20.32 and 20.41).
package install

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// START: fakes

// fakePM is a package manager that logs its call and exits with PM_EXIT (0 when unset).
func fakePM(name string) string {
	return "#!/bin/sh\nprintf '" + name + " %s\\n' \"$*\" >> \"$PM_LOG\"\nexit ${PM_EXIT:-0}\n"
}

// fakeSudoRuns logs the call and runs the command it was given, so the fake package manager logs it too.
const fakeSudoRuns = "#!/bin/sh\nprintf 'sudo %s\\n' \"$*\" >> \"$SUDO_LOG\"\nexec \"$@\"\n"

// fakeCurlBodies serves the goremi binary, the yt-dlp binary and its SHA2-256SUMS from three files, and fails on request.
const fakeCurlBodies = `#!/bin/sh
out=""
url=""
while [ $# -gt 0 ]; do
	case "$1" in
	-o) out="$2"; shift 2 ;;
	-*) shift ;;
	*) url="$1"; shift ;;
	esac
done
printf '%s\n' "$url" >> "$CURL_LOG"
case "$url" in
*SHA2-256SUMS) cp "$FAKE_SUMS" "$out" ;;
*/yt-dlp_*) [ -n "$FAKE_YTDLP_FAIL" ] && exit 22; cp "$FAKE_YTDLP" "$out" ;;
*) cp "$FAKE_BODY" "$out" ;;
esac
`

const ytBody = "#!/bin/sh\necho yt-dlp\n"

// ytAsset is the yt-dlp release file for a system (install.md 2026-10-06).
func ytAsset(osName, arch string) string {
	switch {
	case osName == "Darwin":
		return "yt-dlp_macos"
	case arch == "aarch64" || arch == "arm64":
		return "yt-dlp_linux_aarch64"
	}
	return "yt-dlp_linux"
}

// withDeps runs install.sh with the given package managers on PATH, a sudo that runs its command, and a curl that serves the yt-dlp binary with a SHA2-256SUMS that matches it, or not.
func withDeps(t *testing.T, r Run, goodSums bool, pms ...string) Result {
	t.Helper()
	tmp := t.TempDir()
	hash := sha256.Sum256([]byte(ytBody))
	sum := hex.EncodeToString(hash[:])
	if !goodSums {
		sum = strings.Repeat("0", 64)
	}
	files := map[string]string{"yt": ytBody, "sums": sum + "  " + ytAsset(r.OS, r.Arch) + "\n"}
	for name, text := range files {
		if err := os.WriteFile(filepath.Join(tmp, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	tools := map[string]string{"sudo": fakeSudoRuns, "curl": fakeCurlBodies}
	for _, pm := range pms {
		tools[pm] = fakePM(pm)
	}
	for name, text := range r.Tools {
		tools[name] = text
	}
	r.Tools = tools
	r.Env = append(r.Env, "FAKE_YTDLP="+filepath.Join(tmp, "yt"), "FAKE_SUMS="+filepath.Join(tmp, "sums"))
	return runInstall(t, r, "")
}

func answer(s string) *string { return &s }

// installed tells whether a file exists under the install directory of the result.
func installed(res Result, name string) bool {
	_, err := os.Stat(filepath.Join(res.Home, ".local", "bin", name))
	return err == nil
}

// END: fakes

// START: TestDepsPresentAsksNothing

func TestDepsPresentAsksNothing(t *testing.T) {
	res := withDeps(t, Run{OS: "Linux", Arch: "x86_64", Mpv: true, YtDlp: true, Answer: answer("y\n")}, true, "apt")
	if strings.Contains(res.Stdout, "Install") || strings.Contains(res.Stdout, "mpv") || strings.Contains(res.Stdout, "yt-dlp") {
		t.Fatalf("stdout = %q, want nothing about dependencies", res.Stdout)
	}
	if len(res.PMCalls) != 0 || len(res.CurlURLs) != 1 {
		t.Fatalf("package manager calls %v, curl URLs %v; want none and only the goremi download", res.PMCalls, res.CurlURLs)
	}
}

// END: TestDepsPresentAsksNothing

// START: TestMissingListsCommandsThenAsks

func TestMissingListsCommandsThenAsks(t *testing.T) {
	res := withDeps(t, Run{OS: "Linux", Arch: "x86_64"}, true, "apt")
	out := res.Stdout
	question := "Install mpv and yt-dlp? [y/N]"
	if strings.Count(out, question) != 1 {
		t.Fatalf("stdout = %q, want the question %q once", out, question)
	}
	for _, want := range []string{"sudo apt install -y mpv", "~/.local/bin/yt-dlp", "yt-dlp_linux"} {
		if i := strings.Index(out, want); i < 0 || i > strings.Index(out, question) {
			t.Errorf("stdout = %q, want %q before the question", out, want)
		}
	}
}

// END: TestMissingListsCommandsThenAsks

// START: TestAnswerComesFromTerminalNotStdin

func TestAnswerComesFromTerminalNotStdin(t *testing.T) {
	res := withDeps(t, Run{OS: "Linux", Arch: "x86_64", YtDlp: true, Stdin: true, Answer: answer("y\n")}, true, "apt")
	if !has(res.PMCalls, "apt install -y mpv") {
		t.Fatalf("package manager calls = %v, want apt install -y mpv (the answer y came from the terminal while the script came from stdin)", res.PMCalls)
	}
}

// END: TestAnswerComesFromTerminalNotStdin

// START: no install

// noInstall checks a run whose answer is not yes: nothing installed, the list still printed, exit 0.
func noInstall(t *testing.T, ans string) {
	t.Helper()
	res := withDeps(t, Run{OS: "Linux", Arch: "x86_64", Answer: answer(ans)}, true, "apt")
	if res.Exit != 0 || len(res.PMCalls) != 0 || installed(res, "yt-dlp") {
		t.Fatalf("exit %d, package manager calls %v, yt-dlp installed %v; want exit 0 and nothing installed", res.Exit, res.PMCalls, installed(res, "yt-dlp"))
	}
	if !strings.Contains(res.Stdout, "sudo apt install -y mpv") {
		t.Fatalf("stdout = %q, want the list of commands", res.Stdout)
	}
}

func TestAnswerNInstallsNothing(t *testing.T)     { noInstall(t, "n\n") }
func TestEmptyAnswerInstallsNothing(t *testing.T) { noInstall(t, "\n") }
func TestEndOfInputInstallsNothing(t *testing.T)  { noInstall(t, "") }

// END: no install

// START: mpv by package manager

// wantMpvCall checks a yes run with one package manager: the call it got and who ran it.
func wantMpvCall(t *testing.T, r Run, pm, call string, viaSudo bool) {
	t.Helper()
	r.Mpv, r.YtDlp = false, true
	r.Answer = answer("y\n")
	res := withDeps(t, r, true, pm)
	if !has(res.PMCalls, call) {
		t.Fatalf("package manager calls = %v, want %q", res.PMCalls, call)
	}
	if viaSudo && !strings.Contains(strings.Join(res.SudoCalls, " "), "sudo "+call) {
		t.Fatalf("sudo calls = %v, want %q", res.SudoCalls, "sudo "+call)
	}
	if !viaSudo && len(res.SudoCalls) != 0 {
		t.Fatalf("sudo calls = %v, want none for %s", res.SudoCalls, pm)
	}
}

func TestYesInstallsMpvWithApt(t *testing.T) {
	wantMpvCall(t, Run{OS: "Linux", Arch: "x86_64"}, "apt", "apt install -y mpv", true)
}
func TestYesInstallsMpvWithDnf(t *testing.T) {
	wantMpvCall(t, Run{OS: "Linux", Arch: "x86_64"}, "dnf", "dnf install -y mpv", true)
}
func TestYesInstallsMpvWithPacman(t *testing.T) {
	wantMpvCall(t, Run{OS: "Linux", Arch: "x86_64"}, "pacman", "pacman -S --noconfirm mpv", true)
}
func TestYesInstallsMpvWithBrewNoSudo(t *testing.T) {
	wantMpvCall(t, Run{OS: "Darwin", Arch: "arm64"}, "brew", "brew install mpv", false)
}

// END: mpv by package manager

// START: TestAnswerIsNeverRun

func TestAnswerIsNeverRun(t *testing.T) {
	res := withDeps(t, Run{OS: "Linux", Arch: "x86_64", Answer: answer("y; touch \"$HOME/pwned\"\n")}, true, "apt")
	if _, err := os.Stat(filepath.Join(res.Home, "pwned")); err == nil {
		t.Fatal("the answer was run as a command: $HOME/pwned exists")
	}
	if len(res.PMCalls) != 0 {
		t.Fatalf("package manager calls = %v, want none: the answer is not exactly y", res.PMCalls)
	}
}

// END: TestAnswerIsNeverRun

// START: yt-dlp binary

// wantYtdlp checks a yes run that installs only yt-dlp: the URL asked for, then the file.
func wantYtdlp(t *testing.T, osName, arch string) Result {
	t.Helper()
	res := withDeps(t, Run{OS: osName, Arch: arch, Mpv: true, Answer: answer("y\n")}, true)
	want := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/" + ytAsset(osName, arch)
	if !has(res.CurlURLs, want) {
		t.Fatalf("curl URLs = %v, want %s", res.CurlURLs, want)
	}
	return res
}

func TestYtdlpBinaryLinuxAmd64(t *testing.T) {
	res := wantYtdlp(t, "Linux", "x86_64")
	info, err := os.Stat(filepath.Join(res.Home, ".local", "bin", "yt-dlp"))
	if err != nil || info.Mode().Perm() != 0o755 {
		t.Fatalf("yt-dlp file: %v, %v; want mode 755", info, err)
	}
	if got, _ := os.ReadFile(filepath.Join(res.Home, ".local", "bin", "yt-dlp")); string(got) != ytBody {
		t.Fatalf("yt-dlp content = %q, want the served binary", got)
	}
}
func TestYtdlpBinaryLinuxArm64(t *testing.T) { wantYtdlp(t, "Linux", "aarch64") }
func TestYtdlpBinaryMacos(t *testing.T)      { wantYtdlp(t, "Darwin", "arm64") }

func TestYtdlpHttpsOnly(t *testing.T) {
	res := wantYtdlp(t, "Linux", "x86_64")
	for _, u := range res.CurlURLs {
		if !strings.HasPrefix(u, "https://") {
			t.Errorf("URL %q does not start with https://", u)
		}
	}
}

// END: yt-dlp binary

// START: TestYesFlagInstalls

func TestYesFlagInstalls(t *testing.T) {
	res := withDeps(t, Run{OS: "Linux", Arch: "x86_64", Stdin: true, Args: []string{"-y"}}, true, "apt")
	if !has(res.PMCalls, "apt install -y mpv") || !installed(res, "yt-dlp") {
		t.Fatalf("package manager calls %v, yt-dlp installed %v; want both installed with no answer at all", res.PMCalls, installed(res, "yt-dlp"))
	}
}

// END: TestYesFlagInstalls

// START: TestNoPackageManagerPrintsCommand

func TestNoPackageManagerPrintsCommand(t *testing.T) {
	res := withDeps(t, Run{OS: "Linux", Arch: "x86_64", YtDlp: true, Answer: answer("y\n")}, true)
	if res.Exit != 0 || len(res.SudoCalls) != 0 || len(res.PMCalls) != 0 {
		t.Fatalf("exit %d, sudo calls %v, package manager calls %v; want exit 0 and no call", res.Exit, res.SudoCalls, res.PMCalls)
	}
	if !strings.Contains(res.Stdout, "sudo apt install -y mpv") || !installed(res, "goremi") {
		t.Fatalf("stdout = %q, goremi installed %v; want the command printed and goremi installed", res.Stdout, installed(res, "goremi"))
	}
}

// END: TestNoPackageManagerPrintsCommand

// START: no terminal

func TestNoTerminalPrintsListOnly(t *testing.T) {
	res := withDeps(t, Run{OS: "Linux", Arch: "x86_64"}, true, "apt")
	if res.Exit != 0 || len(res.PMCalls) != 0 || installed(res, "yt-dlp") {
		t.Fatalf("exit %d, calls %v; want exit 0 and nothing installed", res.Exit, res.PMCalls)
	}
	if !strings.Contains(res.Stdout, "sudo apt install -y mpv") || !strings.Contains(res.Stdout, "-y") {
		t.Fatalf("stdout = %q, want the list and a hint to run again with -y", res.Stdout)
	}
}

func TestNoTerminalThenTerminalInstalls(t *testing.T) {
	first := withDeps(t, Run{OS: "Linux", Arch: "x86_64"}, true, "apt")
	if len(first.PMCalls) != 0 {
		t.Fatalf("first run (no terminal) called %v, want nothing", first.PMCalls)
	}
	second := withDeps(t, Run{OS: "Linux", Arch: "x86_64", Answer: answer("y\n")}, true, "apt")
	if !has(second.PMCalls, "apt install -y mpv") || !installed(second, "yt-dlp") {
		t.Fatalf("second run (terminal back) called %v, yt-dlp installed %v; want it installed", second.PMCalls, installed(second, "yt-dlp"))
	}
}

// END: no terminal

// START: checksum

func TestChecksumMatchInstalls(t *testing.T) {
	res := wantYtdlp(t, "Linux", "x86_64")
	sums := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/SHA2-256SUMS"
	if !has(res.CurlURLs, sums) || !installed(res, "yt-dlp") {
		t.Fatalf("curl URLs %v, yt-dlp installed %v; want the sums file asked for and yt-dlp installed", res.CurlURLs, installed(res, "yt-dlp"))
	}
}

func TestChecksumMismatchRefused(t *testing.T) {
	res := withDeps(t, Run{OS: "Linux", Arch: "x86_64", Mpv: true, Answer: answer("y\n")}, false)
	if installed(res, "yt-dlp") || !strings.Contains(strings.ToLower(res.Stderr), "checksum") || !installed(res, "goremi") {
		t.Fatalf("yt-dlp installed %v, stderr %q, goremi installed %v; want no yt-dlp, a checksum error and goremi installed", installed(res, "yt-dlp"), res.Stderr, installed(res, "goremi"))
	}
}

func TestChecksumMismatchLeavesNoPartial(t *testing.T) {
	res := withDeps(t, Run{OS: "Linux", Arch: "x86_64", Mpv: true, Answer: answer("y\n")}, false)
	entries, _ := os.ReadDir(filepath.Join(res.Home, ".local", "bin"))
	if len(entries) != 1 || entries[0].Name() != "goremi" {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("install directory holds %v, want only goremi", names)
	}
}

// END: checksum

// START: failed dependency

func TestFailedMpvInstallKeepsGoremi(t *testing.T) {
	r := Run{OS: "Linux", Arch: "x86_64", Answer: answer("y\n"), Env: []string{"PM_EXIT=1"}}
	res := withDeps(t, r, true, "apt")
	if !installed(res, "goremi") || !strings.Contains(res.Stderr, "mpv") || !installed(res, "yt-dlp") {
		t.Fatalf("goremi %v, yt-dlp %v, stderr %q; want goremi and yt-dlp installed and an error naming mpv", installed(res, "goremi"), installed(res, "yt-dlp"), res.Stderr)
	}
}

func TestFailedYtdlpDownloadKeepsGoremi(t *testing.T) {
	r := Run{OS: "Linux", Arch: "x86_64", Mpv: true, Answer: answer("y\n"), Env: []string{"FAKE_YTDLP_FAIL=1"}}
	res := withDeps(t, r, true)
	entries, _ := os.ReadDir(filepath.Join(res.Home, ".local", "bin"))
	if !installed(res, "goremi") || !strings.Contains(res.Stderr, "yt-dlp") || len(entries) != 1 {
		t.Fatalf("goremi %v, stderr %q, %d files; want goremi only and an error naming yt-dlp", installed(res, "goremi"), res.Stderr, len(entries))
	}
}

// END: failed dependency
