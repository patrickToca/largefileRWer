//go:build !linux

package write

import (
	"log/slog"
	"os"
	"runtime"
)

// enableDirectIO is a no-op on non-Linux platforms
func (w *Writer) enableDirectIO(file *os.File) error {
	slog.Debug("Direct I/O not supported on this platform, skipping",
		"file", file.Name(),
		"os", runtime.GOOS)
	return nil
}
