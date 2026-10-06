// Tests of install.ps1, run by pwsh with the download and the tool lookup replaced by functions: the URL, the file, unsupported systems and the dependency report.
package install

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const ps1Wrapper = `function Invoke-WebRequest { param($Uri, $OutFile) Add-Content -LiteralPath $env:WEB_LOG -Value $Uri; Set-Content -LiteralPath $env:PROGRESS_LOG -Value $ProgressPreference; if ($env:FAKE_FAIL) { Set-Content -LiteralPath $OutFile -Value 'half'; throw 'network down' }; Copy-Item -LiteralPath $env:FAKE_BODY -Destination $OutFile }
function Get-Command { if (($env:FAKE_HAVE -split ',') -contains $args[0]) { [pscustomobject]@{ Name = $args[0] } } }
& $env:INSTALL_PS1
exit $LASTEXITCODE`

// START: runPs1

// Ps1Run says how one install.ps1 run looks: the architecture, the tools Get-Command finds, and the download.
type Ps1Run struct {
	Arch     string
	Arch6432 string // PROCESSOR_ARCHITEW6432, empty = unset
	Have     string // comma-separated tool names Get-Command finds
	Fail     bool   // Invoke-WebRequest writes half a file and throws
	Existing string // content of an already installed goremi.exe, empty = none
	PathHas  bool   // the install directory is in Path
}

// Ps1Result is what a run left behind.
type Ps1Result struct {
	LocalAppData, Stdout, Stderr, Progress string
	Exit                                   int
	URLs                                   []string
}

