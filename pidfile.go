package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type PIDFile string

func NewPIDFile(pathes ...string) PIDFile {
	var path string
	if len(pathes) != 0 {
		path = pathes[0]
	}

	if path == "" {
		path = os.Args[0]
	}
	path = strings.TrimSuffix(path, ".exe")
	if !strings.HasSuffix(path, ".pid") {
		path += ".pid"
	}
	return PIDFile(path)
}

func (p PIDFile) Path() string {
	return string(p)
}

func (p PIDFile) Read() (int, error) {
	fd, err := os.OpenFile(p.Path(), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return -1, err
	}
	defer fd.Close()

	return readPID(fd)
}

func (p PIDFile) Write() error {
	err := os.MkdirAll(filepath.Dir(p.Path()), 0755)
	if err != nil {
		return err
	}
	fd, err := os.OpenFile(p.Path(), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	pid, err := readPID(fd)
	if err != nil {
		return err
	}
	if pid > 0 && IsProcessExist(pid) {
		return fmt.Errorf("pid file exist and process with pid %d is alive", pid)
	}

	err = truncFile(fd)
	if err != nil {
		return err
	}
	_, err = fd.Write([]byte(strconv.Itoa(os.Getpid())))
	return err
}

func (p PIDFile) Remove() error {
	return os.Remove(p.Path())
}

func IsProcessExist(pid int) bool {
	return isProcessExist(pid)
}

func truncFile(fd *os.File) error {
	_, err := fd.Seek(0, 0)
	if err == nil {
		err = fd.Truncate(0)
	}
	return err
}

func readPID(fd *os.File) (int, error) {
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return -1, err
	}

	if len(data) == 0 {
		return -1, nil
	}
	return strconv.Atoi(string(data))
}
