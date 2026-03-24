//go:build !linux

package write

import (
	"os"
	"log/slog"
)

// enableDirectIO is a no-op on non-Linux platforms
func (w *Writer) enableDirectIO(file *os.File) error {
	slog.Debug("Direct I/O not supported on this platform, skipping", "file", file.Name())
	return nil
}