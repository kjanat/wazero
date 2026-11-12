package sysfs

import (
	"io"
	"io/fs"
	"syscall"
	"testing"
	"time"

	experimentalsys "github.com/kjanat/wazero/experimental/sys"
	"github.com/kjanat/wazero/internal/testing/require"
)

type fakeDirEntry struct {
	name string
	err  error
}

func (e fakeDirEntry) Name() string      { return e.name }
func (e fakeDirEntry) IsDir() bool       { return false }
func (e fakeDirEntry) Type() fs.FileMode { return 0 }
func (e fakeDirEntry) Info() (fs.FileInfo, error) {
	if e.err != nil {
		return nil, e.err
	}
	return fakeFileInfo{name: e.name}, nil
}

type fakeFileInfo struct{ name string }

func (i fakeFileInfo) Name() string       { return i.name }
func (i fakeFileInfo) Size() int64        { return 0 }
func (i fakeFileInfo) Mode() fs.FileMode  { return 0 }
func (i fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (i fakeFileInfo) IsDir() bool        { return false }
func (i fakeFileInfo) Sys() any           { return nil }

// fakeDir returns its entries like os.File.ReadDir does: at most n per call
// when n > 0, and io.EOF once a positive n finds nothing left.
type fakeDir struct{ entries []fs.DirEntry }

func (d *fakeDir) ReadDir(n int) ([]fs.DirEntry, error) {
	if n <= 0 || n >= len(d.entries) {
		entries := d.entries
		d.entries = nil
		if n > 0 && len(entries) == 0 {
			return nil, io.EOF
		}
		return entries, nil
	}
	entries := d.entries[:n]
	d.entries = d.entries[n:]
	return entries, nil
}

func TestReadDirEntries(t *testing.T) {
	// On darwin, lstat of the kqueue descriptor in /dev/fd returns EBADF.
	// See https://go.dev/issue/80143
	unstattable := &fs.PathError{Op: "lstat", Path: "kqueue", Err: syscall.EBADF}
	removed := &fs.PathError{Op: "lstat", Path: "removed", Err: fs.ErrNotExist}

	names := func(dirents []experimentalsys.Dirent) (names []string) {
		for _, dirent := range dirents {
			names = append(names, dirent.Name)
		}
		return names
	}

	t.Run("skips entries the host cannot stat", func(t *testing.T) {
		dir := &fakeDir{entries: []fs.DirEntry{
			fakeDirEntry{name: "0"},
			fakeDirEntry{name: "kqueue", err: unstattable},
			fakeDirEntry{name: "removed", err: removed},
			fakeDirEntry{name: "1"},
		}}

		dirents, errno := readDirEntries(dir, "", -1)
		require.EqualErrno(t, 0, errno)
		require.Equal(t, []string{"0", "1"}, names(dirents))
	})

	t.Run("a skipped batch is not the end of the directory", func(t *testing.T) {
		dir := &fakeDir{entries: []fs.DirEntry{
			fakeDirEntry{name: "kqueue", err: unstattable},
			fakeDirEntry{name: "0"},
		}}

		dirents, errno := readDirEntries(dir, "", 1)
		require.EqualErrno(t, 0, errno)
		require.Equal(t, []string{"0"}, names(dirents))

		dirents, errno = readDirEntries(dir, "", 1)
		require.EqualErrno(t, 0, errno)
		require.Zero(t, len(dirents))
	})

	t.Run("only skipped entries read as an empty directory", func(t *testing.T) {
		dir := &fakeDir{entries: []fs.DirEntry{fakeDirEntry{name: "kqueue", err: unstattable}}}

		dirents, errno := readDirEntries(dir, "", 1)
		require.EqualErrno(t, 0, errno)
		require.Zero(t, len(dirents))
	})

	t.Run("other stat errors fail the read", func(t *testing.T) {
		dir := &fakeDir{entries: []fs.DirEntry{
			fakeDirEntry{name: "0"},
			fakeDirEntry{name: "broken", err: &fs.PathError{Op: "lstat", Path: "broken", Err: syscall.EIO}},
		}}

		_, errno := readDirEntries(dir, "", -1)
		require.EqualErrno(t, experimentalsys.EIO, errno)
	})
}
