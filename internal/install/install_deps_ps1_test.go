// Tests of install.ps1 offering to install mpv and yt-dlp: the list, the one question, the answer, winget and scoop, the yt-dlp binary and its checksum (tasks 20.33 to 20.41).
package install

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ps1DepsWrapper fakes the download, the tool lookup, the question, winget and scoop around install.ps1; -Yes is passed when FAKE_YES is set.
const ps1DepsWrapper = `function Invoke-WebRequest { param($Uri, $OutFile) Add-Content -LiteralPath $env:WEB_LOG -Value $Uri; if ($Uri -like '*SHA2-256SUMS') { Copy-Item -LiteralPath $env:FAKE_SUMS -Destination $OutFile } elseif ($Uri -like '*yt-dlp.exe') { if ($env:FAKE_YTDLP_FAIL) { throw 'network down' }; Copy-Item -LiteralPath $env:FAKE_YTDLP -Destination $OutFile } else { Copy-Item -LiteralPath $env:FAKE_BODY -Destination $OutFile } }
function Get-Command { if (($env:FAKE_HAVE -split ',') -contains $args[0]) { [pscustomobject]@{ Name = $args[0] } } }
function Read-Host { param($Prompt) Add-Content -LiteralPath $env:PROMPT_LOG -Value $Prompt; if ($env:FAKE_NO_TTY) { throw 'no console' }; $env:FAKE_ANSWER }
function winget { Add-Content -LiteralPath $env:PM_LOG -Value "winget $args"; $global:LASTEXITCODE = [int]$env:PM_EXIT }
function scoop { Add-Content -LiteralPath $env:PM_LOG -Value "scoop $args"; $global:LASTEXITCODE = [int]$env:PM_EXIT }
if ($env:FAKE_YES) { & $env:INSTALL_PS1 -Yes } else { & $env:INSTALL_PS1 }
exit $LASTEXITCODE`

const ps1YtBody = "fake yt-dlp exe"

// DepsRun says how one run looks: the tools found, the answer and what fails.
type DepsRun struct {
	Have      string // comma-separated tool names Get-Command finds
	Answer    string // what Read-Host returns
	NoTTY     bool   // Read-Host throws, as in a non-interactive host
	Yes       bool   // the script is run with -Yes
	BadSums   bool   // SHA2-256SUMS holds a wrong hash
	YtdlpFail bool   // the yt-dlp.exe download throws
	PMExit    string // exit code of winget and scoop, empty = 0
}

// DepsResult is what a run left behind.
type DepsResult struct {
	Dir, Stdout, Stderr string
	Exit                int
	URLs, PM, Prompts   []string
}

