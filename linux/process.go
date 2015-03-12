package linux

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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

type ProcessIO struct {
	RChar               uint64 `json:"rchar" field:"rchar"`                                 // chars read
	WChar               uint64 `json:"wchar" field:"wchar"`                                 // chars written
	Syscr               uint64 `json:"syscr" field:"syscr"`                                 // read syscalls
	Syscw               uint64 `json:"syscw" field:"syscw"`                                 // write syscalls
	ReadBytes           uint64 `json:"read_bytes" field:"read_bytes"`                       // bytes read
	WriteBytes          uint64 `json:"write_bytes" field:"write_bytes"`                     // bytes written
	CancelledWriteBytes uint64 `json:"cancelled_write_bytes" field:"cancelled_write_bytes"` // bytes truncated
}

type ProcessStatm struct{}
type ProcessStat struct{}

func ReadProcess(pid uint64, path string) (*Process, error) {

	p := filepath.Join(path, strconv.FormatUint(pid, 10))

	if _, err := os.Stat(p); err != nil {
		return nil, err
	}

	process := Process{}

	ioPath := filepath.Join(p, "io")
	statPath := filepath.Join(p, "stat")
	statmPath := filepath.Join(p, "statm")
	statusPath := filepath.Join(p, "status")

	io, err := ReadProcessIO(ioPath)

	if err != nil {
		return nil, err
	}

	_ = statPath
	_ = statmPath
	_ = statusPath

	process.IO = *io

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

func ReadProcessIO(path string) (*ProcessIO, error) {

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	// Maps a io metric to its value (i.e. rchar --> 100000)
	m := map[string]uint64{}

	var io ProcessIO = ProcessIO{}

	lines := strings.Split(string(b), "\n")

	for _, line := range lines {

		if strings.Index(line, ": ") == -1 {
			continue
		}

		l := strings.Split(line, ": ")

		k := l[0]
		v, err := strconv.ParseUint(l[1], 10, 0)

		if err != nil {
			return nil, err
		}

		m[k] = v

	}

	e := reflect.ValueOf(&io).Elem()
	t := e.Type()

	for i := 0; i < e.NumField(); i++ {

		k := t.Field(i).Tag.Get("field")

		v, ok := m[k]

		if ok {
			e.Field(i).SetUint(v)
		}

	}

	return &io, nil
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
