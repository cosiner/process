// +build linux solaris

package process

import (
	"os"
	"strconv"
	"syscall"
)

func isProcessExist(pid int) bool {
	_, err := os.Stat("/proc/" + strconv.Itoa(pid))
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	err = syscall.Kill(pid, syscall.Signal(0))
	return err == nil || err == syscall.EPERM
}