// runDeps runs install.ps1 under pwsh with the fakes above; a missing pwsh fails the test.
func runDeps(t *testing.T, r DepsRun) DepsResult {
	t.Helper()
	pwsh, err := exec.LookPath("pwsh")
	if err != nil {
		t.Fatalf("pwsh not found: %v", err)
	}
	tmp := t.TempDir()
	local := filepath.Join(tmp, "LocalAppData")
	if err := os.MkdirAll(local, 0o755); err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256([]byte(ps1YtBody))
	sum := hex.EncodeToString(hash[:])
	if r.BadSums {
		sum = strings.Repeat("0", 64)
	}
	for name, text := range map[string]string{"body": "fake exe body", "yt": ps1YtBody, "sums": sum + "  yt-dlp.exe\n"} {
		if err := os.WriteFile(filepath.Join(tmp, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	script, _ := filepath.Abs(filepath.Join("..", "..", "install.ps1"))
	logs := map[string]string{"web": "WEB_LOG", "pm": "PM_LOG", "prompt": "PROMPT_LOG"}
	pairs := []string{"DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=1", "LocalAppData=" + local, "PROCESSOR_ARCHITECTURE=AMD64", "PROCESSOR_ARCHITEW6432=", "INSTALL_PS1=" + script,
		"FAKE_HAVE=" + r.Have, "FAKE_ANSWER=" + r.Answer, "FAKE_BODY=" + filepath.Join(tmp, "body"), "FAKE_YTDLP=" + filepath.Join(tmp, "yt"), "FAKE_SUMS=" + filepath.Join(tmp, "sums"),
		"FAKE_NO_TTY=", "FAKE_YES=", "FAKE_YTDLP_FAIL=", "PM_EXIT=" + r.PMExit}
	for name, env := range logs {
		pairs = append(pairs, env+"="+filepath.Join(tmp, name+".log"))
	}
	if r.NoTTY {
		pairs = append(pairs, "FAKE_NO_TTY=1")
	}
	if r.Yes {
		pairs = append(pairs, "FAKE_YES=1")
	}
	if r.YtdlpFail {
		pairs = append(pairs, "FAKE_YTDLP_FAIL=1")
	}
	cmd := exec.Command(pwsh, "-NoProfile", "-NonInteractive", "-Command", ps1DepsWrapper)
	cmd.Dir = tmp
	cmd.Env = withEnv(os.Environ(), pairs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	res := DepsResult{Dir: filepath.Join(local, "Programs", "goremi")}
	if err := cmd.Run(); err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("pwsh did not run: %v", err)
		}
		res.Exit = ee.ExitCode()
	}
	res.Stdout, res.Stderr = stdout.String(), stderr.String()
	res.URLs, res.PM, res.Prompts = logLines(filepath.Join(tmp, "web.log")), logLines(filepath.Join(tmp, "pm.log")), logLines(filepath.Join(tmp, "prompt.log"))
	return res
}

// exists tells whether a file exists under the install directory of the result.
func (r DepsResult) exists(name string) bool {
	_, err := os.Stat(filepath.Join(r.Dir, name))
	return err == nil
}

// START: TestPs1PresentAsksNothing

func TestPs1PresentAsksNothing(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "mpv,yt-dlp,winget", Answer: "y"})
	if len(res.Prompts) != 0 || len(res.PM) != 0 || len(res.URLs) != 1 {
		t.Fatalf("prompts %v, winget/scoop calls %v, URLs %v; want no question, no call and only the goremi download", res.Prompts, res.PM, res.URLs)
	}
}

// END: TestPs1PresentAsksNothing

// START: TestPs1ListsThenAsks

func TestPs1ListsThenAsks(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "winget"})
	for _, want := range []string{"winget install --id shinchiro.mpv", "yt-dlp.exe"} {
		if !strings.Contains(res.Stdout, want) {
			t.Errorf("stdout = %q, want the command list to hold %q", res.Stdout, want)
		}
	}
	if len(res.Prompts) != 1 || res.Prompts[0] != "Install mpv and yt-dlp? [y/N]" {
		t.Fatalf("prompts = %v, want one question: Install mpv and yt-dlp? [y/N]", res.Prompts)
	}
}

// END: TestPs1ListsThenAsks

// START: no install

// ps1NoInstall checks a run whose answer is not yes: nothing installed, exit 0.
func ps1NoInstall(t *testing.T, r DepsRun) {
	t.Helper()
	r.Have = "winget"
	res := runDeps(t, r)
	if res.Exit != 0 || len(res.PM) != 0 || res.exists("yt-dlp.exe") {
		t.Fatalf("exit %d, winget/scoop calls %v, yt-dlp.exe %v; want exit 0 and nothing installed", res.Exit, res.PM, res.exists("yt-dlp.exe"))
	}
}

func TestPs1AnswerNInstallsNothing(t *testing.T)     { ps1NoInstall(t, DepsRun{Answer: "n"}) }
func TestPs1EmptyAnswerInstallsNothing(t *testing.T) { ps1NoInstall(t, DepsRun{Answer: ""}) }
func TestPs1NoConsoleInstallsNothing(t *testing.T)   { ps1NoInstall(t, DepsRun{NoTTY: true}) }

// END: no install

// START: mpv by winget or scoop

func TestPs1YesInstallsMpvWithWinget(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "winget,yt-dlp", Answer: "y"})
	if !has(res.PM, "winget install --id shinchiro.mpv") {
		t.Fatalf("calls = %v, want winget install --id shinchiro.mpv", res.PM)
	}
}

func TestPs1YesInstallsMpvWithScoop(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "scoop,yt-dlp", Answer: "y"})
	if !has(res.PM, "scoop bucket add extras") || !has(res.PM, "scoop install mpv") {
		t.Fatalf("calls = %v, want scoop bucket add extras and scoop install mpv", res.PM)
	}
}

func TestPs1AnswerIsNeverRun(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "winget,yt-dlp", Answer: "y; New-Item -ItemType File pwned"})
	if len(res.PM) != 0 {
		t.Fatalf("calls = %v, want none: the answer is not exactly y", res.PM)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(res.Dir), "pwned")); err == nil {
		t.Fatal("the answer was run as a command")
	}
}

