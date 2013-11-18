package linux

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Process struct {
}

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

// process info reader
// https://github.com/sharyanto/scripts/blob/master/explain-proc-stat
