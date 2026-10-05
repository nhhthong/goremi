//go:build windows

// IPC endpoint on Windows: a named pipe, one per process.
package player

import (
	"fmt"
	"net"
	"os"

	"github.com/Microsoft/go-winio"
)

// START: newEndpoint

// newEndpoint returns the pipe name \\.\pipe\goremi-<pid>; there is no directory to remove.
func newEndpoint() (string, func(), error) {
	return fmt.Sprintf(`\\.\pipe\goremi-%d`, os.Getpid()), func() {}, nil
}

// END: newEndpoint

// START: dial

func dial(endpoint string) (net.Conn, error) { return winio.DialPipe(endpoint, nil) }

// END: dial
