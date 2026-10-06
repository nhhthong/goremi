//go:build windows

// IPC endpoint on Windows: a named pipe, one per Start.
package player

import (
	"fmt"
	"net"
	"os"
	"sync/atomic"

	"github.com/Microsoft/go-winio"
)

// START: newEndpoint

// pipeSeq makes the pipe name differ per Start, so an mpv that is still shutting down never answers for the next one.
var pipeSeq atomic.Uint64

// newEndpoint returns the pipe name \\.\pipe\goremi-<pid>-<n>; there is no directory to remove.
func newEndpoint() (string, func(), error) {
	return fmt.Sprintf(`\\.\pipe\goremi-%d-%d`, os.Getpid(), pipeSeq.Add(1)), func() {}, nil
}

// END: newEndpoint

// START: dial

func dial(endpoint string) (net.Conn, error) { return winio.DialPipe(endpoint, nil) }

// END: dial
