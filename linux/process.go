package linux

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func ReadMaxPID(path string) (uint64, error) {

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return 0, err
	}

	s := strings.TrimSpace(string(b))

	i, err := strconv.ParseUint(s, 10, 0)

	if err != nil {
		return 0, err
	}

	return i, nil

}

func ListPID(path string, max uint64) ([]uint64, error) {

	l := make([]uint64, 0, 5)

	for i := uint64(1); i <= max; i++ {

		p := filepath.Join(path, strconv.FormatUint(i, 10))

		s, err := os.Stat(p)

		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		if err != nil || !s.IsDir() {
			continue
		}

		l = append(l, i)

	}

	return l, nil
}

// process info reader
// https://github.com/sharyanto/scripts/blob/master/explain-proc-stat
