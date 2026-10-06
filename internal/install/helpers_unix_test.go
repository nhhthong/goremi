//go:build !windows

// Helpers for the installer tests: the repository root, a sealed PATH with fake uname, curl and sudo, and a runner for install.sh.
package install

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// START: repoRoot

// repoRoot is the repository root, two directories above this package.
func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Dir(filepath.Dir(wd))
}

// END: repoRoot

// START: sealedPath

// systemTools are the real programs install.sh may use; everything else (mpv, yt-dlp, sudo) is absent from the sealed PATH.
var systemTools = []string{"mkdir", "chmod", "mv", "rm", "mktemp", "cat", "cp", "dirname", "basename", "tr", "sed", "ls"}

const fakeUname = `#!/bin/sh
case "$1" in
-s) printf '%s\n' "$FAKE_OS" ;;
-m) printf '%s\n' "$FAKE_ARCH" ;;
esac
`

// fakeCurl logs the URL, then fails or copies FAKE_BODY to its -o target as the environment says.
const fakeCurl = `#!/bin/sh
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
if [ -n "$FAKE_CURL_FAIL" ]; then
	[ -n "$FAKE_CURL_PARTIAL" ] && printf 'half' > "$out"
	exit 22
fi
cp "$FAKE_BODY" "$out"
`

const fakeSudo = `#!/bin/sh
printf 'sudo %s\n' "$*" >> "$SUDO_LOG"
exit 1
`

// sealedPath makes a directory that holds links to systemTools, the fake uname, curl and sudo, and mpv and yt-dlp when asked, and returns it.
func sealedPath(t *testing.T, mpv, ytdlp bool) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range systemTools {
		real, err := exec.LookPath(name)
		if err != nil {
			t.Fatalf("system tool %s not found: %v", name, err)
		}
		if err := os.Symlink(real, filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	tools := map[string]string{"uname": fakeUname, "curl": fakeCurl, "sudo": fakeSudo}
	if mpv {
		tools["mpv"] = "#!/bin/sh\nexit 0\n"
	}
	if ytdlp {
		tools["yt-dlp"] = "#!/bin/sh\nexit 0\n"
	}
	for name, body := range tools {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// END: sealedPath

// START: runInstall

// Run says how one install.sh run looks to the host: the fake uname answers, curl fails or not.
type Run struct {
	OS, Arch string
	Fail     bool // curl exits 22
	Partial  bool // with Fail: curl writes half a file first
	Body     string
	Mpv      bool   // a fake mpv is on PATH
	YtDlp    bool   // a fake yt-dlp is on PATH
	Stdin    bool   // the script text goes to sh on stdin, as in curl | sh
	Cut      int    // with Stdin: only the first Cut bytes of the script are sent (0 = all)
	Sealed   string // a PATH directory from sealedPath to reuse (empty = make one)
	PathHome bool   // ~/.local/bin of the home is on PATH
}

// Result is what a run left behind.
type Result struct {
	Home, Stdout, Stderr string
	Exit                 int
	CurlURLs             []string
	SudoCalls            []string
}

// runInstall runs install.sh in a sealed environment; the home is a new temp directory unless home is given.
func runInstall(t *testing.T, r Run, home string) Result {
	t.Helper()
	tmp := t.TempDir()
	if home == "" {
		home = filepath.Join(tmp, "home")
		if err := os.MkdirAll(home, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	body := filepath.Join(tmp, "body")
	if r.Body == "" {
		r.Body = "#!/bin/sh\necho goremi\n"
	}
	if err := os.WriteFile(body, []byte(r.Body), 0o644); err != nil {
		t.Fatal(err)
	}
	curlLog, sudoLog := filepath.Join(tmp, "curl.log"), filepath.Join(tmp, "sudo.log")
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(repoRoot(t), "install.sh")
	cmd := exec.Command(sh, script)
	if r.Stdin {
		text, err := os.ReadFile(script)
		if err != nil {
			t.Fatal(err)
		}
		if r.Cut > 0 && r.Cut < len(text) {
			text = text[:r.Cut]
		}
		cmd = exec.Command(sh)
		cmd.Stdin = bytes.NewReader(text)
	}
	cmd.Dir = tmp
	sealed := r.Sealed
	if sealed == "" {
		sealed = sealedPath(t, r.Mpv, r.YtDlp)
	}
	path := sealed
	if r.PathHome {
		path += ":" + filepath.Join(home, ".local", "bin")
	}
	cmd.Env = []string{"HOME=" + home, "PATH=" + path, "FAKE_OS=" + r.OS, "FAKE_ARCH=" + r.Arch, "FAKE_BODY=" + body, "CURL_LOG=" + curlLog, "SUDO_LOG=" + sudoLog}
	if r.Fail {
		cmd.Env = append(cmd.Env, "FAKE_CURL_FAIL=1")
	}
	if r.Partial {
		cmd.Env = append(cmd.Env, "FAKE_CURL_PARTIAL=1")
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	res := Result{Home: home}
	if err := cmd.Run(); err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("install.sh did not run: %v", err)
		}
		res.Exit = ee.ExitCode()
	}
	res.Stdout, res.Stderr = stdout.String(), stderr.String()
	res.CurlURLs, res.SudoCalls = lines(curlLog), lines(sudoLog)
	return res
}

// lines returns the non-empty lines of a file, none when it does not exist.
func lines(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return strings.Fields(strings.ReplaceAll(string(data), " ", " "))
}

// END: runInstall
