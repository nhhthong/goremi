//go:build !windows

// IPC endpoint on Linux and macOS: a socket inside a private directory.
package player

import (
	"net"
	"os"
	"path/filepath"
)

// START: newEndpoint

// newEndpoint makes a private directory (mode 0700, a new random path each time) and returns the socket path in it and a function that removes the directory.
func newEndpoint() (string, func(), error) {
	dir, err := os.MkdirTemp("", "goremi-")
	if err != nil {
		return "", nil, err
	}
	return filepath.Join(dir, "mpv.sock"), func() { _ = os.RemoveAll(dir) }, nil
}

// END: newEndpoint

// START: dial

func dial(endpoint string) (net.Conn, error) { return net.Dial("unix", endpoint) }

// END: dial
