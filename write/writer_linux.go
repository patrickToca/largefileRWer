//go:build linux

package write

import (
	"os"
	"golang.org/x/sys/unix"
)

// enableDirectIO enables direct I/O on Linux
func (w *Writer) enableDirectIO(file *os.File) error {
	fd := file.Fd()
	flags, err := unix.FcntlInt(fd, unix.F_GETFL, 0)
	if err != nil {
		return err
	}
	_, err = unix.FcntlInt(fd, unix.F_SETFL, flags|unix.O_DIRECT)
	if err != nil {
		return err
	}
	slog.Debug("Direct I/O enabled", "file", file.Name())
	return nil
}