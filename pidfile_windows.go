// +build windows

package process

import "os"

func isProcessExist(pid int) bool {
	_, err := os.FindProcess(pid)
	return err == nil
}
