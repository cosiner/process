package process

import "syscall"

func isProcessExist(pid int) bool {
	err := syscall.Kill(pid, syscall.Signal(0))
	return err == nil || err == syscall.EPERM
}
