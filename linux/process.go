package linux

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type Process struct {
	Status ProcessStatus `json:"status"`
	Statm  ProcessStatm  `json:"statm"`
	Stat   ProcessStat   `json:"stat"`
	IO     ProcessIO     `json:"io"`
}

type ProcessStatus struct{}
type ProcessIO struct{}
type ProcessStatm struct{}
type ProcessStat struct{}

func ReadProcess(pid int, baseDir string) (*Process, error) {
	var pidDir = filepath.Join(baseDir, strconv.Itoa(pid))

	if _, err := os.Stat(pidDir); err != nil {
		return nil, err
	}
	process := Process{}

	var ioFile = filepath.Join(pidDir, "io")
	var statFile = filepath.Join(pidDir, "stat")
	var statmFile = filepath.Join(pidDir, "statm")
	var statusFile = filepath.Join(pidDir, "status")

	_ = ioFile
	_ = statFile
	_ = statmFile
	_ = statusFile

	return &process, nil
}

func ReadProcessStatus(path string) (*ProcessStatus, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	status := ProcessStatus{}
	_ = b
	return &status, nil
}

// process info reader
// https://github.com/sharyanto/scripts/blob/master/explain-proc-stat
