// The log file of goremi.
package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// START: LogPath

// LogPath is the log file: <os.UserCacheDir()>/goremi/goremi.log.
func LogPath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "goremi", "goremi.log"), nil
}

// END: LogPath

// START: Logger

// Logger writes log lines, each starting with the PID.
type Logger struct {
	mu sync.Mutex
	w  io.Writer
}

// OpenLog opens the file in append mode. It never returns a nil logger: when the file cannot be opened the logger discards, and the error says why.
func OpenLog(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return &Logger{w: io.Discard}, err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return &Logger{w: io.Discard}, err
	}
	return &Logger{w: f}, nil
}

// Printf writes the message, one line per line of the message, each starting with the PID.
func (l *Logger) Printf(format string, a ...any) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, line := range strings.Split(strings.TrimSuffix(fmt.Sprintf(format, a...), "\n"), "\n") {
		fmt.Fprintf(l.w, "%d %s\n", os.Getpid(), line)
	}
}

// END: Logger