// runPs1 runs install.ps1 under pwsh; a missing pwsh fails the test, it does not skip it.
func runPs1(t *testing.T, r Ps1Run) Ps1Result {
	t.Helper()
	pwsh, err := exec.LookPath("pwsh")
	if err != nil {
		t.Fatalf("pwsh not found: %v", err)
	}
	tmp := t.TempDir()
	local := filepath.Join(tmp, "LocalAppData")
	installDir := filepath.Join(local, "Programs", "goremi")
	if err := os.MkdirAll(local, 0o755); err != nil {
		t.Fatal(err)
	}
	if r.Existing != "" {
		if err := os.MkdirAll(installDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(installDir, "goremi.exe"), []byte(r.Existing), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	body := filepath.Join(tmp, "body")
	if err := os.WriteFile(body, []byte("fake exe body"), 0o644); err != nil {
		t.Fatal(err)
	}
	log, progress := filepath.Join(tmp, "web.log"), filepath.Join(tmp, "progress.log")
	script, _ := filepath.Abs(filepath.Join("..", "..", "install.ps1"))
	path := os.Getenv("PATH")
	if r.PathHas {
		path += string(os.PathListSeparator) + installDir
	}
	cmd := exec.Command(pwsh, "-NoProfile", "-NonInteractive", "-Command", ps1Wrapper)
	cmd.Dir = tmp
	cmd.Env = append(os.Environ(), "DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=1", "LocalAppData="+local, "Path="+path, "PROCESSOR_ARCHITECTURE="+r.Arch, "PROCESSOR_ARCHITEW6432="+r.Arch6432, "FAKE_HAVE="+r.Have, "FAKE_BODY="+body, "WEB_LOG="+log, "PROGRESS_LOG="+progress, "INSTALL_PS1="+script)
	if r.Fail {
		cmd.Env = append(cmd.Env, "FAKE_FAIL=1")
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	res := Ps1Result{LocalAppData: local}
	if err := cmd.Run(); err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("pwsh did not run: %v", err)
		}
		res.Exit = ee.ExitCode()
	}
	res.Stdout, res.Stderr = stdout.String(), stderr.String()
	if data, err := os.ReadFile(progress); err == nil {
		res.Progress = strings.TrimSpace(string(data))
	}
	if data, err := os.ReadFile(log); err == nil {
		for _, l := range strings.Split(string(data), "\n") {
			if l = strings.TrimSpace(l); l != "" {
				res.URLs = append(res.URLs, l)
			}
		}
	}
	return res
}

// END: runPs1

// START: download and file

func TestPs1DownloadUrl(t *testing.T) {
	res := runPs1(t, Ps1Run{Arch: "AMD64", Have: "mpv,yt-dlp"})
	want := "https://github.com/nhhthong/goremi/releases/latest/download/goremi_windows_amd64.exe"
	if res.Exit != 0 || len(res.URLs) != 1 || res.URLs[0] != want {
		t.Fatalf("exit = %d, URLs = %q, stderr = %q; want one URL %s", res.Exit, res.URLs, res.Stderr, want)
	}
}

func TestPs1InstallsFile(t *testing.T) {
	res := runPs1(t, Ps1Run{Arch: "AMD64", Have: "mpv,yt-dlp"})
	got, err := os.ReadFile(filepath.Join(res.LocalAppData, "Programs", "goremi", "goremi.exe"))
	if res.Exit != 0 || err != nil || string(got) != "fake exe body" {
		t.Fatalf("exit = %d, file = %q, %v; want the served body", res.Exit, got, err)
	}
}

func TestPs1UnsupportedArch(t *testing.T) {
	res := runPs1(t, Ps1Run{Arch: "ARM64"})
	if res.Exit != 1 || !strings.Contains(res.Stderr+res.Stdout, "ARM64") {
		t.Fatalf("exit = %d, output = %q %q; want exit 1 naming ARM64", res.Exit, res.Stdout, res.Stderr)
	}
	if len(res.URLs) != 0 {
		t.Fatalf("downloaded: %q", res.URLs)
	}
	if _, err := os.Stat(filepath.Join(res.LocalAppData, "Programs")); err == nil {
		t.Fatal("Programs exists after a refused run")
	}
}

// END: download and file

// START: security

func TestPs1HostileArchNotExecuted(t *testing.T) {
	// the value names a file under a directory the test creates; if it ran, that file would exist
	probe := filepath.Join(t.TempDir(), "pwned")
	res := runPs1(t, Ps1Run{Arch: "AMD64; New-Item -ItemType File '" + probe + "'"})
	if res.Exit != 1 || len(res.URLs) != 0 {
		t.Fatalf("exit = %d, URLs = %q; want exit 1 and no download", res.Exit, res.URLs)
	}
	if _, err := os.Stat(probe); err == nil {
		t.Fatal("the architecture value was run as a command: pwned exists")
	}
}

func TestPs1HttpsOnly(t *testing.T) {
	res := runPs1(t, Ps1Run{Arch: "AMD64", Have: "mpv,yt-dlp"})
	if len(res.URLs) == 0 {
		t.Fatal("nothing was downloaded")
	}
	for _, u := range res.URLs {
		if !strings.HasPrefix(u, "https://") {
			t.Fatalf("URL %q does not start with https://", u)
		}
	}
}

func TestPs1WritesOnlyInstallDir(t *testing.T) {
	res := runPs1(t, Ps1Run{Arch: "AMD64", Have: "mpv,yt-dlp"})
	var files []string
	_ = filepath.WalkDir(res.LocalAppData, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(res.LocalAppData, p)
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	if len(files) != 1 || files[0] != "Programs/goremi/goremi.exe" {
		t.Fatalf("files under LocalAppData = %q, want only Programs/goremi/goremi.exe", files)
	}
}

// END: security

// START: dependency report

// ps1Reported runs a successful install and returns its stdout.
func ps1Reported(t *testing.T, have string) string {
	t.Helper()
	res := runPs1(t, Ps1Run{Arch: "AMD64", Have: have})
	if res.Exit != 0 {
		t.Fatalf("exit = %d, want 0 (stderr %q)", res.Exit, res.Stderr)
	}
	return res.Stdout
}

func TestPs1DepsBothMissing(t *testing.T) {
	out := ps1Reported(t, "")
	m, y := strings.Index(out, "mpv"), strings.Index(out, "yt-dlp")
	if m < 0 || y < 0 || m > y {
		t.Fatalf("stdout = %q, want a line for mpv, then one for yt-dlp", out)
	}
}

func TestPs1DepsOnlyMpvMissing(t *testing.T) {
	out := ps1Reported(t, "yt-dlp")
	if !strings.Contains(out, "mpv") || strings.Contains(out, "yt-dlp") {
		t.Fatalf("stdout = %q, want mpv only", out)
	}
}

func TestPs1DepsOnlyYtdlpMissing(t *testing.T) {
	out := ps1Reported(t, "mpv")
	if !strings.Contains(out, "yt-dlp") || strings.Contains(out, "mpv") {
		t.Fatalf("stdout = %q, want yt-dlp only", out)
	}
}

func TestPs1DepsNoneMissing(t *testing.T) {
	out := ps1Reported(t, "mpv,yt-dlp")
	if strings.Contains(out, "mpv") || strings.Contains(out, "yt-dlp") {
		t.Fatalf("stdout = %q, want neither name", out)
	}
}

// END: dependency report

// START: architecture

func TestPs1Wow64Installs(t *testing.T) {
	res := runPs1(t, Ps1Run{Arch: "x86", Arch6432: "AMD64", Have: "mpv,yt-dlp"})
	if res.Exit != 0 || len(res.URLs) != 1 || !strings.HasSuffix(res.URLs[0], "/goremi_windows_amd64.exe") {
		t.Fatalf("exit = %d, URLs = %q, stderr = %q; want exit 0 and one goremi_windows_amd64.exe download", res.Exit, res.URLs, res.Stderr)
	}
}

func TestPs1Real32bitRefused(t *testing.T) {
	res := runPs1(t, Ps1Run{Arch: "x86"})
	if res.Exit != 1 || !strings.Contains(res.Stderr+res.Stdout, "x86") || len(res.URLs) != 0 {
		t.Fatalf("exit = %d, URLs = %q, output = %q %q; want exit 1 naming x86 and no download", res.Exit, res.URLs, res.Stdout, res.Stderr)
	}
}

// END: architecture

// START: download behaviour

func TestPs1ProgressSilent(t *testing.T) {
	res := runPs1(t, Ps1Run{Arch: "AMD64", Have: "mpv,yt-dlp"})
	if res.Progress != "SilentlyContinue" {
		t.Fatalf("ProgressPreference during the download = %q, want SilentlyContinue", res.Progress)
	}
}

func TestPs1FailedDownloadKeepsOld(t *testing.T) {
	res := runPs1(t, Ps1Run{Arch: "AMD64", Fail: true, Existing: "old exe"})
	if res.Exit != 1 || strings.TrimSpace(res.Stderr) == "" {
		t.Fatalf("exit = %d, stderr = %q; want exit 1 and a message", res.Exit, res.Stderr)
	}
	dir := filepath.Join(res.LocalAppData, "Programs", "goremi")
	got, err := os.ReadFile(filepath.Join(dir, "goremi.exe"))
	if err != nil || string(got) != "old exe" {
		t.Fatalf("goremi.exe = %q, %v; want the old content", got, err)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 1 {
		t.Fatalf("install directory holds %d entries, want only goremi.exe", len(entries))
	}
}

// END: download behaviour

// START: hints

func TestPs1PathHintWhenMissing(t *testing.T) {
	out := ps1Reported(t, "mpv,yt-dlp")
	if !strings.Contains(out, filepath.Join("Programs", "goremi")) || !strings.Contains(out, "PATH") {
		t.Fatalf("stdout = %q, want a line naming the install directory and PATH", out)
	}
}

func TestPs1NoPathHintWhenPresent(t *testing.T) {
	res := runPs1(t, Ps1Run{Arch: "AMD64", Have: "mpv,yt-dlp", PathHas: true})
	if res.Exit != 0 || strings.Contains(res.Stdout, "PATH") {
		t.Fatalf("exit = %d, stdout = %q; want no PATH line", res.Exit, res.Stdout)
	}
}

func TestPs1MpvInstallHints(t *testing.T) {
	out := ps1Reported(t, "yt-dlp")
	for _, want := range []string{"winget install --id shinchiro.mpv", "scoop bucket add extras", "scoop install mpv"} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout = %q, want %q", out, want)
		}
	}
}

func TestPs1YtdlpInstallHints(t *testing.T) {
	out := ps1Reported(t, "mpv")
	for _, want := range []string{"winget install --id yt-dlp.yt-dlp", "scoop install yt-dlp"} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout = %q, want %q", out, want)
		}
	}
}

// END: hints