// END: mpv by winget or scoop

// START: yt-dlp exe

func TestPs1YtdlpExeDownloaded(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "mpv", Answer: "y"})
	want := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	got, _ := os.ReadFile(filepath.Join(res.Dir, "yt-dlp.exe"))
	if !has(res.URLs, want) || string(got) != ps1YtBody {
		t.Fatalf("URLs %v, file %q; want %s asked for and the served file in the install directory", res.URLs, got, want)
	}
}

func TestPs1YtdlpHttpsOnly(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "mpv", Answer: "y"})
	for _, u := range res.URLs {
		if !strings.HasPrefix(u, "https://") {
			t.Errorf("URL %q does not start with https://", u)
		}
	}
}

// END: yt-dlp exe

// START: TestPs1YesSwitchInstalls

func TestPs1YesSwitchInstalls(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "winget", Yes: true})
	if len(res.Prompts) != 0 || !has(res.PM, "winget install --id shinchiro.mpv") || !res.exists("yt-dlp.exe") {
		t.Fatalf("prompts %v, calls %v, yt-dlp.exe %v; want no question and both installed", res.Prompts, res.PM, res.exists("yt-dlp.exe"))
	}
}

// END: TestPs1YesSwitchInstalls

// START: TestPs1NoWingetNoScoopPrintsCommand

func TestPs1NoWingetNoScoopPrintsCommand(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "yt-dlp", Answer: "y"})
	if res.Exit != 0 || len(res.PM) != 0 || !res.exists("goremi.exe") {
		t.Fatalf("exit %d, calls %v, goremi.exe %v; want exit 0, no call and goremi installed", res.Exit, res.PM, res.exists("goremi.exe"))
	}
	for _, want := range []string{"winget install --id shinchiro.mpv", "scoop install mpv"} {
		if !strings.Contains(res.Stdout, want) {
			t.Errorf("stdout = %q, want %q", res.Stdout, want)
		}
	}
}

// END: TestPs1NoWingetNoScoopPrintsCommand

// START: TestPs1ChecksumMismatchRefused

func TestPs1ChecksumMismatchRefused(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "mpv", Answer: "y", BadSums: true})
	entries, _ := os.ReadDir(res.Dir)
	if res.exists("yt-dlp.exe") || !strings.Contains(strings.ToLower(res.Stderr), "checksum") || !res.exists("goremi.exe") || len(entries) != 1 {
		t.Fatalf("yt-dlp.exe %v, stderr %q, %d files; want no yt-dlp.exe, a checksum error and goremi.exe only", res.exists("yt-dlp.exe"), res.Stderr, len(entries))
	}
}

// END: TestPs1ChecksumMismatchRefused

// START: failed dependency

func TestPs1FailedMpvInstallKeepsGoremi(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "winget", Answer: "y", PMExit: "1"})
	if !res.exists("goremi.exe") || !strings.Contains(res.Stderr, "mpv") || !res.exists("yt-dlp.exe") {
		t.Fatalf("goremi.exe %v, yt-dlp.exe %v, stderr %q; want both installed and an error naming mpv", res.exists("goremi.exe"), res.exists("yt-dlp.exe"), res.Stderr)
	}
}

func TestPs1FailedYtdlpDownloadKeepsGoremi(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "mpv", Answer: "y", YtdlpFail: true})
	entries, _ := os.ReadDir(res.Dir)
	if !res.exists("goremi.exe") || !strings.Contains(res.Stderr, "yt-dlp") || len(entries) != 1 {
		t.Fatalf("goremi.exe %v, stderr %q, %d files; want goremi.exe only and an error naming yt-dlp", res.exists("goremi.exe"), res.Stderr, len(entries))
	}
}

// END: failed dependency

// START: TestPs1ChecksumMismatchLeavesNoPartial

func TestPs1ChecksumMismatchLeavesNoPartial(t *testing.T) {
	res := runDeps(t, DepsRun{Have: "mpv", Answer: "y", BadSums: true})
	entries, _ := os.ReadDir(res.Dir)
	if len(entries) != 1 || entries[0].Name() != "goremi.exe" {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("install directory holds %v, want only goremi.exe", names)
	}
}

// END: TestPs1ChecksumMismatchLeavesNoPartial
