//go:build !linux

package read

import (
	"log/slog"
	"os"
)

// enableDirectIO is a no-op on non-Linux platforms
func (r *Reader) enableDirectIO(file *os.File) error {
	slog.Debug("Direct I/O not supported on this platform, skipping", "file", file.Name())
	return nil
}
