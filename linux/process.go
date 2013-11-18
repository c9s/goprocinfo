package linux

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Process struct {
}

type ProcessStatus struct{}

func ReadProcess(pid int, baseDir string) (*Process, error) {
	var pidDir = filepath.Join(baseDir, strconv.Itoa(pid))

	if _, err := os.Stat(pidDir); err != nil {
		return nil, error
	}
	process := Process{}

	var ioFile = filepath.Join(pidDir, "io")
	var statFile = filepath.Join(pidDir, "stat")
	var statmFile = filepath.Join(pidDir, "statm")
	var statusFile = filepath.Join(pidDir, "status")

	return &process, nil
}

func ReadProcessStatus(path string) (*ProcessStatus, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	status := ProcessStatus{}
	return &status, nil
}

// process info reader
// https://github.com/sharyanto/scripts/blob/master/explain-proc-stat
