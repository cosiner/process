package process

import (
	"os"
	"runtime"
	"testing"
)

func TestPidFile(t *testing.T) {
	p := NewPIDFile()

	err := p.Write()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Remove()

	pid, err := p.Read()
	if err != nil {
		t.Error(err)
		return
	}
	if pid != os.Getpid() {
		t.Error("illegal pidfile content")
		return
	}
	if !IsProcessExist(pid) {
		t.Error("check process exist failed")
		return
	}
}

func TestIsProcessExist(t *testing.T) {
	if runtime.GOOS != "windows" {
		if !IsProcessExist(1) {
			t.Fatal("pid 1 should exist")
		}
	}
}
