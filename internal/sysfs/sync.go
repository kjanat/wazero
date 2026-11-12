//go:build !windows

package sysfs

import (
	"os"

	"github.com/kjanat/wazero/experimental/sys"
)

func fsync(f *os.File) sys.Errno {
	return sys.UnwrapOSError(f.Sync())
}
